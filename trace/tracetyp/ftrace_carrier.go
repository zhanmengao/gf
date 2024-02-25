package tracetyp

import (
	"bytes"
	"encoding/json"
	"github.com/zhanmengao/gf/proto/go/gerror"
)

// Carrier 是 TextMapPropagator 使用的存储介质
type Carrier map[string]interface{}

func NewCarrier(data ...map[string]interface{}) Carrier {
	if len(data) > 0 && data[0] != nil {
		return data[0]
	}
	return make(map[string]interface{})
}

func (c Carrier) Get(k string) string {
	return c[k].(string)
}

func (c Carrier) Set(k, v string) {
	c[k] = v
}

func (c Carrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

// MustMarshal .returns the JSON encoding of c
func (c Carrier) MustMarshal() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return b
}

func (c Carrier) String() string {
	return string(c.MustMarshal())
}

func (c Carrier) UnmarshalJSON(b []byte) error {
	carrier := NewCarrier(nil)
	err := c.unmarshalUseNumber(b, carrier)
	return err
}

func (c Carrier) unmarshalUseNumber(data []byte, v interface{}) (err error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err = decoder.Decode(v)
	if err != nil {
		err = gerror.ErrServerDecode().SetBasicErr(err)
		return
	}
	return
}
