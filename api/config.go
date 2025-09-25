package api

// APIConfig 定义 API 服务器配置
type APIConfig struct {
	IP      string `yaml:"ip" json:"ip"`           // 服务器IP地址
	Port    int32  `yaml:"port" json:"port"`       // 服务器端口
	Timeout int32  `yaml:"timeout" json:"timeout"` // 超时时间（秒）
}
