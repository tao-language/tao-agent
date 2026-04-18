package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"tao-agent/internal/provider"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UI interface for the engine to interact with.
type UI interface {
	Print(message string)
	Ask(question string) string
	PromptStream(chunks <-chan provider.Chunk, errs <-chan error) string
}

// BubbleTeaUI implements the UI interface using charmbracelet/bubbletea.
// Note: For MVP, we'll keep it simple by wrapping the tea.Program.
type BubbleTeaUI struct {
	program *tea.Program
}

func NewBubbleTeaUI() *BubbleTeaUI {
	m := initialModel()
	p := tea.NewProgram(m)
	return &BubbleTeaUI{program: p}
}

// Styles
var (
	thinkingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true)
	contentStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	systemStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
)

type model struct {
	history  []string
	current  string
	thinking string
	input    textinput.Model
	waiting  bool
	err      error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	return model{
		history: []string{},
		input:   ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var s strings.Builder
	for _, h := range m.history {
		s.WriteString(h + "\n")
	}
	s.WriteString(m.input.View())
	return s.String()
}

// These methods will actually run the tea.Program in a way that blocks or handles messages.
// This is a bit complex with BubbleTea because it's usually the driver.
// We'll implement a simpler version for the MVP that uses BubbleTea for the TUI display.

func (b *BubbleTeaUI) Print(message string) {
	fmt.Println(systemStyle.Render("Tao: ") + message)
}

func (b *BubbleTeaUI) Ask(question string) string {
	fmt.Print(systemStyle.Render("Tao: ") + question + " ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (b *BubbleTeaUI) PromptStream(chunks <-chan provider.Chunk, errs <-chan error) string {
	var fullContent strings.Builder
	var fullThinking strings.Builder

	fmt.Print(systemStyle.Render("Agent: "))

	for {
		select {
		case err := <-errs:
			if err != nil {
				fmt.Printf("\nError: %v\n", err)
				return ""
			}
		case chunk, ok := <-chunks:
			if !ok {
				fmt.Println()
				return fullContent.String()
			}

			if chunk.Thinking != "" {
				// Print thinking in gray
				fmt.Print(thinkingStyle.Render(chunk.Thinking))
				fullThinking.WriteString(chunk.Thinking)
			}
			if chunk.Content != "" {
				// Print content in normal color
				fmt.Print(contentStyle.Render(chunk.Content))
				fullContent.WriteString(chunk.Content)
			}

			if chunk.Done {
				fmt.Println()
				return fullContent.String()
			}
		}
	}
}
