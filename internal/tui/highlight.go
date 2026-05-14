package tui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var highlightStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("36"))

// highlight 渲染命中的搜索关键词。
func highlight(text string, query string) string {
	keywords := strings.Fields(strings.ToLower(strings.TrimSpace(query)))
	if len(keywords) == 0 || text == "" {
		return text
	}

	matches := findMatches(strings.ToLower(text), keywords)
	if len(matches) == 0 {
		return text
	}

	return renderMatches(text, matches)
}

type matchRange struct {
	start int
	end   int
}

func findMatches(text string, keywords []string) []matchRange {
	var matches []matchRange

	for _, keyword := range keywords {
		if keyword == "" {
			continue
		}

		offset := 0
		for {
			pos := strings.Index(text[offset:], keyword)
			if pos < 0 {
				break
			}

			start := offset + pos
			end := start + len(keyword)
			matches = append(matches, matchRange{start: start, end: end})
			offset = end
		}
	}

	return mergeMatches(matches)
}

func mergeMatches(matches []matchRange) []matchRange {
	if len(matches) <= 1 {
		return matches
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].start < matches[j].start
	})

	merged := []matchRange{matches[0]}
	for _, current := range matches[1:] {
		last := &merged[len(merged)-1]
		if current.start <= last.end {
			last.end = max(last.end, current.end)
			continue
		}

		merged = append(merged, current)
	}

	return merged
}

func renderMatches(text string, matches []matchRange) string {
	var b strings.Builder
	cursor := 0

	for _, m := range matches {
		if m.start < cursor {
			continue
		}

		b.WriteString(text[cursor:m.start])
		b.WriteString(highlightStyle.Render(text[m.start:m.end]))
		cursor = m.end
	}

	b.WriteString(text[cursor:])

	return b.String()
}

// truncateHighlighted 按终端显示宽度截断文本，并保留命中高亮。
func truncateHighlighted(text string, query string, width int) string {
	truncated := truncate(text, width)
	return highlight(truncated, query)
}
