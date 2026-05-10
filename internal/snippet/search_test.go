package snippet

import (
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	snippets := []Snippet{
		{
			RawSnippet: RawSnippet{
				Command: "conda env list",
				Desc:    "查看所有 conda 环境",
			},
			SourceFile: "conda.toml",
		},
		{
			RawSnippet: RawSnippet{
				Command: "conda activate myenv",
				Desc:    "激活指定 conda 环境",
			},
			SourceFile: "conda.toml",
		},
		{
			RawSnippet: RawSnippet{
				Command: "git status",
				Desc:    "查看当前工作区状态",
			},
			SourceFile: "git.toml",
		},
		{
			RawSnippet: RawSnippet{
				Command: "docker ps",
				Desc:    "查看正在运行的容器",
			},
			SourceFile: "docker.toml",
		},
	}

	tests := []struct {
		name     string
		query    string
		expected []Snippet
	}{
		{
			name:     "empty query returns all snippets",
			query:    "",
			expected: snippets,
		},
		{
			name:  "match command",
			query: "git",
			expected: []Snippet{
				snippets[2],
			},
		},
		{
			name:  "match desc",
			query: "容器",
			expected: []Snippet{
				snippets[3],
			},
		},
		{
			name:  "match source file",
			query: "docker",
			expected: []Snippet{
				snippets[3],
			},
		},
		{
			name:  "match multiple keywords with and logic",
			query: "conda env",
			expected: []Snippet{
				snippets[0],
				snippets[1],
			},
		},
		{
			name:  "case insensitive",
			query: "GIT",
			expected: []Snippet{
				snippets[2],
			},
		},
		{
			name:     "no match",
			query:    "kubectl",
			expected: []Snippet{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Search(snippets, tt.query)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("Search() = %#v, expected %#v", got, tt.expected)
			}
		})
	}
}