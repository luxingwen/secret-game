package conf

import (
	"time"
)

type Conf struct {
	DB *DBConfig `timl:"mysql"`
}

// Config mysql config.
type DBConfig struct {
	DSN         string        `toml:"dsn"`         // data source name.
	Active      int           `toml:"active"`      // pool
	Idle        int           `toml:"idle"`        // pool
	IdleTimeout time.Duration `toml:"idletimeout"` // connect max life time.
}
