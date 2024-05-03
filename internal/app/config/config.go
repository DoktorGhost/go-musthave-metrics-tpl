package config

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	Host string
	Port int
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
	addr := new(Config)
	_ = flag.Value(addr)
	flag.Var(addr, "a", "Net address host:port")
	flag.Parse()

	if addr.Host == "" {
		addr.Host = "localhost"
	}
	if addr.Port == 0 {
		addr.Port = 8080
	}
	return addr
}

func ParseConfigClient() (*Config, int, int) {
	addr := new(Config)
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

	return addr, ReportInterval, PollInterval
}
