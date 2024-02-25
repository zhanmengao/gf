package cfgtp

import "strings"

type DBType int

const (
	DBTypeUnknown DBType = -1
	DBTypeRedis   DBType = 1
	DBTypeMysql   DBType = 2
)

func (d DBType) String() string {
	switch d {
	case DBTypeRedis:
		return "redis"
	case DBTypeMysql:
		return "mysql"
	}
	return "unknown"
}

func DBTypeFromString(s string) DBType {
	ss := strings.ToLower(s)
	switch ss {
	case "redis":
		return DBTypeRedis
	case "mysql":
		return DBTypeMysql

	}
	return DBTypeUnknown
}

//DBConfig 数据库配置
type DBConfig struct {
	Type           string   `json:"type" yaml:"type"`
	Addrs          []string `json:"addr" yaml:"addr"`
	DBName         string   `json:"db_name" yaml:"db_name"`
	User           string   `json:"user" yaml:"user"`
	Password       string   `json:"password" yaml:"password"`
	TimeoutMs      int      `json:"timeout_ms" yaml:"timeout_ms"`
	RetryNum       int      `json:"retry_num" yaml:"retry_num" `
	MaxActive      int      `json:"max_active" yaml:"max_active" `
	ReadTimeoutMs  int      `json:"read_timeout_ms" yaml:"read_timeout_ms"`
	WriteTimeoutMs int      `json:"write_timeout_ms" yaml:"write_timeout_ms"`
	CharSet        string   `json:"charset" yaml:"charset"`
}
