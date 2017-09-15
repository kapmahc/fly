package nut

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/google/uuid"
)

const (
	// RoleAdmin admin role
	RoleAdmin = "admin"
	// RoleRoot root role
	RoleRoot = "root"
	// UserTypeEmail email user
	UserTypeEmail = "email"

	// DefaultResourceType default resource type
	DefaultResourceType = "-"
	// DefaultResourceID default resourc id
	DefaultResourceID = 0
)

// User user
type User struct {
	ID              uint       `orm:"column(id)" json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	UID             string     `orm:"column(uid)" json:"uid"`
	Password        string     `json:"-"`
	ProviderID      string     `orm:"column(provider_id)" json:"-"`
	ProviderType    string     `json:"providerType"`
	Logo            string     `json:"logo"`
	SignInCount     uint       `json:"signInCount"`
	LastSignInAt    *time.Time `json:"lastSignInAt"`
	LastSignInIP    string     `orm:"column(last_sign_in_ip)" json:"lastSignInIp"`
	CurrentSignInAt *time.Time `json:"currentSignInAt"`
	CurrentSignInIP string     `orm:"column(current_sign_in_ip)" json:"currentSignInIp"`
	ConfirmedAt     *time.Time `json:"confirmedAt"`
	LockedAt        *time.Time `json:"lockedAt"`
	UpdatedAt       time.Time  `orm:"auto_now" json:"updatedAt"`
	CreatedAt       time.Time  `orm:"auto_now_add" json:"createdAt"`
}

// TableName table name
func (*User) TableName() string {
	return "users"
}

// IsConfirm is confirm?
func (p *User) IsConfirm() bool {
	return p.ConfirmedAt != nil
}

// IsLock is lock?
func (p *User) IsLock() bool {
	return p.LockedAt != nil
}

//SetGravatarLogo set logo by gravatar
func (p *User) SetGravatarLogo() {
	buf := md5.Sum([]byte(strings.ToLower(p.Email)))
	p.Logo = fmt.Sprintf("https://gravatar.com/avatar/%s.png", hex.EncodeToString(buf[:]))
}

//SetUID generate uid
func (p *User) SetUID() {
	p.UID = uuid.New().String()
}

func (p User) String() string {
	return fmt.Sprintf("%s<%s>", p.Name, p.Email)
}

// Attachment attachment
type Attachment struct {
	ID           uint      `orm:"column(id)" json:"id"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Length       int64     `json:"length"`
	MediaType    string    `json:"mediaType"`
	ResourceID   uint      `orm:"column(resource_id)" json:"resourceId"`
	ResourceType string    `json:"resourceType"`
	UpdatedAt    time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt    time.Time `orm:"auto_now_add" json:"createdAt"`

	User *User `orm:"rel(fk)" json:"-"`
}

// TableName table name
func (*Attachment) TableName() string {
	return "attachments"
}

// IsPicture is picture?
func (p *Attachment) IsPicture() bool {
	return strings.HasPrefix(p.MediaType, "image/")
}

// Log log
type Log struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Message   string    `json:"message"`
	IP        string    `orm:"column(ip)" json:"ip"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`

	User *User `orm:"rel(fk)" json:"-"`
}

// TableName table name
func (*Log) TableName() string {
	return "logs"
}

func (p Log) String() string {
	return fmt.Sprintf("%s: [%s]\t %s", p.CreatedAt.Format(time.ANSIC), p.IP, p.Message)
}

// Policy policy
type Policy struct {
	ID        uint `orm:"column(id)" json:"id"`
	StartUp   time.Time
	ShutDown  time.Time
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`

	User *User `orm:"rel(fk)"`
	Role *Role `orm:"rel(fk)"`
}

//Enable is enable?
func (p *Policy) Enable() bool {
	now := time.Now()
	return now.After(p.StartUp) && now.Before(p.ShutDown)
}

// TableName table name
func (*Policy) TableName() string {
	return "policies"
}

// Role role
type Role struct {
	ID           uint `orm:"column(id)" json:"id"`
	Name         string
	ResourceID   uint `orm:"column(resource_id)"`
	ResourceType string
	UpdatedAt    time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt    time.Time `orm:"auto_now_add" json:"createdAt"`
}

// TableName table name
func (*Role) TableName() string {
	return "roles"
}

func (p Role) String() string {
	return fmt.Sprintf("%s@%s://%d", p.Name, p.ResourceType, p.ResourceID)
}

// Vote vote
type Vote struct {
	ID           uint `orm:"column(id)" json:"id"`
	Point        int
	ResourceID   uint `orm:"column(resource_id)"`
	ResourceType string
	UpdatedAt    time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt    time.Time `orm:"auto_now_add" json:"createdAt"`
}

// TableName table name
func (*Vote) TableName() string {
	return "votes"
}

// LeaveWord leave-word
type LeaveWord struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`
}

// TableName table name
func (*LeaveWord) TableName() string {
	return "leave_words"
}

// Link link
type Link struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Loc       string    `json:"loc"`
	Href      string    `json:"href"`
	Label     string    `json:"label"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`
}

// TableName table name
func (*Link) TableName() string {
	return "links"
}

// Card card
type Card struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Loc       string    `json:"loc"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Type      string    `json:"type"`
	Href      string    `json:"href"`
	Logo      string    `json:"logo"`
	SortOrder int       `json:"sortOrder"`
	Action    string    `json:"action"`
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`
}

// TableName table name
func (*Card) TableName() string {
	return "cards"
}

// FriendLink friend_links
type FriendLink struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Title     string    `json:"title"`
	Home      string    `json:"home"`
	Logo      string    `json:"logo"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`
}

// TableName table name
func (*FriendLink) TableName() string {
	return "friend_links"
}

func init() {
	orm.RegisterModel(
		new(User), new(Log), new(Role), new(Policy),
		new(Attachment), new(Vote),
		new(LeaveWord), new(FriendLink),
		new(Card), new(Link),
	)
}
