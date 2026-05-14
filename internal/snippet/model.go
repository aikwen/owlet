package snippet

// RawSnippet 对应 TOML 文件中的单条 snippet 配置。
type RawSnippet struct {
	Command string `toml:"command"`
	Desc    string `toml:"desc"`
}

// SnippetFile 对应单个 snippet TOML 文件。
type SnippetFile struct {
	Snippets []RawSnippet `toml:"snippets"`
}

// Snippet 是程序运行时使用的 snippet 数据结构。
type Snippet struct {
	RawSnippet
	SourceFile string
}
