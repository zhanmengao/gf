package cfgtp

type IServerConfig interface {
	GetNet() ListenConfig
}

// ListenConfig 网络侦听器
type ListenConfig struct {
	TCPPort  int `yaml:"tcp_port" json:"tcp_port"`
	HTTPPort int `yaml:"http_port" json:"http_port"`
	GRPCPort int `yaml:"grpc_port" json:"grpc_port"`
	UDPPort  int `yaml:"udp_port" json:"udp_port"`
}

// ServerConfigCommon 服务器每个srv的yml 配置的格式, 业务中自己一定要定义类似的结构，
type ServerConfigCommon struct {
	Net        ListenConfig `json:"net" yaml:"net"`                 // 必须定义！！！！
	LogPath    string       `json:"log_path" yaml:"log_path"`       // 可选定义，根据业务需要
	LogLevel   string       `json:"log_level" yaml:"log_level"`     // 可选定义，根据业务需要
	WorkerNum  int          `json:"worker_num" yaml:"worker_num"`   // 可选定义，根据业务需要
	KafkaAddr  string       `json:"kafka_addr" yaml:"kafka_addr"`   // 可选定义，根据业务需要
	KafkaTopic string       `json:"kafka_topic" yaml:"kafka_topic"` // 可选定义，根据业务需要
}

func (s *ServerConfigCommon) GetNet() ListenConfig {
	return s.Net
}
