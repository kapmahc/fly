package mail

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"time"

	"github.com/astaxie/beego/orm"
)

// https://www.linode.com/docs/email/postfix/email-with-postfix-dovecot-and-mysql
// http://wiki.dovecot.org/Authentication/PasswordSchemes
// https://mad9scientist.com/dovecot-password-creation-php/

// Domain domain
type Domain struct {
	ID        uint `orm:"column(id)" json:"id"`
	Name      string
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`
}

// TableName table name
func (*Domain) TableName() string {
	return "mail_domains"
}

// User user
type User struct {
	ID       uint `orm:"column(id)" json:"id"`
	FullName string
	Email    string
	Password string
	Enable   bool

	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`

	Domain *Domain `orm:"rel(fk)"`
}

// TableName table name
func (*User) TableName() string {
	return "mail_users"
}

func (p *User) sum(password string, salt []byte) string {
	buf := sha512.Sum512(append([]byte(password), salt...))
	return base64.StdEncoding.EncodeToString(append(buf[:], salt...))
}

// SetPassword set  password (SSHA512-CRYPT)
func (p *User) SetPassword(password string) error {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	p.Password = p.sum(password, salt)
	return nil
}

// ChkPassword check password
func (p *User) ChkPassword(password string) bool {
	buf, err := base64.StdEncoding.DecodeString(p.Password)
	if err != nil {
		return false
	}

	return len(buf) > sha512.Size && p.Password == p.sum(password, buf[sha512.Size:])
}

// Alias alias
type Alias struct {
	ID          uint `orm:"column(id)" json:"id"`
	Source      string
	Destination string
	UpdatedAt   time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt   time.Time `orm:"auto_now_add" json:"createdAt"`

	Domain *Domain `orm:"rel(fk)"`
}

// TableName table name
func (*Alias) TableName() string {
	return "mail_aliases"
}

func init() {
	orm.RegisterModel(new(Alias), new(User), new(Domain))
}
