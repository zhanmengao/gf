package cfgtp

type IBasicConfig interface {
	GetAppName() string
	GetLogLevel() string
	GetLogPath() string
	GetWorkerNum() int64
	GetInnerTimeout() int64
	GetTimeZone() string
	GetCoordinatorAddr() string
	GetMonitorProcTime() bool
	GetEnv() string
	GetK8s() bool
	GetPushGateway() string
	GetFLogDAddr() string
	GetCompressType() int64
	GetRPCTimeout() int64
}

// BasicConfig 全局的配置，服务启动的时候会使用里面的参数
type BasicConfig struct {
	AppName         string `yaml:"app_name" json:"app_name"`                   //
	LogLevel        string `yaml:"log_level" json:"log_level"`                 // 日志级别
	LogPath         string `yaml:"log_path" json:"log_path"`                   // 日志路径
	WorkerNum       int64  `yaml:"worker_num" json:"worker_num"`               // 工作队列长度
	InnerTimeout    int64  `yaml:"inner_timeout" json:"inner_timeout"`         // 全局的内网超时，单位：毫秒
	TimeZone        string `yaml:"time_zone" json:"time_zone"`                 // 服务器时区
	CoordinatorAddr string `yaml:"coordinator_addr" json:"coordinator_addr"`   // zk等服务发现组件地址
	MonitorProcTime bool   `yaml:"monitor_proc_time" json:"monitor_proc_time"` // 是否开启proc超时监控
	Env             string `yaml:"env" json:"env"`                             // 当前的环境: test/beta/prod
	K8s             bool   `yaml:"k8s" json:"k8s"`                             // 是否运行在容器模式下
	PushGateway     string `yaml:"pushgateway" json:"pushgateway"`             //普罗米修斯网关
	FLogDAddr       string `yaml:"flogd_addr" json:"flogd_addr"`               //FLOGD地址
	CompressType    int64  `yaml:"compress_type" json:"compress_type"`         //FRPC压缩解压缩方式
	RPCTimeout      int64  `yaml:"rpc_timeout" json:"rpc_timeout"`             //RPC超时
}

func (b *BasicConfig) GetAppName() string {
	return b.AppName
}
func (b *BasicConfig) GetLogLevel() string {
	return b.LogLevel
}
func (b *BasicConfig) GetLogPath() string {
	return b.LogPath
}
func (b *BasicConfig) GetWorkerNum() int64 {
	return b.WorkerNum
}
func (b *BasicConfig) GetInnerTimeout() int64 {
	return b.InnerTimeout
}
func (b *BasicConfig) GetTimeZone() string {
	return b.TimeZone
}
func (b *BasicConfig) GetCoordinatorAddr() string {
	return b.CoordinatorAddr
}
func (b *BasicConfig) GetMonitorProcTime() bool {
	return b.MonitorProcTime
}
func (b *BasicConfig) GetEnv() string {
	return b.Env
}
func (b *BasicConfig) GetK8s() bool {
	return b.K8s
}
func (b *BasicConfig) GetPushGateway() string {
	return b.PushGateway
}
func (b *BasicConfig) GetFLogDAddr() string {
	return b.FLogDAddr
}
func (b *BasicConfig) GetCompressType() int64 {
	return b.CompressType
}
func (b *BasicConfig) GetRPCTimeout() int64 {
	return b.RPCTimeout
}
