package snippet

import (
	"strings"
)

// Search 根据 query 过滤 snippets。
//
// 搜索规则：
//   - query 为空时返回全部 snippets；
//   - 大小写不敏感；
//   - query 会按空白字符拆成多个关键词；
//   - 每个关键词都必须命中 command、desc 或 sourceFile 中的任意一个。
func Search(snippets []Snippet, query string) []Snippet {
	keywords := parseKeywords(query)
	if len(keywords) == 0 {
		return snippets
	}

	results := make([]Snippet, 0, len(snippets))
	for _, item := range snippets {
		if matchSnippet(item, keywords) {
			results = append(results, item)
		}
	}

	return results
}

// parseKeywords 将用户输入拆成搜索关键词。
func parseKeywords(query string) []string {
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(query)))
	if len(fields) == 0 {
		return nil
	}

	return fields
}

// matchSnippet 判断 snippet 是否匹配所有关键词。
func matchSnippet(item Snippet, keywords []string) bool {
	text := strings.ToLower(strings.Join([]string{
		item.Command,
		item.Desc,
		item.SourceFile,
	}, " "))

	for _, keyword := range keywords {
		if !strings.Contains(text, keyword) {
			return false
		}
	}

	return true
}