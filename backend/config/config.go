package config

import (
	"backend/pkg/log"
	"backend/utils"

	"github.com/spf13/viper"
)

var cfg CfgInfo

type CfgInfo struct {
	Debug  bool
	Server struct {
		Addr string
	}

	Database struct {
		DSN string
	}
}

// Init 讀取設定並初始化
func Init(cfgPath string) (err error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if cfgPath == "" {
		const defaultCfgFName = "setting.yaml"
		v.SetConfigName(defaultCfgFName)
		v.AddConfigPath(utils.EXEDir())
		v.AddConfigPath("./")
	} else {
		v.SetConfigFile(cfgPath)
	}

	if err = v.ReadInConfig(); err != nil {
		log.Errorf("read config failed, err=%s", err)
		return
	}

	if err = v.Unmarshal(&cfg); err != nil {
		log.Errorf("unmarshal config failed, err=%s", err)
		return
	}

	return
}

// GetCfg 取得 config
func GetCfg() CfgInfo {
	return cfg
}
