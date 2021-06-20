package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/luxingwen/secret-game/conf"

	"context"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

var dao *Dao

func init() {
	var conf conf.Conf

	b, err := ioutil.ReadFile("conf/app.conf")
	if err != nil {
		panic(err)
	}

	if _, err := toml.Decode(b, &conf); err != nil {
		// handle error
		panic(err)
	}
	dao = NewDao(&conf)

	dao.AutoMigrate(&model.Team{}, &model.TeamUserMap{}, &model.Subject{}, &model.TeamTest{}, &model.TeamTestLog{})
}

func GetDao() *Dao {
	return dao
}

// Dao dao
type Dao struct {
	c  *conf.Conf
	DB *gorm.DB
}

func NewDao(c *conf.Conf) *Dao {
	return &Dao{
		c:  c,
		DB: NewMySQL(c.DB),
	}
}

// Ping verify server is ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.DB.DB().Ping(); err != nil {
		log.Error("dao.cloudDB.Ping() error(%v)", err)
		return
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.DB.Close()
}

func NewMySQL(c *conf.DBConfig) (db *gorm.DB) {
	db, err := gorm.Open("mysql", c.DSN)
	if err != nil {
		log.Error("db dsn(%s) error: %v", c.DSN, err)
		panic(err)
	}
	db.DB().SetMaxIdleConns(c.Idle)
	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout) / time.Second)
	return
}
