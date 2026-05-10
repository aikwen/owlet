package tui

import "github.com/aikwen/owlet/internal/snippet"

// item 适配 Bubbles list.Item。
type item struct {
	snippet snippet.Snippet
	query   string
}

func newItem(s snippet.Snippet, query string) item {
	return item{
		snippet: s,
		query:   query,
	}
}

// FilterValue 返回列表过滤文本。
func (i item) FilterValue() string {
	return i.snippet.Command
}
