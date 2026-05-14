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
}

// New 创建 TUI model。
func New(snippets []snippet.Snippet) Model {
	// input 组件
	input := textinput.New()
	input.Placeholder = "Search snippets..."
	input.Focus()
	input.CharLimit = 256
	input.Width = defaultWidth
	input.Prompt = "> "
	input.PromptStyle = promptStyle

	items := toItems(snippets, "")

	// list 组件
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
	b.WriteString("\n\n")
	b.WriteString(m.list.View())

	if m.expanded {
		b.WriteString("\n\n")
		b.WriteString(m.renderDetail())
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓ move · tab detail · enter select · esc quit"))

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
	count := len(m.list.Items())

	contentHeight := min(count, defaultListHeight)
	if contentHeight <= 0 {
		contentHeight = 1
	}

	return contentHeight + 1
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
