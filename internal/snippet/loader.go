package snippet

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// LoadAll 从指定目录加载所有 snippet TOML 文件。
//
// 加载规则：
//   - 只读取 *.toml 文件；
//   - 按文件名排序，保证加载顺序稳定；
//   - 每条 snippet 必须包含 command；
//   - desc 可以为空；
//   - SourceFile 使用 TOML 文件名，例如 conda.toml。
func LoadAll(snippetDir string) ([]Snippet, error) {
	paths, err := filepath.Glob(filepath.Join(snippetDir, "*.toml"))
	if err != nil {
		return nil, fmt.Errorf("glob snippet files: %w", err)
	}

	sort.Strings(paths)

	snippets := make([]Snippet, 0)

	for _, path := range paths {
		loaded, err := loadFile(path)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, loaded...)
	}

	return snippets, nil
}

// loadFile 加载单个 snippet TOML 文件。
func loadFile(path string) ([]Snippet, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read snippet file %q: %w", path, err)
	}

	var file SnippetFile
	if err := toml.Unmarshal(raw, &file); err != nil {
		return nil, fmt.Errorf("parse snippet file %q: %w", path, err)
	}

	sourceFile := filepath.Base(path)

	snippets := make([]Snippet, 0, len(file.Snippets))
	for i, rawSnippet := range file.Snippets {
		command := strings.TrimSpace(rawSnippet.Command)
		if command == "" {
			return nil, fmt.Errorf("invalid snippet file %q: snippets[%d].command is required", path, i)
		}

		snippets = append(snippets, Snippet{
			RawSnippet: RawSnippet{
				Command: command,
				Desc:    strings.TrimSpace(rawSnippet.Desc),
			},
			SourceFile: sourceFile,
		})
	}

	return snippets, nil
}