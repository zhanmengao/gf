package tracetyp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhanmengao/gf/basedef"
	"go.opentelemetry.io/otel/attribute"
)

var (
	UseJson = true
)

func GetCurrentSpan(ctx context.Context) *Span {
	if ctx == nil {
		return nil
	}
	if span, ok := ctx.Value(basedef.KeyNameCurrentSpan).(*Span); ok {
		return span
	}
	return nil
}

func GetBaseSpan(ctx context.Context) *Span {
	if c := GetBaseCtx(ctx); c != nil {
		return GetCurrentSpan(c)
	}
	return nil
}

func GetBaseCtx(ctx context.Context) context.Context {
	iCtx := ctx.Value(basedef.KeyNameBaseSpanCtx)
	if iCtx != nil {
		if c, ok := iCtx.(context.Context); ok {
			return c
		}
	}
	return nil
}

func GetAttributes(key string, value interface{}) attribute.KeyValue {
	switch v := value.(type) {
	case string:
		return attribute.String(key, v)
	case []string:
		return attribute.StringSlice(key, v)
	case int:
		return attribute.Int(key, v)
	case []int:
		return attribute.IntSlice(key, v)
	case int64:
		return attribute.Int64(key, v)
	case []int64:
		return attribute.Int64Slice(key, v)
	case int32:
		return attribute.Int(key, int(v))
	case float64:
		return attribute.Float64(key, v)
	case []float64:
		return attribute.Float64Slice(key, v)
	case float32:
		return attribute.Float64(key, float64(v))
	case bool:
		return attribute.Bool(key, v)
	case []bool:
		return attribute.BoolSlice(key, v)
	case fmt.Stringer:
		return attribute.String(key, getString(v))
	case []fmt.Stringer:
		ss := make([]string, 0, len(v))
		for _, s := range v {
			if s != nil {
				ss = append(ss, getString(s))
			}
		}
		return attribute.StringSlice(key, ss)
	case []interface{}:
		return getStringSlice(key, v)
	case *[]interface{}:
		return getStringSlice(key, *v)
	default:
		return attribute.String(key, getInterface(v))
	}
}

func getInterface(t interface{}) string {
	if s, ok := t.(fmt.Stringer); ok {
		return getString(s)
	} else if t != nil {
		bt, _ := json.Marshal(t)
		return zerocopy.BtsToString(bt)
	} else {
		return "nil"
	}
}

func getString(ser fmt.Stringer) string {
	if UseJson {
		bt, _ := json.Marshal(ser)
		s := string(bt)
		return s
	} else {
		return ser.String()
	}
}

func getStringSlice(key string, slice []interface{}) attribute.KeyValue {
	if len(slice) == 0 {
		return attribute.String(key, "empty")
	} else if len(slice) == 1 {
		return attribute.String(key, getInterface(slice[0]))
	} else {
		ss := make([]string, 0, len(slice))
		for _, t := range slice {
			ss = append(ss, getInterface(t))
		}
		return attribute.StringSlice(key, ss)
	}
}
