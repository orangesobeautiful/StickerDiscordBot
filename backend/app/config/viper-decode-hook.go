package config

import (
	"encoding/base64"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/xerrors"
)

func StringToByteSlice() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf([]byte{}) {
			return data, nil
		}

		bs64, err := base64.StdEncoding.DecodeString(data.(string))
		if err != nil {
			return nil, xerrors.Errorf("decode base64 failed: %w", err)
		}

		return bs64, nil
	}
}
