package cfgtp

type CfgType int32

const (
	INI CfgType = iota << 1
	XLSX
	TXT
	Json
	YAML
	Toml
	VYaml
)

type WatchConfigCall func(key string, oldConfig interface{}, newConfig interface{})

// IConfig Key 只允许一级
type IConfig interface {
	Get(key string, cfg interface{}) (err error)
	GetAndWatch(key string, cfg interface{}, cb WatchConfigCall) (err error)
}
