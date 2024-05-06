package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"strconv"
	"strings"
)

type Config struct {
	Host string
	Port int
}

type HostPort struct {
	Hp             []string `env:"ADDRESS" envSeparator:":"`
	ReportInterval int      `env:"REPORT_INTERVAL"`
	PollInterval   int      `env:"POLL_INTERVAL"`
}

var ReportInterval int
var PollInterval int

func (c *Config) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) Set(value string) error {
	hp := strings.Split(value, ":")
	if len(hp) != 2 {
		return fmt.Errorf("invalid host:port format: %s", value)
	}

	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	c.Host = hp[0]
	c.Port = port
	return nil
}

func ParseConfigServer() *Config {
	var cfg HostPort
	env.Parse(&cfg)
	addr := new(Config)
	if len(cfg.Hp) == 0 {
		_ = flag.Value(addr)
		flag.Var(addr, "a", "Net address host:port")
		flag.Parse()
		if addr.Host == "" {
			addr.Host = "localhost"
		}
		if addr.Port == 0 {
			addr.Port = 8080
		}
	} else {
		addr.Host = cfg.Hp[0]
		port, err := strconv.Atoi(cfg.Hp[1])
		if err != nil {
			return nil
		}
		addr.Port = port
	}

	return addr
}

func ParseConfigClient() (*Config, int, int) {
	var cfg HostPort
	env.Parse(&cfg)
	addr := new(Config)

	if len(cfg.Hp) == 0 {
		_ = flag.Value(addr)
		flag.Var(addr, "a", "Net address host:port")
		flag.IntVar(&ReportInterval, "r", 10, "частоту отправки метрик на сервер")
		flag.IntVar(&PollInterval, "p", 2, "частоту опроса метрик из пакета runtime")
		flag.Parse()
		if addr.Host == "" {
			addr.Host = "localhost"
		}
		if addr.Port == 0 {
			addr.Port = 8080
		}
	} else {
		addr.Host = cfg.Hp[0]
		port, err := strconv.Atoi(cfg.Hp[1])
		if err != nil {
			return nil, 0, 0
		}
		addr.Port = port
	}

	if cfg.ReportInterval != 0 {
		ReportInterval = cfg.ReportInterval
	}

	if cfg.PollInterval != 0 {
		PollInterval = cfg.PollInterval
	}

	return addr, ReportInterval, PollInterval
}
