package svc

import (
	"simplest_script/core"
	"simplest_script/core/conf"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var lock sync.Mutex
var db = make(map[string]*gorm.DB)
var rds = make(map[string]*redis.Client)

func NewDb(adapter string) *gorm.DB {
	lock.Lock()
	defer lock.Unlock()
	retry := 3
	var g *gorm.DB
	var ok bool
	var err error
	var dbLink string
conn:
	if g, ok = db[adapter]; !ok {
		retry--
		switch adapter {
		case core.DBMain:
			dbLink = conf.Conf.Mysql.TestMain
		case core.DBConsole:
			dbLink = conf.Conf.Mysql.TestConsole
		default:
			dbLink = conf.Conf.Mysql.TestMain
		}
		//启动Gorm支持
		g, err = gorm.Open(mysql.Open(dbLink), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",   // 表名前缀，`User` 的表名应该是 `t_users`
				SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
			},
		})
		//如果出错就GameOver了
		if err != nil {
			hlog.Info("数据库连接失败", err)
		}

		db[adapter] = g
	}
	//g = g.Debug()
	gdb, _ := g.DB()

	if gdb.Ping() != nil && retry > 0 {
		delete(db, adapter)
		goto conn
	}

	return g
}

func NewRedis(adapter string) *redis.Client {
	lock.Lock()
	defer lock.Unlock()
	retry := 3
	var rdb *redis.Client
	var ok bool
conn:
	if rdb, ok = rds[adapter]; !ok {
		retry--
		switch adapter {
		case core.RDSData:
			rdb = redis.NewClient(&redis.Options{
				Addr:     conf.Conf.Redis.Data.Addr,
				Password: conf.Conf.Redis.Data.Pass,
				DB:       conf.Conf.Redis.Data.Db,
			})
		default:
			rdb = redis.NewClient(&redis.Options{
				Addr:     conf.Conf.Redis.Default.Addr,
				Password: conf.Conf.Redis.Default.Pass,
				DB:       conf.Conf.Redis.Default.Db,
			})
		}

		rds[adapter] = rdb
	}

	_, err := rdb.Ping().Result()

	if err != nil && retry > 0 {
		delete(rds, adapter)
		goto conn
	}

	return rdb
}
