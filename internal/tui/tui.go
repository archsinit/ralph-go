package tui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Message represents a message to display in the transcript.
type Message struct {
	Author string
	Text   string
}

// StreamStart starts rendering a streamed message from the given author.
type StreamStart struct {
	Author string
}

// StreamToken adds a token to the in-progress streamed message.
type StreamToken struct {
	Text string
}

// StreamEnd finalizes the in-progress streamed message.
type StreamEnd struct{}

// msgSubmit is an internal message sent when the user presses Enter.
type msgSubmit struct {
	text string
}

// colorPalette is a list of distinct colors for authors.
var colorPalette = []string{"10", "12", "11", "9", "13", "14", "15"}

// getAuthorStyle returns a lipgloss style for the given author, assigning a new color on first appearance.
func (m *Model) getAuthorStyle(author string) lipgloss.Style {
	if style, ok := m.authorColors[author]; ok {
		return style
	}

	// Assign a new color from the palette
	colorIdx := len(m.authorList) % len(colorPalette)
	color := colorPalette[colorIdx]
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	m.authorColors[author] = style
	m.authorList = append(m.authorList, author)

	return style
}

// Model is the TUI model for the plan mode chatroom.
type Model struct {
	viewport viewport.Model
	input    textinput.Model
	spinner  spinner.Model
	program  *tea.Program

	width  int
	height int

	messages         []Message
	userScrolled     bool // Track if user has manually scrolled up
	lastViewSize     int  // Track content size to detect new messages
	inProgressAuthor string
	inProgressText   string
	isStreaming      bool
	onSubmit         func(text string) // Called when user presses Enter
	onCancel         func()            // Called when user presses ctrl+c
	authorColors     map[string]lipgloss.Style
	authorList       []string // Track order of first appearance for consistent coloring
}

// NewModel creates a new TUI model with optional initial messages and a submit callback.
func NewModel(opts ...Option) *Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your message..."
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	m := &Model{
		input:        ti,
		spinner:      s,
		messages:     []Message{},
		onSubmit:     func(string) {}, // Default no-op
		onCancel:     func() {},       // Default no-op
		authorColors: make(map[string]lipgloss.Style),
		authorList:   []string{},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Option is a functional option for configuring NewModel.
type Option func(*Model)

// WithSubmitCallback sets a callback function for when the user submits text.
func WithSubmitCallback(cb func(text string)) Option {
	return func(m *Model) {
		m.onSubmit = cb
	}
}

// WithCancelCallback sets a callback function for when the user presses ctrl+c during streaming.
func WithCancelCallback(cb func()) Option {
	return func(m *Model) {
		m.onCancel = cb
	}
}

// WithInitialMessages sets the initial transcript to display.
func WithInitialMessages(msgs []Message) Option {
	return func(m *Model) {
		m.messages = append([]Message{}, msgs...) // Copy to avoid external mutation
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

		// Clamp dimensions for small terminals (min 5 lines for viewport, 3 for status/input)
		viewportHeight := msg.Height - 3
		if viewportHeight < 5 {
			viewportHeight = 5
		}

		m.viewport.Width = msg.Width
		m.viewport.Height = viewportHeight

		m.input.Width = msg.Width

		// Re-wrap viewport content on resize
		m.updateViewport()

	case StreamStart:
		m.inProgressAuthor = msg.Author
		m.inProgressText = ""
		m.isStreaming = true
		m.userScrolled = false
		m.updateViewport()
		return m, m.spinner.Tick

	case StreamToken:
		m.inProgressText += msg.Text
		m.updateViewport()

	case StreamEnd:
		// StreamEnd only clears the in-progress stream; the engine will call AddMessage
		// to persist the message after handling it (for crash-safety and model ownership)
		m.inProgressAuthor = ""
		m.inProgressText = ""
		m.isStreaming = false
		m.updateViewport()

	case msgSubmit:
		// Don't add the message locally; the orchestrator is responsible.
		// Call the submit callback to notify the orchestrator.
		m.onSubmit(msg.text)
		m.input.Reset()
		m.userScrolled = false

	case msgAddMessage:
		m.messages = append(m.messages, Message{
			Author: msg.author,
			Text:   msg.text,
		})
		m.updateViewport()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// Signal cancel to orchestrator; it decides whether to quit or retry
			m.onCancel()
			return m, nil
		case "enter":
			if m.input.Value() != "" {
				return m, func() tea.Msg {
					return msgSubmit{text: m.input.Value()}
				}
			}
		case "up", "pgup", "down", "pgdn", "home", "end":
			// Forward scroll keys to viewport
			m.viewport, _ = m.viewport.Update(msg)
			// Track userScrolled based on actual viewport position
			m.userScrolled = !m.viewport.AtBottom()
			return m, nil
		}
	}

	// Handle spinner and input
	var cmd tea.Cmd
	if m.isStreaming {
		m.spinner, cmd = m.spinner.Update(msg)
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)

	if cmd != nil && inputCmd != nil {
		return m, tea.Batch(cmd, inputCmd)
	} else if cmd != nil {
		return m, cmd
	}
	return m, inputCmd
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
		authorStyle := m.getAuthorStyle(msg.Author)
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
		authorStyle := m.getAuthorStyle(m.inProgressAuthor)
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

// wrapText wraps text to the specified width, breaking long unbroken tokens if needed.
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
			// Word fits on current line
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		} else if currentLine.Len() == 0 {
			// Word doesn't fit and line is empty: break the word
			for len(word) > width {
				lines = append(lines, word[:width])
				word = word[width:]
			}
			currentLine.WriteString(word)
		} else {
			// Word doesn't fit but line has content: start new line
			lines = append(lines, currentLine.String())
			currentLine.Reset()

			// If word itself is too long, break it
			for len(word) > width {
				lines = append(lines, word[:width])
				word = word[width:]
			}
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// Handle is a handle to a running TUI program that allows external code to interact with it.
type Handle struct {
	program *tea.Program
	model   *Model
	mu      sync.RWMutex // Protects program and model access
	errCh   chan error
}

// StartStream notifies the TUI that a stream from the given author is starting.
func (h *Handle) StartStream(author string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.program != nil {
		h.program.Send(StreamStart{Author: author})
	}
}

// SendToken sends a token to the in-progress stream.
func (h *Handle) SendToken(text string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.program != nil {
		h.program.Send(StreamToken{Text: text})
	}
}

// EndStream finalizes the in-progress stream.
func (h *Handle) EndStream() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.program != nil {
		h.program.Send(StreamEnd{})
	}
}

// AddMessage appends a finalized message to the transcript.
// Author and Text are required; Text may be empty.
func (h *Handle) AddMessage(author, text string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.program != nil {
		h.program.Send(msgAddMessage{author: author, text: text})
	}
}

// Wait blocks until the TUI program exits and returns any error.
func (h *Handle) Wait() error {
	return <-h.errCh
}

// Quit sends a quit signal to the TUI program.
func (h *Handle) Quit() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.program != nil {
		h.program.Send(tea.Quit())
	}
}

// msgAddMessage is an internal message to add a message to the transcript.
type msgAddMessage struct {
	author string
	text   string
}

// Run starts the TUI with the given options and returns a Handle to control it immediately.
// The TUI runs in a background goroutine. Call Handle.Wait() to wait for it to exit.
func Run(opts ...Option) (*Handle, error) {
	m := NewModel(opts...)
	p := tea.NewProgram(m, tea.WithAltScreen())
	m.program = p
	h := &Handle{program: p, model: m, errCh: make(chan error, 1)}

	// Run the program in a goroutine so this returns immediately
	go func() {
		_, err := p.Run()
		h.mu.Lock()
		h.program = nil // Clear program to indicate it has exited
		h.mu.Unlock()
		h.errCh <- err
	}()

	return h, nil
}
