package util

import (
	"github.com/spf13/cast"
	"github.com/zhanmengao/gf/errors"
	gxe "golang.org/x/xerrors"
	"reflect"
	"strings"
)

func SetField(val reflect.Value, str string) (err error) {
	typ := val.Type()
	switch typ.Kind() {
	case reflect.String:
		val.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var iv int64
		iv, err = cast.ToInt64E(str)
		if err != nil {
			err = errors.NewWarp(60, "Err Server Bad Param", "Err Server Bad Param", gxe.Errorf("WARP ERR: %w", err))
			return
		}
		val.SetInt(iv)
	case reflect.Slice:
		arr := strings.Split(str, ",")
		for i, v := range arr {
			arr[i] = strings.Trim(v, " \t\r\n")
		}
		var ia []int
		typName := typ.String()
		if typName == "[]int" {
			if ia, err = cast.ToIntSliceE(arr); err != nil {
				err = errors.NewWarp(60, "Err Server Bad Param", "Err Server Bad Param", gxe.Errorf("WARP ERR: %w", err))
				return
			}
			//如果不是int类型的，就直接挂掉吧
			val.Set(reflect.ValueOf(ia))
		} else if typName == "[]string" {
			val.Set(reflect.ValueOf(arr))
		} else {
			err = errors.New(60, "Err Server Bad Param", "Err Server Bad Param").Format("only []int and []string are supported")
			return
		}

	case reflect.Bool:
		var iv bool
		if iv, err = cast.ToBoolE(str); err != nil {
			err = errors.NewWarp(60, "Err Server Bad Param", "Err Server Bad Param", gxe.Errorf("WARP ERR: %w", err))
			return
		}
		val.SetBool(iv)
	default:
		err = errors.New(60, "Err Server Bad Param", "Err Server Bad Param").Format("unsupported config field name:%s, type:%s", typ.Name(), typ.Kind())

	}
	return
}
