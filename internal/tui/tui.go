package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// renderedMsg represents a message ready for display in the viewport.
type renderedMsg struct {
	Author string
	Text   string
}

// Stream message types for feeding tokens into the TUI.
type msgStreamStart struct {
	author string
}

type msgStreamToken struct {
	text string
}

type msgStreamEnd struct {}

// authorColors maps author names to lipgloss styles.
var authorColors = map[string]lipgloss.Style{
	"user":  lipgloss.NewStyle().Foreground(lipgloss.Color("10")),  // Green
	"claude": lipgloss.NewStyle().Foreground(lipgloss.Color("12")), // Blue
	"codex": lipgloss.NewStyle().Foreground(lipgloss.Color("11")),  // Cyan
}

// getAuthorStyle returns a lipgloss style for the given author.
func getAuthorStyle(author string) lipgloss.Style {
	if style, ok := authorColors[author]; ok {
		return style
	}
	// Default to yellow for unknown authors.
	return lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
}

// Model is the TUI model for the plan mode chatroom.
type Model struct {
	viewport viewport.Model
	input    textinput.Model
	spinner  spinner.Model
	program  *tea.Program

	width  int
	height int

	messages         []renderedMsg
	userScrolled     bool // Track if user has manually scrolled up
	lastViewSize     int  // Track content size to detect new messages
	inProgressAuthor string
	inProgressText   string
	isStreaming      bool
}

// NewModel creates a new TUI model.
func NewModel() *Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your message..."
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return &Model{
		input:    ti,
		spinner:  s,
		messages: []renderedMsg{},
	}
}

// Init initializes the model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages from the Bubble Tea runtime.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 3 // Reserve space for input

		m.input.Width = msg.Width

		// Re-wrap viewport content on resize
		m.updateViewport()

	case msgStreamStart:
		m.inProgressAuthor = msg.author
		m.inProgressText = ""
		m.isStreaming = true
		m.userScrolled = false
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(nil)
		m.updateViewport()
		return m, cmd

	case msgStreamToken:
		m.inProgressText += msg.text
		m.updateViewport()

	case msgStreamEnd:
		if m.isStreaming && m.inProgressAuthor != "" {
			m.messages = append(m.messages, renderedMsg{
				Author: m.inProgressAuthor,
				Text:   m.inProgressText,
			})
			m.inProgressAuthor = ""
			m.inProgressText = ""
			m.isStreaming = false
			m.updateViewport()
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.input.Value() != "" {
				m.messages = append(m.messages, renderedMsg{
					Author: "user",
					Text:   m.input.Value(),
				})
				m.input.Reset()
				m.userScrolled = false
				m.updateViewport()
			}
		case "up", "pgup":
			// Detect when user scrolls up
			m.userScrolled = true
		case "down", "pgdn":
			// Check if we're near the bottom
			if m.viewport.AtBottom() {
				m.userScrolled = false
			}
		}
	}

	// Handle spinner ticks during streaming
	if m.isStreaming {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if cmd != nil {
			var inputCmd tea.Cmd
			m.input, inputCmd = m.input.Update(msg)
			return m, tea.Batch(cmd, inputCmd)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the model to the screen.
func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	viewportContent := m.viewport.View()
	inputContent := m.input.View()

	// Build status line
	var statusLine string
	if m.isStreaming {
		statusLine = fmt.Sprintf("%s waiting for %s...", m.spinner.View(), m.inProgressAuthor)
	} else {
		statusLine = "your turn"
	}
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	statusContent := statusStyle.Render(statusLine)

	// Render viewport on top, status, input below
	return lipgloss.JoinVertical(
		lipgloss.Top,
		viewportContent,
		strings.Repeat("─", m.width),
		statusContent,
		inputContent,
	)
}

// updateViewport refreshes the viewport content from messages.
func (m *Model) updateViewport() {
	var content strings.Builder
	for _, msg := range m.messages {
		authorStyle := getAuthorStyle(msg.Author)
		label := authorStyle.Render(msg.Author + ":")

		// Wrap text to viewport width, accounting for author label
		lines := wrapText(msg.Text, m.viewport.Width-4)
		if len(lines) == 0 {
			lines = []string{""}
		}

		// First line with author label
		content.WriteString(fmt.Sprintf("%s %s\n", label, lines[0]))

		// Subsequent lines indented
		for _, line := range lines[1:] {
			content.WriteString("  " + line + "\n")
		}
		content.WriteString("\n")
	}

	// Append in-progress message if streaming
	if m.isStreaming && m.inProgressAuthor != "" {
		authorStyle := getAuthorStyle(m.inProgressAuthor)
		label := authorStyle.Render(m.inProgressAuthor + ":")

		lines := wrapText(m.inProgressText, m.viewport.Width-4)
		if len(lines) == 0 {
			lines = []string{""}
		}

		// First line with author label
		content.WriteString(fmt.Sprintf("%s %s", label, lines[0]))
		if m.inProgressText != "" {
			content.WriteString(" ▌") // Cursor indicator
		}
		content.WriteString("\n")

		// Subsequent lines indented
		for _, line := range lines[1:] {
			content.WriteString("  " + line + "\n")
		}
	}

	contentStr := content.String()
	m.viewport.SetContent(contentStr)

	// Auto-scroll to bottom if user hasn't scrolled up manually
	if !m.userScrolled {
		m.viewport.GotoBottom()
	}

	m.lastViewSize = len(contentStr)
}

// wrapText wraps text to the specified width.
func wrapText(text string, width int) []string {
	if width <= 0 {
		width = 80
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return lines
	}

	var currentLine strings.Builder
	for _, word := range words {
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= width {
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		} else {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
			}
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// Run starts the TUI and blocks until the user quits.
func Run() error {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	m.program = p
	_, err := p.Run()
	return err
}

// StreamTokens feeds tokens from a channel into the model.
// Call this in a goroutine alongside the model's Run method.
func (m *Model) StreamTokens(author string, tokenChan <-chan string) {
	if m.program == nil {
		return
	}

	m.program.Send(msgStreamStart{author: author})

	for token := range tokenChan {
		m.program.Send(msgStreamToken{text: token})
	}

	m.program.Send(msgStreamEnd{})
}
