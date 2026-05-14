package tui

import (
	"strings"

	"github.com/aikwen/owlet/internal/snippet"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultWidth      = 80
	defaultHeight     = 24
	defaultListHeight = 10
)

var (
	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33"))
	helpStyle        = lipgloss.NewStyle().Faint(true)
	detailLabelStyle = lipgloss.NewStyle().Bold(true)
	detailTextStyle  = lipgloss.NewStyle().PaddingLeft(2)
)

// Model 管理 owlet 的 TUI 状态。
type Model struct {
	input    textinput.Model
	list     list.Model
	snippets []snippet.Snippet

	query    string
	expanded bool
	selected string
	ok       bool
	width    int
	height   int
}

// New 创建 TUI model。
func New(snippets []snippet.Snippet) Model {
	input := textinput.New()
	input.Placeholder = "Search snippets..."
	input.Focus()
	input.CharLimit = 256
	input.Width = defaultWidth
	input.Prompt = "> "
	input.PromptStyle = promptStyle

	items := toItems(snippets, "")

	l := list.New(items, newDelegate(), defaultWidth, defaultListHeight)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowPagination(true)
	l.SetFilteringEnabled(false)

	m := Model{
		input:    input,
		list:     l,
		snippets: snippets,
		width:    defaultWidth,
		height:   defaultHeight,
	}
	m.resizeList()

	return m
}

// Init 初始化 TUI 命令。
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update 处理 TUI 消息并更新状态。
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width
		m.resizeList()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if selected, ok := m.selectedItem(); ok {
				m.selected = selected.snippet.Command
				m.ok = true
				return m, tea.Quit
			}

			return m, nil

		case "tab":
			if _, ok := m.selectedItem(); !ok {
				return m, nil
			}

			m.expanded = !m.expanded
			m.resizeList()
			return m, nil
		}
	}

	oldQuery := m.input.Value()

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	if m.input.Value() != oldQuery {
		m.query = m.input.Value()
		m.expanded = false

		items := toItems(snippet.Search(m.snippets, m.query), m.query)
		cmds = append(cmds, m.list.SetItems(items))
		m.list.Select(0)
		m.resizeList()
	}

	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}

// View 渲染 TUI 视图。
func (m Model) View() string {
	var b strings.Builder

	b.WriteString(m.input.View())

	listHeight := m.listHeight()
	detailHeight := m.rawDetailHeight()
	detailVisible := m.showDetail(listHeight, detailHeight)
	helpVisible := m.showHelp(listHeight, detailVisible, detailHeight)

	if listHeight > 0 {
		if m.showListGap() {
			b.WriteString("\n\n")
		} else {
			b.WriteString("\n")
		}

		b.WriteString(m.list.View())
	}

	if detailVisible {
		b.WriteString("\n\n")
		b.WriteString(m.renderDetail())
	}

	if helpVisible {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("↑/↓ move · tab detail · enter select · esc quit"))
	}

	return b.String()
}

// Run 启动 TUI 并返回选中的命令。
func Run(snippets []snippet.Snippet) (string, bool, error) {
	program := tea.NewProgram(New(snippets))

	finalModel, err := program.Run()
	if err != nil {
		return "", false, err
	}

	m, ok := finalModel.(Model)
	if !ok {
		return "", false, nil
	}

	return m.selected, m.ok, nil
}

// selectedItem 返回当前选中的列表项。
func (m Model) selectedItem() (item, bool) {
	selected := m.list.SelectedItem()
	if selected == nil {
		return item{}, false
	}

	it, ok := selected.(item)
	return it, ok
}

// resizeList 调整列表显示尺寸。
func (m *Model) resizeList() {
	m.list.SetSize(m.width, m.listHeight())
}

// listHeight 返回列表显示高度。
func (m Model) listHeight() int {
	if m.height <= 0 {
		return m.defaultTotalListHeight()
	}

	availableHeight := m.height - 1
	if availableHeight <= 0 {
		return 0
	}

	availableHeight -= m.listGapLines() - 1
	if availableHeight <= 0 {
		return 0
	}

	return min(m.defaultTotalListHeight(), availableHeight)
}

// defaultTotalListHeight 返回默认列表总高度，包括分页区域。
func (m Model) defaultTotalListHeight() int {
	return m.defaultContentListHeight() + 1
}

// defaultContentListHeight 返回默认列表内容高度。
func (m Model) defaultContentListHeight() int {
	count := len(m.list.Items())

	contentHeight := min(count, defaultListHeight)
	if contentHeight <= 0 {
		return 1
	}

	return contentHeight
}

// showListGap 返回是否在输入框和列表之间保留空行。
func (m Model) showListGap() bool {
	return m.height <= 0 || m.height >= 4
}

// listGapLines 返回输入框和列表之间的换行数。
func (m Model) listGapLines() int {
	if m.showListGap() {
		return 2
	}

	return 1
}

// appendedHeight 返回追加内容占用的显示高度。
func appendedHeight(newlines int, contentHeight int) int {
	if contentHeight <= 0 {
		return 0
	}

	return newlines + contentHeight - 1
}

// baseViewHeight 返回输入框和列表占用的显示高度。
func (m Model) baseViewHeight(listHeight int) int {
	height := 1
	if listHeight > 0 {
		height += appendedHeight(m.listGapLines(), listHeight)
	}

	return height
}

// showDetail 返回是否渲染详情区域。
func (m Model) showDetail(listHeight int, detailHeight int) bool {
	if !m.expanded {
		return false
	}

	if detailHeight <= 0 {
		return false
	}

	if m.height <= 0 {
		return true
	}

	return m.baseViewHeight(listHeight)+appendedHeight(2, detailHeight) <= m.height
}

// showHelp 返回是否渲染快捷键提示。
func (m Model) showHelp(listHeight int, detailVisible bool, detailHeight int) bool {
	if m.height <= 0 {
		return true
	}

	usedHeight := m.baseViewHeight(listHeight)
	if detailVisible {
		usedHeight += appendedHeight(2, detailHeight)
	}

	return usedHeight+appendedHeight(1, 1) <= m.height
}

// rawDetailHeight 返回详情区域原始显示高度。
func (m Model) rawDetailHeight() int {
	detail := m.renderDetail()
	if strings.TrimSpace(detail) == "" {
		return 0
	}

	return lipgloss.Height(detail)
}

// renderDetail 渲染当前选中项详情。
func (m Model) renderDetail() string {
	it, ok := m.selectedItem()
	if !ok {
		return ""
	}

	var b strings.Builder

	writeDetailBlock(&b, "Source", it.snippet.SourceFile, m.width)
	writeDetailBlock(&b, "Command", it.snippet.Command, m.width)

	if strings.TrimSpace(it.snippet.Desc) != "" {
		writeDetailBlock(&b, "Desc", it.snippet.Desc, m.width)
	}

	return b.String()
}

// toItems 将 snippets 转换为 Bubbles list items。
func toItems(snippets []snippet.Snippet, query string) []list.Item {
	items := make([]list.Item, 0, len(snippets))
	for _, s := range snippets {
		items = append(items, newItem(s, query))
	}

	return items
}

// writeDetailBlock 渲染详情字段。
func writeDetailBlock(b *strings.Builder, label string, value string, width int) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}

	if b.Len() > 0 {
		b.WriteString("\n")
	}

	b.WriteString(detailLabelStyle.Render(label + ":"))
	b.WriteString("\n")

	contentWidth := max(width-2, 0)
	text := truncate(value, contentWidth)
	b.WriteString(detailTextStyle.Render(text))
}
