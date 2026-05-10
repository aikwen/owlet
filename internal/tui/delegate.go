package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	itemPrefix         = "  "
	selectedItemPrefix = "> "
)

var (
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().Bold(true)
)

// delegate 渲染 snippet 列表项。
type delegate struct{}

func newDelegate() delegate {
	return delegate{}
}

// Height 返回列表项渲染高度。
func (d delegate) Height() int {
	return 1
}

// Spacing 返回列表项之间的间距。
func (d delegate) Spacing() int {
	return 0
}

// Update 处理列表项渲染器消息。
func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render 渲染单个列表项。
func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(item)
	if !ok {
		return
	}

	selected := index == m.Index()
	width := m.Width()

	fmt.Fprint(w, d.renderCommand(it, selected, width))
}

func (d delegate) renderCommand(it item, selected bool, width int) string {
	prefix := itemPrefix
	style := itemStyle

	if selected {
		prefix = selectedItemPrefix
		style = selectedItemStyle
	}

	contentWidth := max(width-lipgloss.Width(prefix), 0)
	command := truncateHighlighted(it.snippet.Command, it.query, contentWidth)

	return style.Render(prefix) + command
}

// truncate 按终端显示宽度截断文本。
func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}

	if lipgloss.Width(s) <= width {
		return s
	}

	ellipsis := "…"
	ellipsisWidth := lipgloss.Width(ellipsis)

	if width <= ellipsisWidth {
		return ellipsis
	}

	var b strings.Builder
	currentWidth := 0

	for _, r := range s {
		runeWidth := lipgloss.Width(string(r))

		if currentWidth+runeWidth+ellipsisWidth > width {
			break
		}

		b.WriteRune(r)
		currentWidth += runeWidth
	}

	b.WriteString(ellipsis)
	return b.String()
}
