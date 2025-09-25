package action

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// ConfigExampleAction 处理 config example 命令
func ConfigExampleAction(ctx *cli.Context) error {
	outputPath := ctx.String("output")
	force := ctx.Bool("force")

	// 源文件路径
	sourceFile := "config/config.example.yaml"

	// 检查源文件是否存在
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return fmt.Errorf("example config file not found: %s", sourceFile)
	}

	// 检查输出文件是否已存在
	if _, err := os.Stat(outputPath); err == nil && !force {
		return fmt.Errorf("output file already exists: %s (use --force to overwrite)", outputPath)
	}

	// 创建输出目录（如果不存在）
	if dir := filepath.Dir(outputPath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// 复制文件
	if err := copyFile(sourceFile, outputPath); err != nil {
		return fmt.Errorf("failed to copy config file: %w", err)
	}

	fmt.Printf("Example configuration file created at: %s\n", outputPath)
	return nil
}
