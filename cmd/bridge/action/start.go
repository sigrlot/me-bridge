package action

import (
	"fmt"
	"os"

	"github.com/st-chain/me-bridge/log"
	"github.com/st-chain/me-bridge/server"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// StartAction 处理 start 命令
func StartAction(ctx *cli.Context) error {
	configPath := ctx.String("config")

	// 读取 YAML 配置文件
	yamlData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file '%s': %w", configPath, err)
	}

	var serverConfig *server.ServerConfig
	if err := yaml.Unmarshal(yamlData, &serverConfig); err != nil {
		return fmt.Errorf("failed to parse YAML config: %w", err)
	}

	// Initialize logger with config from file if available
	if serverConfig.Logger != nil {
		if err := log.SetLogger(serverConfig.Logger); err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}
	}

	log.Infof("Configuration %s loaded successfully", configPath)

	// 创建并启动服务器
	srv := server.NewServerWithConfig(serverConfig)

	log.Info("Starting me-bridge server...")
	if err := srv.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	log.Info("me-bridge server started successfully")
	return nil
}
