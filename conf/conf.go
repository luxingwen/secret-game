package conf

import (
	"time"
)

type Conf struct {
	DB *DBConfig
}

// Config mysql config.
type DBConfig struct {
	DSN         string        // data source name.
	Active      int           // pool
	Idle        int           // pool
	IdleTimeout time.Duration // connect max life time.
}
