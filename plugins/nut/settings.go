package nut

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego/orm"
)

// Setting k-v
type Setting struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Key       string    `json:"key"`
	Val       string    `json:"val"`
	Encode    bool      `json:"encode"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
}

// TableName table name
func (*Setting) TableName() string {
	return "settings"
}

// Set set k-v
func Set(o orm.Ormer, k string, v interface{}, e bool) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if e {
		if buf, err = AES().Encrypt(buf); err != nil {
			return err
		}
	}

	var it Setting
	err = o.QueryTable(&it).Filter("key", k).One(&it, "id")
	if err == nil {
		_, err = o.QueryTable(&it).
			Filter("id", it.ID).
			Update(orm.Params{
				"encode": e,
				"val":    string(buf),
			})
	} else if err == orm.ErrNoRows {
		it.Key = k
		it.Val = string(buf)
		it.Encode = e
		_, err = o.Insert(&it)
	}
	return err
}

// Get by key
func Get(k string, v interface{}) error {
	var it Setting
	o := orm.NewOrm()
	err := o.QueryTable(&it).
		Filter("key", k).
		One(&it, "val", "encode")
	if err != nil {
		return err
	}
	buf := []byte(it.Val)
	if it.Encode {
		buf, err = AES().Decrypt(buf)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(buf, v)
}

func init() {
	orm.RegisterModel(new(Setting))
}
