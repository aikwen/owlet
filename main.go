package main

import (
	"fmt"
	"os"

	"github.com/aikwen/owlet/internal/clipboard"
	"github.com/aikwen/owlet/internal/config"
	"github.com/aikwen/owlet/internal/snippet"
	"github.com/aikwen/owlet/internal/tui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run 执行 owlet 主流程。
func run() error {
	// 加载配置。
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// 确保运行目录存在。
	if err := cfg.EnsureDirs(); err != nil {
		return fmt.Errorf("create owlet dirs: %w", err)
	}

	// 加载 snippet 文件。
	snippets, err := snippet.LoadAll(cfg.SnippetDir)
	if err != nil {
		return fmt.Errorf("load snippets: %w", err)
	}

	if len(snippets) == 0 {
		return fmt.Errorf("no snippets found, create TOML files under %s", cfg.SnippetDir)
	}

	// 启动 TUI 并获取用户选择的命令。
	command, ok, err := tui.Run(snippets)
	if err != nil {
		return fmt.Errorf("run tui: %w", err)
	}

	if !ok || command == "" {
		return nil
	}

	writeCommand(command)

	return nil
}

// writeCommand 输出命令，并尝试写入系统剪贴板。
func writeCommand(command string) {
	if err := clipboard.WriteText(command); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to copy command: %v\n", err)
	}

	fmt.Println(command)
}