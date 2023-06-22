package hs

import (
	"errors"
	"reflect"
	"strconv"
)

var ErrNotVariable = errors.New("not a variable")

func setVariable[T any](ss []string,
	parseFunc func(string) (T, error), setFunc func(T)) error {
	if len(ss) > 1 {
		return ErrNotVariable
	} else if len(ss) == 0 {
		return nil
	}

	v, err := parseFunc(ss[0])
	if err != nil {
		return err
	}

	setFunc(v)
	return nil
}

func setInt64(v reflect.Value, isPtr bool, ss []string) error {
	setFunc := v.SetInt
	if isPtr {
		vPtr := reflect.New(v.Type().Elem())
		setFunc = vPtr.Elem().SetInt
		v.Set(vPtr)
	}
	return setVariable(ss, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	}, setFunc)
}

func setUint64(v reflect.Value, isPtr bool, ss []string) error {
	setFunc := v.SetUint
	if isPtr {
		vPtr := reflect.New(v.Type().Elem())
		setFunc = vPtr.Elem().SetUint
		v.Set(vPtr)
	}
	return setVariable(ss, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 64)
	}, setFunc)
}

func setFloat64(v reflect.Value, isPtr bool, ss []string) error {
	setFunc := v.SetFloat
	if isPtr {
		vPtr := reflect.New(v.Type().Elem())
		setFunc = vPtr.Elem().SetFloat
		v.Set(vPtr)
	}
	return setVariable(ss, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	}, setFunc)
}

func setBool(v reflect.Value, isPtr bool, ss []string) error {
	setFunc := v.SetBool
	if isPtr {
		vPtr := reflect.New(v.Type().Elem())
		setFunc = vPtr.Elem().SetBool
		v.Set(vPtr)
	}
	return setVariable(ss, strconv.ParseBool, setFunc)
}

func setComplex128(v reflect.Value, isPtr bool, ss []string) error {
	setFunc := v.SetComplex
	if isPtr {
		vPtr := reflect.New(v.Type().Elem())
		setFunc = vPtr.Elem().SetComplex
		v.Set(vPtr)
	}
	return setVariable(ss, func(s string) (complex128, error) {
		return strconv.ParseComplex(s, 128)
	}, setFunc)
}

func setString(v reflect.Value, isPtr bool, ss []string) error {
	setFunc := v.SetString
	if isPtr {
		vPtr := reflect.New(v.Type().Elem())
		setFunc = vPtr.Elem().SetString
		v.Set(vPtr)
	}
	return setVariable(ss, func(s string) (string, error) {
		return s, nil
	}, setFunc)
}
