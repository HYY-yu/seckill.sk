package config

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	MySQL struct {
		Base struct {
			MaxOpenConn     int           `toml:"maxOpenConn"`
			MaxIdleConn     int           `toml:"maxIdleConn"`
			ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
			Addr            string        `toml:"addr"`
			User            string        `toml:"user"`
			Pass            string        `toml:"pass"`
			Name            string        `toml:"name"`
		} `toml:"base"`
	} `toml:"mysql"`

	Redis struct {
		Addr         string `toml:"addr"`
		Pass         string `toml:"pass"`
		Db           int    `toml:"db"`
		MaxRetries   int    `toml:"maxRetries"`
		PoolSize     int    `toml:"poolSize"`
		MinIdleConns int    `toml:"minIdleConns"`
	} `toml:"redis"`

	JWT struct {
		Secret          string        `toml:"secret"`
		ExpireDuration  time.Duration `toml:"expireDuration"`
		Type            string        `toml:"type"`
		RefreshDuration time.Duration `toml:"refreshDuration"`
	} `toml:"jwt"`

	Log struct {
		LogPath    string `toml:"logPath"`
		Level      string `toml:"level"`
		Stdout     bool   `toml:"stdout"`
		JsonFormat bool   `toml:"jsonFormat"`
	} `toml:"log"`

	Server struct {
		ServerName string `toml:"serverName"`
		Host       string `toml:"host"`
		Pprof      bool   `toml:"pprof"`
		Grpc       struct {
			ShopHost string `toml:"shopHost"`
		} `toml:"grpc"`
	} `toml:"server"`

	Jaeger struct {
		UdpEndpoint string `toml:"udpEndpoint"`
		StdOut      bool   `toml:"stdOut"`
	} `toml:"jaeger"`

	ElasticJobETCD struct {
		Endpoints   []string      `toml:"endpoints"`
		DialTimeout time.Duration `toml:"dialTimeout"`
		Username    string        `toml:"username"`
		Password    string        `toml:"password"`
	} `toml:"elasticJobETCD"`
}

func InitConfig() {
	viper.SetConfigName("cfg")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./internal/service/sk/config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}

func Get() Config {
	return *config
}
