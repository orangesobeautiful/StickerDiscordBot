package config

import (
	"fmt"
	"strings"

	"backend/app/utils"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config interface {
	GetDebug() bool
	GetServer() Server
	GetDatabase() Database
	GetRedis() Redis
	GetVectorDatabase() VectorDatabase
	GetFullTextSearchDB() FullTextSearchDatabase
	GetOpenai() Openai
	GetObjectStorage() ObjectStorage
	GetDiscord() Discord
}

var _ Config = (*config)(nil)

type config struct {
	Debug bool

	Server *server

	Database *database

	Redis *redis

	VectorDatabase *vectorDatabase

	FullTextSearchDatabase *fullTextSearchDatabase

	Openai *openai

	ObjectStorage *objectStorage

	Discord *discord
}

func New() (configInterface Config, err error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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

func (c *config) GetVectorDatabase() VectorDatabase {
	return c.VectorDatabase
}

func (c *config) GetFullTextSearchDB() FullTextSearchDatabase {
	return c.FullTextSearchDatabase
}

func (c *config) GetOpenai() Openai {
	return c.Openai
}

func (c *config) GetObjectStorage() ObjectStorage {
	return c.ObjectStorage
}

func (c *config) GetDiscord() Discord {
	return c.Discord
}
