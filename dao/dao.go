package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/luxingwen/secret-game/conf"
	"github.com/luxingwen/secret-game/model"

	"context"
	"github.com/BurntSushi/toml"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

const (
	TableTeam        = "teams"
	TableTeamUser    = "team_user_maps"
	TableTeamTest    = "team_tests"
	TableTeamTestLog = "team_test_logs"
	TableSubject     = "subjects"
	TableWxUser      = "wx_users"
	TableWxCode      = "wx_codes"
)

var dao *Dao

func init() {
	var conf conf.Conf

	b, err := ioutil.ReadFile("conf/app.conf")
	if err != nil {
		panic(err)
	}

	if _, err := toml.Decode(string(b), &conf); err != nil {
		// handle error
		panic(err)
	}
	dao = NewDao(&conf)
	dao.DB.LogMode(true)

	dao.DB.AutoMigrate(&model.Team{}, &model.TeamUserMap{}, &model.Subject{}, &model.TeamTest{}, &model.TeamTestLog{}, &model.WxUser{}, &model.WxCode{})
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
