package modules

import "time"

type PostgreConfig struct {
	HOST        string
	Port        string
	Username    string
	Password    string
	DBName      string
	SSLMode     string
	ExecTimeout time.Duration
}
