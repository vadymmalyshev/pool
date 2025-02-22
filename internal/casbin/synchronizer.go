package casbin

import (
	"fmt"
	"time"

	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/auth"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	redisadapter "github.com/casbin/redis-adapter"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	// init postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const defaultDelay = 5

// Syncronizer is synzhronizing casbin data between Postgres and Redis
type Syncronizer struct {
	db           *gorm.DB
	redisConn    redigo.Conn
	redisAdapter *redisadapter.Adapter
	enforcer     *casbin.Enforcer
	// Delay is a time in seconds between synchronizations
	Delay int
}

// NewSynchronizer returns Syncronizer
func NewSynchronizer() (*Syncronizer, error) {
	db, err := config.Config.Admin.DB.Connect()

	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize IDP db")
	}

	redisConn, err := config.Config.Redis.ConnectCasbin()

	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize redis")
	}

	redisAdapter := config.Config.Redis.Casbin()

	a := gormadapter.NewAdapter("postgres", config.Config.Admin.DB.Connection(), true)
	enforcer := auth.NewAccessEnforcer(a)

	return &Syncronizer{
		db:           db,
		redisConn:    redisConn,
		redisAdapter: redisAdapter,
		enforcer:     enforcer,
		Delay:        defaultDelay,
	}, nil
}

// Start init synchronization between Postgres and Redis
func (s Syncronizer) Start(er chan error) {
	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %ds", s.Delay), s.copyRulesToRedis)
	c.Start()
}

func (s Syncronizer) copyRulesToRedis() {
	var timeDb time.Time

	row := s.db.Raw("SELECT tstamp FROM last_update").Row()
	row.Scan(&timeDb)

	timeRedis, err := redigo.String(s.redisConn.Do("GET", "last_update"))

	if err != nil {
		logrus.Debug("can't get last updates from redis")
	}

	if timeDb.String() != timeRedis {
		logrus.Debugf("sync time – db: %s redis: %s", timeDb.String(), timeRedis)

		s.redisConn.Do("SET", "last_update", timeDb)

		s.enforcer.LoadPolicy()
		model := s.enforcer.GetModel()
		s.redisAdapter.SavePolicy(model)
	}
}
