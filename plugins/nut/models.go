package nut

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

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
	tableName       struct{}   `sql:"users"`
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	UID             string     `json:"uid"`
	Password        string     `json:"-"`
	ProviderID      string     `json:"-"`
	ProviderType    string     `json:"providerType"`
	Logo            string     `json:"logo"`
	SignInCount     uint       `json:"signInCount"`
	LastSignInAt    *time.Time `json:"lastSignInAt"`
	LastSignInIP    string     `json:"lastSignInIp"`
	CurrentSignInAt *time.Time `json:"currentSignInAt"`
	CurrentSignInIP string     `json:"currentSignInIp"`
	ConfirmedAt     *time.Time `json:"confirmedAt"`
	LockedAt        *time.Time `json:"lockedAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
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
	tableName    struct{}  `sql:"attachments"`
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Length       int64     `json:"length"`
	MediaType    string    `json:"mediaType"`
	ResourceID   uint      `json:"resourceId"`
	ResourceType string    `json:"resourceType"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`

	User *User
}

// IsPicture is picture?
func (p *Attachment) IsPicture() bool {
	return strings.HasPrefix(p.MediaType, "image/")
}

// Log log
type Log struct {
	tableName struct{}  `sql:"logs"`
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`

	User *User
}

func (p Log) String() string {
	return fmt.Sprintf("%s: [%s]\t %s", p.CreatedAt.Format(time.ANSIC), p.IP, p.Message)
}

// Policy policy
type Policy struct {
	tableName struct{} `sql:"policies"`
	ID        uint     `json:"id"`
	StartUp   time.Time
	ShutDown  time.Time
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	User *User
	Role *Role
}

//Enable is enable?
func (p *Policy) Enable() bool {
	now := time.Now()
	return now.After(p.StartUp) && now.Before(p.ShutDown)
}

// Role role
type Role struct {
	tableName    struct{} `sql:"roles"`
	ID           uint     `json:"id"`
	Name         string
	ResourceID   uint
	ResourceType string
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (p Role) String() string {
	return fmt.Sprintf("%s@%s://%d", p.Name, p.ResourceType, p.ResourceID)
}

// Vote vote
type Vote struct {
	tableName    struct{} `sql:"votes"`
	ID           uint     `json:"id"`
	Point        int
	ResourceID   uint
	ResourceType string
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

// LeaveWord leave-word
type LeaveWord struct {
	tableName struct{}  `sql:"leave_words"`
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

// Link link
type Link struct {
	tableName struct{}  `sql:"links"`
	ID        uint      `json:"id"`
	Loc       string    `json:"loc"`
	Href      string    `json:"href"`
	Label     string    `json:"label"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// Card card
type Card struct {
	tableName struct{}  `sql:"cards"`
	ID        uint      `json:"id"`
	Loc       string    `json:"loc"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Type      string    `json:"type"`
	Href      string    `json:"href"`
	Logo      string    `json:"logo"`
	SortOrder int       `json:"sortOrder"`
	Action    string    `json:"action"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// FriendLink friend_links
type FriendLink struct {
	tableName struct{}  `sql:"friend_links"`
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Home      string    `json:"home"`
	Logo      string    `json:"logo"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}
