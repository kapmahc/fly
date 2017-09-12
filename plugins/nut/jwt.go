package nut

import (
	"net/http"
	"sync"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
	"github.com/astaxie/beego"
	"github.com/google/uuid"
)

var (
	_jwt    *Jwt
	jwtOnce sync.Once
)

// JWT get a jwt
func JWT() *Jwt {
	jwtOnce.Do(func() {
		_jwt = &Jwt{
			key:    []byte(beego.AppConfig.String("jwtkey")),
			method: crypto.SigningMethodHS512,
		}
	})
	return _jwt
}

// Jwt jwt token helper
type Jwt struct {
	key    []byte
	method crypto.SigningMethod
}

//Parse parse jwt token
func (p *Jwt) Parse(buf []byte) (jwt.Claims, error) {
	tk, err := jws.ParseJWT(buf)
	if err != nil {
		return nil, err
	}
	if err = tk.Validate(p.key, p.method); err != nil {
		return nil, err
	}
	return tk.Claims(), nil
}

// ParseFromRequest parse token from http request
func (p *Jwt) ParseFromRequest(r *http.Request) (jwt.Claims, error) {
	tk, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		return nil, err
	}
	if err = tk.Validate(p.key, p.method); err != nil {
		return nil, err
	}
	return tk.Claims(), nil
}

// Generate generate jwt token
func (p *Jwt) Generate(cm jws.Claims, exp time.Duration) ([]byte, error) {
	kid := uuid.New().String()
	now := time.Now()
	cm.SetNotBefore(now)
	cm.SetExpiration(now.Add(exp))
	cm.Set("kid", kid)
	//TODO using kid

	jt := jws.NewJWT(cm, p.method)
	return jt.Serialize(p.key)
}
