package nut

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// SignIn set sign-in info
func SignIn(o orm.Ormer, lang, ip, email, password string) (*User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if HMAC().Chk([]byte(password), []byte(user.Password)) {
		AddLog(o, user, ip, Tr(lang, "nut.logs.user.sign-in.failed"))
		return nil, errors.New(Tr(lang, "nut.errors.user.email-password-not-match"))
	}

	if !user.IsConfirm() {
		return nil, errors.New(Tr(lang, "nut.errors.user.not-confirm"))
	}

	if user.IsLock() {
		return nil, errors.New(Tr(lang, "nut.errors.user.is-lock"))
	}

	AddLog(o, user, ip, Tr(lang, "nut.logs.user.sign-in.success"))
	user.SignInCount++
	user.LastSignInAt = user.CurrentSignInAt
	user.LastSignInIP = user.CurrentSignInIP
	now := time.Now()
	user.CurrentSignInAt = &now
	user.CurrentSignInIP = ip

	if _, err = o.QueryTable(user).
		Filter("id", user.ID).
		Update(orm.Params{
			"last_sign_in_at":    user.LastSignInAt,
			"last_sign_in_ip":    user.LastSignInIP,
			"current_sign_in_at": user.CurrentSignInAt,
			"current_sign_in_ip": user.CurrentSignInIP,
			"sign_in_count":      user.SignInCount,
			"updated_at":         time.Now(),
		}); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUID get user by uid
func GetUserByUID(uid string) (*User, error) {
	var u User

	if err := orm.NewOrm().QueryTable(&u).Filter("uid", uid).One(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByEmail get user by email
func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := orm.NewOrm().QueryTable(&user).
		Filter("provider_type", UserTypeEmail).
		Filter("provider_id", email).
		One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// AddLog add log
func AddLog(o orm.Ormer, user *User, ip, message string) error {
	_, err := o.Insert(&Log{
		User:    user,
		IP:      ip,
		Message: message,
	})
	return err
}

// AddEmailUser add email user
func AddEmailUser(o orm.Ormer, lang, ip, name, email, password string) (*User, error) {

	user := User{
		Email:           email,
		Password:        string(HMAC().Sum([]byte(password))),
		Name:            name,
		ProviderType:    UserTypeEmail,
		ProviderID:      email,
		LastSignInIP:    "0.0.0.0",
		CurrentSignInIP: "0.0.0.0",
	}
	user.SetUID()
	user.SetGravatarLogo()

	if _, err := o.Insert(&user); err != nil {
		return nil, err
	}
	AddLog(o, &user, ip, Tr(lang, "nut.logs.user.sign-up"))
	return &user, nil
}

// Authority get roles
func Authority(o orm.Ormer, user uint, rty string, rid uint) ([]string, error) {
	var items []*Role

	if _, err := o.QueryTable(new(Role)).
		Filter("resource_type", rty).
		Filter("resource_id", rid).
		All(&items); err != nil {
		return nil, err
	}
	var roles []string
	for _, r := range items {
		var pm Policy
		if err := o.QueryTable(&pm).
			Filter("role_id", r.ID).
			Filter("user_id", user).
			One(&pm); err != nil {
			return nil, err
		}
		if pm.Enable() {
			roles = append(roles, r.Name)
		}
	}
	return roles, nil
}

//Is is role ?
func Is(o orm.Ormer, user uint, names ...string) bool {
	for _, name := range names {
		if Can(o, user, name, "-", 0) {
			return true
		}
	}
	return false
}

//Can can?
func Can(o orm.Ormer, user uint, name string, rty string, rid uint) bool {
	var r Role

	if err := o.QueryTable(&r).
		Filter("name", name).
		Filter("resource_type", rty).
		Filter("resource_id", rid).
		One(&r); err != nil {
		return false
	}
	var pm Policy
	if err := o.QueryTable(&pm).
		Filter("user_id", user).Filter("role_id", r.ID).One(&pm); err != nil {
		return false
	}

	return pm.Enable()
}

// GetRole create role if not exist
func GetRole(o orm.Ormer, name string, rty string, rid uint) (*Role, error) {
	r := Role{}

	err := o.QueryTable(&r).
		Filter("name", name).
		Filter("resource_type", rty).
		Filter("resource_id", rid).
		One(&r)
	if err == nil {
		return &r, nil
	}

	if err != orm.ErrNoRows {
		return nil, err
	}

	r.Name = name
	r.ResourceID = rid
	r.ResourceType = rty
	if _, err = o.Insert(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

//Deny deny permission
func Deny(o orm.Ormer, role uint, user uint) error {
	_, err := o.QueryTable(new(Policy)).
		Filter("role_id", role).
		Filter("user_id", user).
		Delete()
	return err
}

//Allow allow permission
func Allow(o orm.Ormer, user *User, role *Role, years, months, days int) error {
	begin := time.Now()
	end := begin.AddDate(years, months, days)

	var pm Policy
	err := o.QueryTable(&pm).
		Filter("role_id", role.ID).
		Filter("user_id", user.ID).
		One(&pm)
	if err == nil {
		_, err = o.QueryTable(&pm).
			Filter("id", pm.ID).
			Update(orm.Params{
				"start_up":   begin,
				"shut_down":  end,
				"updated_at": time.Now(),
			})

	} else if err == orm.ErrNoRows {
		pm.User = user
		pm.Role = role
		pm.StartUp = begin
		pm.ShutDown = end
		_, err = o.Insert(&pm)
	}
	return err
}

// ListUserByResource list users by resource
func ListUserByResource(o orm.Ormer, role, rty string, rid uint) ([]uint, error) {
	ror, err := GetRole(o, role, rty, rid)
	if err != nil {
		return nil, err
	}

	var ids []uint
	var policies []Policy
	if _, err := o.QueryTable(new(Policy)).
		Filter("role_id", ror.ID).All(&policies); err != nil {
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
func ListResourcesIds(o orm.Ormer, user uint, role, rty string) ([]uint, error) {
	var ids []uint
	var policies []Policy

	if _, err := o.QueryTable(new(Policy)).
		Filter("user", user).
		All(&policies); err != nil {
		return nil, err
	}
	for _, pm := range policies {
		if pm.Enable() {
			var ror Role
			if err := o.QueryTable(&ror).
				Filter("id", pm.Role.ID).
				One(&ror); err != nil {
				return nil, err
			}
			if ror.Name == role && ror.ResourceType == rty {
				ids = append(ids, ror.ResourceID)
			}
		}
	}
	return ids, nil
}

func confirmUser(o orm.Ormer, lang, ip string, user *User) error {
	now := time.Now()
	if _, err := o.QueryTable(user).Filter("id", user.ID).Update(orm.Params{
		"confirmed_at": now,
		"updated_at":   now,
	}); err != nil {
		return err
	}
	return AddLog(o, user, ip, Tr(lang, "nut.logs.user.confirm"))
}

func setUserPassword(o orm.Ormer, lang, ip string, user *User, password string) error {
	now := time.Now()
	if _, err := o.QueryTable(user).Filter("id", user.ID).Update(orm.Params{
		"password":   string(HMAC().Sum([]byte(password))),
		"updated_at": now,
	}); err != nil {
		return err
	}
	return AddLog(o, user, ip, Tr(lang, "nut.logs.user.change-password"))
}
