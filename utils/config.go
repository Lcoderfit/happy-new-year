package utils

import (
	"encoding/json"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/labstack/gommon/log"
)

var (
	// 缓存变量
	viewNumCache cache.Cache
	// 会话变量
	viewNumSession *session.Manager
)

// 链接redis配置
var redisConf = orm.Params{
	"key":      "viewNumCache",
	"conn":     "localhost:6379",
	"dbNum":    "0",
	"password": "124541",
}

// session配置
var sessionConfig = &session.ManagerConfig{
	CookieName:"viewNumSession",
	EnableSetCookie:true,
	Gclifetime:30*24*3600,
	Secure:false,
	CookieLifeTime:30*24*3600,
	ProviderConfig:"127.0.0.1:6379,100,124541,0",
}

func init() {
	BuildCache(redisConf)
	BuildSession()
}

// 建立缓存
func BuildCache(conf orm.Params) {
	confStr, err := json.Marshal(conf)
	if err != nil {
		log.Fatal("ConnRedis->Marshal Error: ", err)
		return
	}
	// 链接redis
	cacheTemp, err := cache.NewCache("redis", string(confStr))
	if err != nil {
		log.Fatal("ConnRedis->NewCache Error: ", err)
		return
	}
	viewNumCache = cacheTemp
}

// 建立会话
func BuildSession() {
	var err error
	viewNumSession, err = session.NewManager("redis", sessionConfig)
	if err != nil {
		log.Fatal("BuildSession->NewManager Error: ", err)
		return
	}
	go viewNumSession.GC()
}