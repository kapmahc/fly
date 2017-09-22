package nut

import (
	"time"

	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/kapmahc/fly/plugins/nut/i18n"
	"github.com/kapmahc/fly/plugins/nut/security"
)

// SignIn set sign-in info
func SignIn(lang, ip, email, password string) (*User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if !security.HmacChk([]byte(password), []byte(user.Password)) {
		AddLog(user, ip, i18n.T(lang, "nut.logs.user.sign-in.failed"))
		return nil, i18n.E(lang, "nut.errors.user.email-password-not-match")
	}

	if !user.IsConfirm() {
		return nil, i18n.E(lang, "nut.errors.user.not-confirm")
	}

	if user.IsLock() {
		return nil, i18n.E(lang, "nut.errors.user.is-lock")
	}

	AddLog(user, ip, i18n.T(lang, "nut.logs.user.sign-in.success"))
	user.SignInCount++
	user.LastSignInAt = user.CurrentSignInAt
	user.LastSignInIP = user.CurrentSignInIP
	now := time.Now()
	user.CurrentSignInAt = &now
	user.CurrentSignInIP = ip
	user.UpdatedAt = now

	if _, err = app.DB().Model(user).
		Column("last_sign_in_at",
			"last_sign_in_ip",
			"current_sign_in_at",
			"current_sign_in_ip",
			"sign_in_count",
			"updated_at",
		).Update(); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUID get user by uid
func GetUserByUID(uid string) (*User, error) {
	var u User
	if err := app.DB().Model(&u).Where("uid = ?", uid).Select(); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByEmail get user by email
func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := app.DB().Model(&user).
		Where("provider_type = ? AND provider_id = ?", UserTypeEmail, email).
		Select(); err != nil {
		return nil, err
	}
	return &user, nil
}

// AddLog add log
func AddLog(user *User, ip, message string) error {
	return app.DB().Insert(&Log{
		User:      user,
		IP:        ip,
		Message:   message,
		CreatedAt: time.Now(),
	})
}

// AddEmailUser add email user
func AddEmailUser(lang, ip, name, email, password string) (*User, error) {
	now := time.Now()
	user := User{
		Email:           email,
		Password:        security.HmacSum([]byte(password)),
		Name:            name,
		ProviderType:    UserTypeEmail,
		ProviderID:      email,
		LastSignInIP:    "0.0.0.0",
		CurrentSignInIP: "0.0.0.0",
		UpdatedAt:       now,
		CreatedAt:       now,
	}
	user.SetUID()
	user.SetGravatarLogo()

	if err := app.DB().Insert(&user); err != nil {
		return nil, err
	}
	AddLog(&user, ip, i18n.T(lang, "nut.logs.user.sign-up"))
	return &user, nil
}

// Authority get roles
func Authority(user uint, rty string, rid uint) ([]string, error) {
	var items []*Role
	db := app.DB()

	if err := db.Model(&items).
		Where("resource_type = ? AND resource_id = ?", rty, rid).
		Select(); err != nil {
		return nil, err
	}
	var roles []string
	for _, r := range items {
		var pm Policy
		if err := db.Model(&pm).
			Where("role_id = ? AND user_id = ?", r.ID, user).
			Select(); err != nil {
			return nil, err
		}
		if pm.Enable() {
			roles = append(roles, r.Name)
		}
	}
	return roles, nil
}

//Is is role ?
func Is(user uint, names ...string) bool {
	for _, name := range names {
		if Can(user, name, DefaultResourceType, DefaultResourceID) {
			return true
		}
	}
	return false
}

//Can can?
func Can(user uint, name string, rty string, rid uint) bool {
	var r Role
	db := app.DB()
	if err := db.Model(&r).
		Where("name = ? AND resource_type = ? AND resource_id = ?", name, rty, rid).
		Select(); err != nil {
		return false
	}
	var pm Policy
	if err := db.Model(&pm).
		Where("user_id = ? AND role_id = ?", user, r.ID).
		Select(); err != nil {
		return false
	}

	return pm.Enable()
}

// GetRole create role if not exist
func GetRole(name string, rty string, rid uint) (*Role, error) {
	r := Role{}
	db := app.DB()
	err := db.Model(&r).
		Where("name = ? AND resource_type = ? AND resource_id", name, rty, rid).
		Select(&r)
	if err == nil {
		return &r, nil
	}

	now := time.Now()
	r.Name = name
	r.ResourceID = rid
	r.ResourceType = rty
	r.CreatedAt = now
	r.UpdatedAt = now
	if err = db.Insert(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

//Deny deny permission
func Deny(role uint, user uint) error {
	_, err := app.DB().Model(new(Policy)).
		Where("role_id = ? AND user_id = ?", role, user).
		Delete()
	return err
}

//Allow allow permission
func Allow(user *User, role *Role, years, months, days int) error {
	now := time.Now()

	db := app.DB()
	var pm Policy
	err := db.Model(&pm).
		Where("role_id = ? AND user_id = ?", role.ID, user.ID).
		Select(&pm)
	pm.StartUp = now
	pm.ShutDown = now.AddDate(years, months, days)
	pm.UpdatedAt = now
	if err == nil {
		_, err = db.Model(&pm).
			Update("start_up",
				"shut_down",
				"updated_at",
			)
	} else {
		pm.User = user
		pm.Role = role
		pm.CreatedAt = now
		err = db.Insert(&pm)
	}
	return err
}

// ListUserByResource list users by resource
func ListUserByResource(role, rty string, rid uint) ([]uint, error) {
	ror, err := GetRole(role, rty, rid)
	if err != nil {
		return nil, err
	}

	var ids []uint
	var policies []Policy
	if err := app.DB().Model(&policies).
		Where("role_id = ?", ror.ID).
		Select(); err != nil {
		return nil, err
	}
	for _, pm := range policies {
		if pm.Enable() {
			ids = append(ids, pm.User.ID)
		}
	}
	return ids, nil
}

// ListResourcesIds list resource ids by user and role
func ListResourcesIds(user uint, role, rty string) ([]uint, error) {
	var ids []uint
	var policies []Policy
	db := app.DB()
	if err := db.Model(&policies).
		Where("user_id = ?", user).
		Select(); err != nil {
		return nil, err
	}
	for _, pm := range policies {
		if pm.Enable() {
			var ror Role
			if err := db.Model(&ror).
				Where("id", pm.Role.ID).
				Select(); err != nil {
				return nil, err
			}
			if ror.Name == role && ror.ResourceType == rty {
				ids = append(ids, ror.ResourceID)
			}
		}
	}
	return ids, nil
}

func confirmUser(lang, ip string, user *User) error {
	now := time.Now()
	user.UpdatedAt = now
	user.ConfirmedAt = &now
	if _, err := app.DB().Model(user).
		Where("id = ?", user.ID).Update(
		"confirmed_at",
		"updated_at",
	); err != nil {
		return err
	}
	return AddLog(user, ip, i18n.T(lang, "nut.logs.user.confirm"))
}

func setUserPassword(lang, ip string, user *User, password string) error {

	user.UpdatedAt = time.Now()
	user.Password = security.HmacSum([]byte(password))
	if _, err := app.DB().Model(user).
		Where("id = ?", user.ID).Update(
		"password",
		"updated_at",
	); err != nil {
		return err
	}
	return AddLog(user, ip, i18n.T(lang, "nut.logs.user.change-password"))
}
