package config

import (
	"time"
)

type Config struct {
	LogLvl string

	Postgresql
	HTTPServer
	Enrich
}

type HTTPServer struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CtxTimeout   time.Duration
}

type Postgresql struct {
	URL string
}

type Enrich struct {
	RequestTimeout time.Duration
}
