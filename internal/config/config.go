package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	DefaultAppDirName     = ".owlet"
	DefaultConfigFileName = "config.toml"
	DefaultSnippetDirName = "snippets"
	DefaultShortcut       = ""
)

type Config struct {
	AppDir     string
	ConfigFile string
	SnippetDir string
	Shortcut   string
}

// fileConfig 对应 config.toml 中允许用户配置的字段。
type fileConfig struct {
	Shortcut string `toml:"shortcut"`
}

// Load 使用固定默认路径加载配置。
//
// 默认配置文件路径为：
//
//	~/.owlet/config.toml
func Load() (Config, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	configFile := filepath.Join(
		userHome,
		DefaultAppDirName,
		DefaultConfigFileName,
	)

	return LoadFromFile(configFile)
}

// LoadFromFile 从指定 config.toml 路径加载配置。
//
// AppDir 会使用 configFile 所在目录。
// SnippetDir 固定为 AppDir/snippets。
// 如果 configFile 不存在，则返回默认配置。
func LoadFromFile(configFile string) (Config, error) {
	configFile = filepath.Clean(configFile)
	appDir := filepath.Dir(configFile)

	cfg := Config{
		AppDir:     appDir,
		ConfigFile: configFile,
		SnippetDir: filepath.Join(appDir, DefaultSnippetDirName),
		Shortcut:   DefaultShortcut,
	}

	raw, err := os.ReadFile(configFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}

		return Config{}, err
	}

	var fc fileConfig
	if err := toml.Unmarshal(raw, &fc); err != nil {
		return Config{}, err
	}

	cfg.Shortcut = fc.Shortcut

	return cfg, nil
}

// EnsureDirs 创建 owlet 运行所需目录。
func (c Config) EnsureDirs() error {
	return os.MkdirAll(c.SnippetDir, 0o755)
}
