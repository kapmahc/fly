package vpn

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/astaxie/beego/orm"
)

// http://chagridsada.blogspot.com/2011/01/openvpn-system-based-on-userpass.html

// User user
type User struct {
	ID        uint `orm:"column(id)" json:"id"`
	FullName  string
	Email     string
	Details   string
	Password  string
	Online    bool
	Enable    bool
	StartUp   time.Time
	ShutDown  time.Time
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName table name
func (*User) TableName() string {
	return "vpn_users"
}

func (p *User) sum(password string, salt []byte) string {
	buf := md5.Sum(append([]byte(password), salt...))
	return base64.StdEncoding.EncodeToString(append(buf[:], salt...))
}

// SetPassword set  password (md5 with salt)
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

	return len(buf) > md5.Size && p.Password == p.sum(password, buf[md5.Size:])
}

// Log log
type Log struct {
	ID          uint   `orm:"column(id)" json:"id"`
	TrustedIP   string `orm:"column(trusted_id)"`
	TrustedPort uint
	RemoteIP    string `orm:"column(remote_ip)"`
	RemotePort  uint
	StartUp     time.Time
	ShutDown    *time.Time
	Received    float64
	Send        float64
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`

	User *User `orm:"rel(fk)"`
}

// TableName table name
func (*Log) TableName() string {
	return "vpn_logs"
}
func init() {
	orm.RegisterModel(new(User), new(Log))
}
