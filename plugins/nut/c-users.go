package nut

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	gomail "gopkg.in/gomail.v2"
)

const (
	actConfirm       = "confirm"
	actUnlock        = "unlock"
	actResetPassword = "reset-password"

	sendEmailJob = "nut.send-email"
)

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
	if err == nil {
		if fm.Password != fm.PasswordConfirmation {
			err = errors.New(Tr(lang, "nut.errors.user.passwords-not-match"))
		}
	}

	var user *User
	ip := p.Ctx.Input.IP()
	if err == nil {
		var cnt int64
		if cnt, err = o.QueryTable(user).
			Filter("provider_type", UserTypeEmail).
			Filter("provider_id", fm.Email).
			Count(); err == nil && cnt > 0 {
			err = errors.New(Tr(lang, "nut.errors.user.email-already-exist"))
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
	to, subject, body, sender := mail["to"], mail["subject"], mail["body"], mail["username"]
	if beego.BConfig.RunMode != beego.PROD {
		beego.Debug("send to", to, ": ", subject, "\n", body)
		return nil
	}
	smtp := make(map[string]interface{})
	if err := Get("site.smtp", &smtp); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dia := gomail.NewDialer(
		smtp["host"].(string),
		smtp["port"].(int),
		sender,
		smtp["password"].(string),
	)

	return dia.DialAndSend(msg)
}

func init() {
	JOBBER().Register(sendEmailJob, doSendMail)
}
