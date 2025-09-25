package trace

// TraceConfig 定义链路追踪配置
type TraceConfig struct {
	JaegerHost string  `yaml:"jaeger_host" json:"jaeger_host"` // Jaeger 服务地址
	Enable     bool    `yaml:"enable" json:"enable"`           // 是否启用追踪
	SampleRate float32 `yaml:"sample_rate" json:"sample_rate"` // 采样率
}
