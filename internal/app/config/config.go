package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host            string
	Port            string
	ReportInterval  int
	PollInterval    int
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

type EnvStruct struct {
	Hp              []string `env:"ADDRESS" envSeparator:":"`
	ReportInterval  int      `env:"REPORT_INTERVAL"`
	PollInterval    int      `env:"POLL_INTERVAL"`
	StoreInterval   int      `env:"STORE_INTERVAL"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH"`
	Restore         bool     `env:"RESTORE"`
}

type HostPort struct {
	Host string
	Port int
}

func (c *HostPort) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *HostPort) Set(value string) error {
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
	var envStruct EnvStruct
	//считываем все переменны окружения в cfg
	env.Parse(&envStruct)

	config := new(Config)
	hostPort := new(HostPort)

	flag.Var(hostPort, "a", "Net address host:port")
	flag.IntVar(&config.StoreInterval, "i", 300, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "полное имя файла, куда сохраняются текущие значения")
	flag.BoolVar(&config.Restore, "r", true, "полное имя файла, куда сохраняются текущие значения")
	flag.Parse()

	_, exists := os.LookupEnv("ADDRESS")
	if exists {
		config.Host = envStruct.Hp[0]
		config.Port = envStruct.Hp[1]
	} else {
		if hostPort.Host == "" {
			config.Host = "localhost"
		}
		if hostPort.Port == 0 {
			config.Port = "8080"
		}
	}

	value, exists := os.LookupEnv("STORE_INTERVAL")
	if exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil
		}
		config.StoreInterval = intValue
	}

	value, exists = os.LookupEnv("FILE_STORAGE_PATH")
	if exists {
		config.FileStoragePath = value
	}

	_, exists = os.LookupEnv("RESTORE")
	if exists {
		config.Restore = envStruct.Restore
	}

	return config
}

func ParseConfigClient() *Config {
	var envStruct EnvStruct
	//считываем все переменны окружения в cfg
	env.Parse(&envStruct)

	config := new(Config)
	hostPort := new(HostPort)

	flag.Var(hostPort, "a", "Net address host:port")
	flag.IntVar(&config.ReportInterval, "r", 10, "частотa отправки метрик на сервер")
	flag.IntVar(&config.PollInterval, "p", 2, "частотa опроса метрик из пакета runtime")
	flag.Parse()

	_, exists := os.LookupEnv("ADDRESS")
	if exists {
		config.Host = envStruct.Hp[0]
		config.Port = envStruct.Hp[1]
	} else {
		if hostPort.Host == "" {
			config.Host = "localhost"
		}
		if hostPort.Port == 0 {
			config.Port = "8080"
		}
	}

	value, exists := os.LookupEnv("REPORT_INTERVAL")
	if exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil
		}
		config.ReportInterval = intValue
	}

	value, exists = os.LookupEnv("POLL_INTERVAL")
	if exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil
		}
		config.PollInterval = intValue
	}

	return config

}
