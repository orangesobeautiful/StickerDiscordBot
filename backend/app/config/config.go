package config

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"backend/app/utils"

	"github.com/spf13/viper"
)

type CfgInfo struct {
	Debug  bool
	Server struct {
		Addr       string
		ImgURL     string
		SessionKey *struct {
			UserAuth *sessionKeyInfo
		}
		Cookie *struct {
			MaxAge   time.Duration
			Secure   bool
			HTTPOnly bool
			SameSite http.SameSite
		}
		CORS *struct {
			AllowOrigins     []string
			AllowMethods     []string
			AllowHeaders     []string
			ExposeHeaders    []string
			AllowCredentials bool
			MaxAge           time.Duration
		}
	}
	Database *struct {
		DSN         string
		AutoMigrate bool
	}
	Redis *struct {
		Addr     string
		Username string
		Password string
		DB       int
	}
	ObjectStorage *struct {
		Endpoint        string
		BucketName      string
		AccessKeyID     string
		AccessKeySecret string
		PublicAccessURL string
	}
	Discord *struct {
		Token string
	}
}

type preprocessingInterface interface {
	Preprocessing() error
}

func preprocessingRecursive(fieldChain []string, v any) error {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)
		fieldChain = append(fieldChain, ft.Name)

		if fv.CanInterface() {
			fvItf := fv.Interface()
			if err := preprocessingRecursive(fieldChain, fvItf); err != nil {
				return err
			}

			if p, ok := fvItf.(preprocessingInterface); ok {
				if err := p.Preprocessing(); err != nil {
					return fmt.Errorf("preprocessing %s failed, err=%s", strings.Join(fieldChain, "."), err)
				}
			}
		}
		fieldChain = fieldChain[:len(fieldChain)-1]
	}

	return nil
}

// Init 讀取設定並初始化
func Init() (cfgInfo *CfgInfo, err error) {
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

	if err = v.Unmarshal(&cfgInfo); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	if err = preprocessingRecursive(nil, cfgInfo); err != nil {
		return nil, fmt.Errorf("preprocessing config failed: %w", err)
	}

	return
}
