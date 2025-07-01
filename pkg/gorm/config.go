package gorm

import "time"

type ConfigDB struct {
	DebugLog            bool
	SSL                 string
	Host                string
	Port                string
	User                string
	Pass                string
	DatabaseName        string
	Timezone            string
	PoolMaxIdleConn     int
	PoolMaxOpenConn     int
	PoolMaxConnLifetime time.Duration
}
