package config

import (
	"fmt"

	"backend/app/utils"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config interface {
	GetDebug() bool
	GetServer() Server
	GetDatabase() Database
	GetRedis() Redis
	GetObjectStorage() ObjectStorage
	GetDiscord() Discord
}

var _ Config = (*config)(nil)

type config struct {
	Debug         bool
	Server        *server
	Database      *database
	Redis         *redis
	ObjectStorage *objectStorage
	Discord       *discord
}

func New() (configInterface Config, err error) {
	v := viper.New()
	v.SetConfigType("yaml")

	cfgPath := ""
	if cfgPath == "" {
		const defaultCfgFName = "setting.yaml"
		v.SetConfigName(defaultCfgFName)
		v.AddConfigPath(utils.EXEDir())
		v.AddConfigPath("./")
	} else {
		v.SetConfigFile(cfgPath)
	}

	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	var cfg config
	if err = v.Unmarshal(&cfg, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			StringToByteSlice(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		))); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	return &cfg, nil
}

func (c *config) GetDebug() bool {
	return c.Debug
}

func (c *config) GetServer() Server {
	return c.Server
}

func (c *config) GetDatabase() Database {
	return c.Database
}

func (c *config) GetRedis() Redis {
	return c.Redis
}

func (c *config) GetObjectStorage() ObjectStorage {
	return c.ObjectStorage
}

func (c *config) GetDiscord() Discord {
	return c.Discord
}
