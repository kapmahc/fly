package nut

import (
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
)

var _cache cache.Cache
var cacheOnce sync.Once

// Cache get cache manager
func Cache() cache.Cache {
	cacheOnce.Do(func() {
		cm, err := cache.NewCache(
			beego.AppConfig.String("cacheprovider"),
			beego.AppConfig.String("cacheproviderconfig"),
		)
		if err != nil {
			beego.Error(err)
			return
		}
		_cache = cm
	})
	return _cache
}
