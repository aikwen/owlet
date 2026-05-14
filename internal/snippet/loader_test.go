package snippet

import (
	"path/filepath"
	"testing"

	"github.com/aikwen/owlet/internal/config"
)

func TestLoadExamples(t *testing.T) {
	cfg, err := config.LoadFromFile(filepath.Join("..", "..", "examples", "config.toml"))
	if err != nil {
		t.Fatal(err)
	}

	snippets, err := LoadAll(cfg.SnippetDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(snippets) == 0 {
		t.Fatal("no snippets loaded from examples")
	}
}
