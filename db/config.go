package db

// PostgresConfig 定义 PostgreSQL 数据库配置
type PostgresConfig struct {
	Host         string `yaml:"host" json:"host"`                     // 数据库主机地址
	Port         string `yaml:"port" json:"port"`                     // 数据库端口
	User         string `yaml:"user" json:"user"`                     // 用户名
	Password     string `yaml:"password" json:"password"`             // 密码
	Database     string `yaml:"database" json:"database"`             // 数据库名
	Timeout      int64  `yaml:"timeout" json:"timeout"`               // 连接超时时间（毫秒）
	MaxOpenConns int32  `yaml:"max_open_conns" json:"max_open_conns"` // 最大打开连接数
	MaxIdleConns int32  `yaml:"max_idle_conns" json:"max_idle_conns"` // 最大空闲连接数
}

