package log

type LoggerConfig struct {
	Level  string `yaml:"level" json:"level"`   // trace, debug, info, warn, error, fatal, panic
	Format string `yaml:"format" json:"format"` // json, console, ethereum
	Output string `yaml:"output" json:"output"` // stdout, stderr, file
	Path   string `yaml:"path" json:"path"`     // log file path (when output is file)
}
