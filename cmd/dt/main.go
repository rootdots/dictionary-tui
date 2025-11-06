package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Version information (populated by goreleaser)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok && version == "dev" {
		if info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
	}
}

// --- API Structures ---
type Phonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type Definition struct {
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type WordEntry struct {
	Word       string     `json:"word"`
	Phonetic   string     `json:"phonetic"`
	Phonetics  []Phonetic `json:"phonetics"`
	Meanings   []Meaning  `json:"meanings"`
	SourceUrls []string   `json:"sourceUrls"`
}

// --- History Item (for bubbles/list) ---
type historyItem string

func (i historyItem) FilterValue() string { return string(i) }
func (i historyItem) Title() string       { return string(i) }
func (i historyItem) Description() string { return "" }

// --- Styles and Constants ---
const (
	maxHistory   = 10
	defaultWidth = 80
)

var (
	// Common styles for both TUI and CLI
	keywordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	defStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	panelBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("51"))

	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Part of speech style
	posStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Bold(true).
			Underline(true)

	// Example style
	exampleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Italic(true)

	// Phonetic style
	phoneticStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Italic(true)

	// Panel style for CLI mode
	cliPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("51")).
			Padding(1, 2).
			Width(defaultWidth)

	// Help styles
	helpHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1).
			PaddingBottom(1)

	helpSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("213")).
				MarginBottom(1).
				MarginTop(1)

	helpTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	helpCommandStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Bold(true)

	helpExampleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("246")).
				Italic(true)

	// Title style for CLI mode
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)
)

// Define the two possible modes for the application
type appMode int

const (
	searchMode appMode = iota
	historyMode
)

type model struct {
	// Large fields first for better memory alignment
	textInput  textinput.Model
	viewport   viewport.Model
	history    list.Model
	definition string
	word       string
	err        error
	mode       appMode
	ready      bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter word..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	vp := viewport.New(defaultWidth, 20)
	vp.SetContent("Welcome! Search for a word to see its definition here.")

	h := list.New([]list.Item{}, list.NewDefaultDelegate(), defaultWidth, 20)
	h.Title = "Search History (Ctrl+H to switch)"
	h.SetShowFilter(false)

	return model{
		mode:       searchMode,
		textInput:  ti,
		viewport:   vp,
		history:    h,
		definition: "Welcome! Search for a word to see its definition here.",
	}
}

// --- Messages for Asynchronous Operations ---
type resultMsg string
type errMsg error

// --- Init, Update, View ---

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit, true

	case tea.KeyCtrlH:
		if m.mode == searchMode {
			m.mode = historyMode
			m.textInput.Blur()
		} else {
			m.mode = searchMode
			m.textInput.Focus()
		}
		return m, nil, true

	case tea.KeyEsc:
		if m.mode == historyMode {
			m.mode = searchMode
			m.textInput.Focus()
			return m, nil, true
		}
	}
	return m, nil, false
}

func (m model) handleSearchMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
		word := strings.TrimSpace(m.textInput.Value())
		if word != "" {
			m.word = word
			m.textInput.SetValue("")
			m.textInput.Blur()
			m.addToHistory(word)
			m.definition = "Searching..."
			m.viewport.SetContent(m.definition)
			m.viewport.GotoTop()
			return m, lookupDefinitionCmd(m.word)
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) handleHistoryMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
		selectedItem := m.history.SelectedItem()
		if selectedItem != nil {
			selectedWord := string(selectedItem.(historyItem))
			m.mode = searchMode
			m.word = selectedWord
			m.textInput.Focus()
			m.definition = fmt.Sprintf("Reviewing definition for: %s...", keywordStyle.Render(selectedWord))
			m.viewport.SetContent(m.definition)
			m.viewport.GotoTop()
			return m, lookupDefinitionCmd(selectedWord)
		}
	}
	m.history, cmd = m.history.Update(msg)
	return m, cmd
}

func (m model) handleWindowSize(msg tea.WindowSizeMsg) model {
	hPad := appStyle.GetHorizontalPadding() * 2
	vPad := appStyle.GetVerticalPadding() * 2
	borderPad := 2

	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())

	availableWidth := msg.Width - hPad - borderPad
	availableHeight := msg.Height - headerHeight - footerHeight - vPad - borderPad

	if !m.ready {
		m.viewport = viewport.New(availableWidth, availableHeight)
		m.history.SetSize(availableWidth, availableHeight)
		m.ready = true
	} else {
		m.viewport.Width = availableWidth
		m.viewport.Height = availableHeight
		m.history.SetSize(availableWidth, availableHeight)
	}

	m.textInput.Width = m.viewport.Width / 3
	m.viewport.SetContent(m.definition)
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		var cmd tea.Cmd
		var handled bool
		var newModel tea.Model
		newModel, cmd, handled = m.handleKeyMsg(msg)
		m = newModel.(model)
		if handled {
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m = m.handleWindowSize(msg)

	case resultMsg:
		m.definition = string(msg)
		m.viewport.SetContent(m.definition)
		m.textInput.Focus()

	case errMsg:
		m.err = msg
		m.definition = fmt.Sprintf("Error fetching definition: %s\nPress Enter to try again.", msg.Error())
		m.viewport.SetContent(m.definition)
		m.textInput.Focus()
	}

	if m.mode == searchMode {
		return m.handleSearchMode(msg)
	}
	return m.handleHistoryMode(msg)
}

func (m model) View() string {
	header := m.headerView()

	var content string
	if m.mode == searchMode {
		content = panelBorder.
			Width(m.viewport.Width + 2).
			Height(m.viewport.Height + 2).
			Render(m.viewport.View())
	} else {
		content = panelBorder.
			Width(m.history.Width() + 2).
			Height(m.history.Height() + 2).
			Render(m.history.View())
	}

	footer := m.footerView()

	return appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	))
}

// --- Helper Methods ---

func (m *model) addToHistory(word string) {
	newItems := []list.Item{historyItem(word)}
	newItems = append(newItems, m.history.Items()...)
	if len(newItems) > maxHistory {
		newItems = newItems[:maxHistory]
	}
	m.history.SetItems(newItems)
}

func (m model) headerView() string {
	var status string
	if m.word != "" {
		status = fmt.Sprintf("Last Search: %s", keywordStyle.Render(m.word))
	} else {
		status = "Dictionary TUI"
	}

	inputWidth := lipgloss.Width(m.textInput.View())
	remainingWidth := m.viewport.Width - inputWidth

	statusText := lipgloss.PlaceHorizontal(
		remainingWidth,
		lipgloss.Right,
		status,
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.textInput.View(),
		statusText,
	)
}

func (m model) footerView() string {
	percent := int(m.viewport.ScrollPercent() * 100)
	info := fmt.Sprintf("%d%%", percent)

	var modeHelp string
	if m.mode == searchMode {
		modeHelp = "(Ctrl+H for History)"
	} else {
		modeHelp = "(Esc to Search)"
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		fmt.Sprintf("Press Ctrl+C to quit %s", modeHelp),
		lipgloss.PlaceHorizontal(m.viewport.Width-48, lipgloss.Right, info),
	)
}

// --- Dictionary API Integration ---

func lookupWord(word string) ([]WordEntry, error) {
	cleanWord := strings.TrimSpace(strings.ToLower(word))
	if cleanWord == "" {
		return nil, fmt.Errorf("please enter a word to search")
	}

	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", cleanWord)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("definition for '%s' not found", word)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var entries []WordEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no data found for '%s'", word)
	}

	return entries, nil
}

func lookupDefinitionCmd(word string) tea.Cmd {
	return func() tea.Msg {
		entries, err := lookupWord(word)
		if err != nil {
			return errMsg(err)
		}
		result := formatDefinition(entries[0])
		return resultMsg(result)
	}
}

func formatDefinition(entry WordEntry) string {
	var builder strings.Builder

	builder.WriteString(keywordStyle.Render(entry.Word))
	if entry.Phonetic != "" {
		builder.WriteString(" " + phoneticStyle.Render(entry.Phonetic))
	}
	builder.WriteString("\n\n")

	for i, meaning := range entry.Meanings {
		builder.WriteString(posStyle.Render(fmt.Sprintf("%d. %s", i+1, meaning.PartOfSpeech)) + "\n")

		for j, def := range meaning.Definitions {
			builder.WriteString(fmt.Sprintf("   %s %s\n", defStyle.Render(fmt.Sprintf("%d.", j+1)), def.Definition))
			if def.Example != "" {
				builder.WriteString("      " + exampleStyle.Render("• "+def.Example) + "\n")
			}
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func formatHelp() string {
	var b strings.Builder

	b.WriteString(helpHeaderStyle.Render(fmt.Sprintf("Dictionary-TUI %s", version)) + "\n")
	b.WriteString(helpTextStyle.Render("A dictionary application with interactive TUI and CLI interfaces.\n"))
	b.WriteString(helpSectionStyle.Render("USAGE: ") + helpTextStyle.Render("dt [FLAGS] [WORD]") + "\n")
	b.WriteString(helpTextStyle.Render("       dt [WORD]") + "\n")
	b.WriteString(helpSectionStyle.Render("FLAGS: ") + fmt.Sprintf("%s  %s\n",
		helpCommandStyle.Render("-w, --word"),
		helpTextStyle.Render("Specify a word to look up"),
	))
	b.WriteString("       " + fmt.Sprintf("%s  %s\n",
		helpCommandStyle.Render("-h, --help"),
		helpTextStyle.Render("Show this help message"),
	))
	b.WriteString("       " + fmt.Sprintf("%s  %s\n",
		helpCommandStyle.Render("--version"),
		helpTextStyle.Render("Show version information"),
	))
	b.WriteString(helpSectionStyle.Render("MODES: ") + helpTextStyle.Render("1. Interactive Mode (TUI):") + "\n")
	b.WriteString(helpTextStyle.Render("       Launch without arguments to enter the interactive interface.") + "\n")
	b.WriteString(helpTextStyle.Render("       • Use Ctrl+H to access search history") + "\n")
	b.WriteString(helpTextStyle.Render("       • Use Ctrl+C to quit") + "\n")
	b.WriteString(helpTextStyle.Render("       2. Command-Line Mode:") + "\n")
	b.WriteString(helpTextStyle.Render("       Provide a word as an argument for quick definition lookup.") + "\n")
	b.WriteString(helpSectionStyle.Render("EXAMPLES: ") + helpCommandStyle.Render("dt") + "\n")
	b.WriteString(helpExampleStyle.Render("          # Launch interactive mode") + "\n")
	b.WriteString(helpCommandStyle.Render("          dt serendipity") + "\n")
	b.WriteString(helpExampleStyle.Render("          # Look up 'serendipity' directly") + "\n")
	b.WriteString(helpCommandStyle.Render("          dt -w ephemeral") + "\n")
	b.WriteString(helpExampleStyle.Render("          # Look up 'ephemeral' using flag syntax"))

	return cliPanelStyle.Render(strings.TrimRight(b.String(), "\n"))
}

func main() {
	// Version flag
	versionFlag := flag.Bool("version", false, "Print version information")

	flag.Usage = func() {
		fmt.Println(formatHelp())
		os.Exit(0)
	}
	wordPtr := flag.String("w", "", "Word to look up")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("dictionary-tui version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", date)
		return
	}

	if *wordPtr != "" || len(flag.Args()) > 0 {
		var word string
		if *wordPtr != "" {
			word = *wordPtr
		} else {
			word = flag.Args()[0]
		}

		entries, err := lookupWord(word)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		title := titleStyle.Render(fmt.Sprintf("Dictionary-TUI: %s", word))
		fmt.Printf("\n%s\n\n%s\n", title, cliPanelStyle.Render(formatDefinition(entries[0])))
		return
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
