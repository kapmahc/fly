package nut

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"golang.org/x/text/language"
	gomail "gopkg.in/gomail.v2"
)

const (
	actConfirm       = "confirm"
	actUnlock        = "unlock"
	actResetPassword = "reset-password"

	sendEmailJob = "nut.send-email"
)

// GetUsersSignIn user sign in
// @router /users/sign-in [get]
func (p *Plugin) GetUsersSignIn() {
	p.LayoutApplication()
	p.Data[TITLE] = Tr(p.Locale(), "nut.users.sign-in.title")
	p.TplName = "nut/users/sign-in.html"
}

type fmSignIn struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

// PostUsersSignIn user sign in
// @router /users/sign-in [post]
func (p *Plugin) PostUsersSignIn() {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	lang := p.Locale()
	var fm fmSignIn
	err := p.ParseForm(&fm)

	var user *User
	if err == nil {
		user, err = SignIn(o, lang, p.Ctx.Input.IP(), fm.Email, fm.Password)
	}
	if err == nil {
		o.Commit()
		p.SetSession("uid", user.UID)
	} else {
		o.Rollback()
	}

	if p.Flash(func() string {
		return Tr(lang, "nut.users.confirm.success")
	}, err) {
		p.Redirect("nut.Plugin.GetHome")
	} else {
		p.Redirect("nut.Plugin.GetUsersSignIn")
	}
}

// GetUsersSignUp user sign up
// @router /users/sign-up [get]
func (p *Plugin) GetUsersSignUp() {
	p.LayoutApplication()
	p.Data[TITLE] = Tr(p.Locale(), "nut.users.sign-up.title")
	p.TplName = "nut/users/sign-up.html"
}

type fmSignUp struct {
	Name                 string `form:"name" valid:"Required"`
	Email                string `form:"email" valid:"Email"`
	Password             string `form:"password" valid:"MinSize(6)"`
	PasswordConfirmation string `form:"passwordConfirmation"`
}

func (p fmSignUp) Valid(v *validation.Validation) {
	if p.Password != p.PasswordConfirmation {
		v.SetError("PasswordConfirmation", Tr(language.AmericanEnglish.String(), "nut.errors.user.passwords-not-match"))
	}
}

// PostUsersSignUp user sign up
// @router /users/sign-up [post]
func (p *Plugin) PostUsersSignUp() {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	lang := p.Locale()
	var fm fmSignUp
	err := p.ParseForm(&fm)

	var user *User
	ip := p.Ctx.Input.IP()
	if err == nil {
		var cnt int64
		if cnt, err = o.QueryTable(user).
			Filter("provider_type", UserTypeEmail).
			Filter("provider_id", fm.Email).
			Count(); err == nil && cnt > 0 {
			err = Te(lang, "nut.errors.user.email-already-exist")
		}
	}
	if err == nil {
		user, err = AddEmailUser(o, lang, ip, fm.Name, fm.Email, fm.Password)
	}

	if err == nil {
		o.Commit()
		if er := p.sendEmail(lang, user, actConfirm); er != nil {
			beego.Error(er)
		}
	} else {
		o.Rollback()
	}

	p.Flash(func() string {
		return Tr(lang, "nut.users.confirm.success")
	}, err)
	p.Redirect("nut.Plugin.GetUsersSignUp")
}

// GetUsersConfirm user confirm
// @router /users/confirm [get]
func (p *Plugin) GetUsersConfirm() {
	p.getUsersEmailForm(actConfirm)
}

func (p *Plugin) getUsersEmailForm(act string) {
	p.LayoutApplication()
	p.Data[TITLE] = Tr(p.Locale(), "nut.users."+act+".title")
	p.TplName = "nut/users/email-form.html"
}

type fmEmail struct {
	Email string `form:"email" valid:"Email"`
}

// PostUsersConfirm user confirm
// @router /users/confirm [post]
func (p *Plugin) PostUsersConfirm() {
	lang := p.Locale()
	var fm fmEmail
	err := p.ParseForm(&fm)
	var user *User
	if err == nil {
		user, err = GetUserByEmail(fm.Email)
	}
	if err == nil && user.IsConfirm() {
		err = Te(lang, "nut.errors.user.already-confirm")
	}

	if err == nil {
		if er := p.sendEmail(lang, user, actConfirm); er != nil {
			beego.Error(er)
		}
	}

	p.Flash(func() string {
		return Tr(lang, "nut.users.confirm.success")
	}, err)
	p.Redirect("nut.Plugin.GetUsersConfirm")
}

// GetUsersConfirmToken confirm user
// @router /users/confirm/:token [get]
func (p *Plugin) GetUsersConfirmToken() {
	user, err := p.parseToken(actConfirm)
	lang := p.Locale()
	if err == nil && user.IsConfirm() {
		err = Te(lang, "nut.errors.user.already-confirm")
	}

	if err == nil {
		err = confirmUser(orm.NewOrm(), lang, p.Ctx.Input.IP(), user)
	}

	p.Flash(func() string {
		return Tr(lang, "nut.emails.user.confirm.success")
	}, err)
	p.Redirect("nut.Plugin.GetUsersSignIn")
}

func (p *Plugin) parseToken(act string) (*User, error) {
	cm, err := JWT().Parse([]byte(p.Ctx.Input.Param(":token")))
	if err != nil {
		return nil, err
	}
	if val := cm.Get("act").(string); val != act {
		return nil, Te(p.Locale(), "errors.bad-action")
	}
	return GetUserByUID(cm.Get("uid").(string))
}

// GetUsersUnlock user unlock
// @router /users/unlock [get]
func (p *Plugin) GetUsersUnlock() {
	p.getUsersEmailForm(actUnlock)
}

// PostUsersUnlock user unlock
// @router /users/unlock [post]
func (p *Plugin) PostUsersUnlock() {
	lang := p.Locale()
	var fm fmEmail
	err := p.ParseForm(&fm)
	var user *User
	if err == nil {
		user, err = GetUserByEmail(fm.Email)
	}
	if err == nil && user.IsConfirm() {
		err = Te(lang, "nut.errors.user.not-lock")
	}

	if err == nil {
		if er := p.sendEmail(lang, user, actUnlock); er != nil {
			beego.Error(er)
		}
	}

	p.Flash(func() string {
		return Tr(lang, "nut.users.unlock.success")
	}, err)
	p.Redirect("nut.Plugin.GetUsersUnlock")
}

// GetUsersUnlockToken unlock user
// @router /users/unlock/:token [get]
func (p *Plugin) GetUsersUnlockToken() {
	user, err := p.parseToken(actUnlock)
	lang := p.Locale()
	if err == nil && !user.IsLock() {
		err = Te(lang, "nut.errors.user.not-lock")
	}
	o := orm.NewOrm()
	if err == nil {
		err = o.Begin()
	}

	if err == nil {
		ip := p.Ctx.Input.IP()
		now := time.Now()
		_, err = o.QueryTable(user).Filter("id", user.ID).Update(orm.Params{
			"locked_at":  nil,
			"updated_at": now,
		})
		if err == nil {
			err = AddLog(o, user, ip, Tr(lang, "nut.logs.user.unlock"))
		}
	}
	if err == nil {
		o.Commit()
	} else {
		o.Rollback()
	}

	p.Flash(func() string {
		return Tr(lang, "nut.emails.user.unlock.success")
	}, err)
	p.Redirect("nut.Plugin.GetUsersSignIn")
}

// GetUsersForgotPassword user forgot password
// @router /users/forgot-password [get]
func (p *Plugin) GetUsersForgotPassword() {
	p.getUsersEmailForm(actResetPassword)
}

// PostUsersForgotPassword forgot password
// @router /users/forgot-password [post]
func (p *Plugin) PostUsersForgotPassword() {
	lang := p.Locale()
	var fm fmEmail
	err := p.ParseForm(&fm)
	var user *User
	if err == nil {
		user, err = GetUserByEmail(fm.Email)
	}

	if err == nil {
		if er := p.sendEmail(lang, user, actResetPassword); er != nil {
			beego.Error(er)
		}
	}

	p.Flash(func() string {
		return Tr(lang, "nut.users.forgot-password.success")
	}, err)
	p.Redirect("nut.Plugin.GetUsersForgotPassword")
}

// GetUsersResetPassword user reset password
// @router /users/reset-password/:token [get]
func (p *Plugin) GetUsersResetPassword() {
	p.LayoutApplication()
	p.Data[TITLE] = Tr(p.Locale(), "nut.users.reset-password.title")
	p.TplName = "nut/users/reset-password.html"
}

type fmResetPassword struct {
	Password             string `form:"password" valid:"MinSize(6)"`
	PasswordConfirmation string `form:"passwordConfirmation"`
}

func (p fmResetPassword) Valid(v *validation.Validation) {
	if p.Password != p.PasswordConfirmation {
		v.SetError("PasswordConfirmation", Tr(language.AmericanEnglish.String(), "nut.errors.user.passwords-not-match"))
	}
}

// PostUsersResetPassword reset user password
// @router /users/reset-password/:token [post]
func (p *Plugin) PostUsersResetPassword() {
	lang := p.Locale()
	var fm fmResetPassword
	err := p.ParseForm(&fm)
	var user *User
	if err == nil {
		user, err = p.parseToken(actResetPassword)
	}
	if err == nil {
		err = setUserPassword(orm.NewOrm(), p.Locale(), p.Ctx.Input.IP(), user, fm.Password)
	}

	if p.Flash(func() string {
		return Tr(lang, "nut.emails.user.reset-password.success")
	}, err) {
		p.Redirect("nut.Plugin.GetUsersSignIn")
	} else {
		p.Redirect("nut.Plugin.GetUsersResetPassword", ":token", p.Ctx.Input.Param(":token"))
	}
}

func (p *Plugin) sendEmail(lang string, user *User, act string) error {
	cm := jws.Claims{}
	cm.Set("act", act)
	cm.Set("uid", user.UID)
	tkn, err := JWT().Generate(cm, time.Hour*6)
	if err != nil {
		return err
	}
	obj := map[string]interface{}{
		"home":  p.HomeURL(),
		"token": string(tkn),
	}

	subject, err := Th(lang, fmt.Sprintf("nut.emails.user.%s.subject", act), obj)
	if err != nil {
		return err
	}
	body, err := Th(lang, fmt.Sprintf("nut.emails.user.%s.body", act), obj)
	if err != nil {
		return err
	}

	// -----------------------
	buf, err := json.Marshal(map[string]string{
		"to":      user.Email,
		"subject": subject,
		"body":    body,
	})
	if err != nil {
		return err
	}
	return JOBBER().Send(1, sendEmailJob, buf)
}

func doSendMail(buf []byte) error {
	var mail map[string]string
	if err := json.Unmarshal(buf, &mail); err != nil {
		return err
	}

	// ---------------------
	to, subject, body := mail["to"], mail["subject"], mail["body"]
	if beego.BConfig.RunMode != beego.PROD {
		beego.Debug("send to", to, ": ", subject, "\n", body)
		return nil
	}
	var smtp SMTP
	if err := Get("site.smtp", &smtp); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", smtp.Sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dia := gomail.NewDialer(
		smtp.Host,
		smtp.Port,
		smtp.Sender,
		smtp.Password,
	)

	return dia.DialAndSend(msg)
}

func init() {
	JOBBER().Register(sendEmailJob, doSendMail)
}
