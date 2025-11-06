# Dictionary-TUI Context Documentation

## Project Overview

Dictionary-TUI is a terminal-based dictionary application written in Go that provides both an interactive TUI (Terminal User Interface) and a command-line interface for looking up word definitions.

### Key Features

- Interactive TUI mode with search history
- Direct CLI mode for quick lookups
- Rich text formatting using Charmbracelet components
- Clean and intuitive user interface
- Error handling and graceful fallbacks

## Usage Modes

### 1. Interactive TUI Mode

Launch the application without arguments to enter the interactive interface:
```sh
dt
```

Features:
- Real-time word lookup
- Search history (Ctrl+H)
- Easy navigation
- Scrollable content
- Clean exit (Ctrl+C)

### 2. Command-Line Mode

For quick word lookups directly from the command line:
```sh
dt serendipity
# or
dt -w serendipity
```

## Technical Stack

### Core Dependencies

- [Bubble Tea](github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](github.com/charmbracelet/lipgloss) - Style definitions
- [Free Dictionary API](https://dictionaryapi.dev/) - Word definitions

### Key Components

1. **TUI Components**
   - TextInput for word entry
   - Viewport for definition display
   - List for search history
   - Custom styled panels

2. **API Integration**
   - RESTful client for dictionary API
   - JSON parsing for responses
   - Error handling

3. **Styling**
   - Custom color schemes
   - Consistent formatting
   - Border designs
   - Typography enhancements

## Code Structure

### Main Components

1. **Data Structures**
   ```go
   type WordEntry struct {
       Word       string     `json:"word"`
       Phonetic   string     `json:"phonetic"`
       Phonetics  []Phonetic `json:"phonetics"`
       Meanings   []Meaning  `json:"meanings"`
       SourceUrls []string   `json:"sourceUrls"`
   }
   ```

2. **TUI Model**
   ```go
   type model struct {
       mode       appMode
       textInput  textinput.Model
       viewport   viewport.Model
       history    list.Model
       err        error
       ready      bool
       word       string
       definition string
   }
   ```

3. **Style Definitions**
   ```go
   var (
       keywordStyle = lipgloss.NewStyle().
           Foreground(lipgloss.Color("205")).
           Bold(true)
       
       panelBorder = lipgloss.NewStyle().
           Border(lipgloss.RoundedBorder()).
           BorderForeground(lipgloss.Color("51"))
   )
   ```

## Error Handling

1. **API Errors**
   - Network connectivity issues
   - Word not found
   - Invalid responses

2. **User Input**
   - Empty searches
   - Invalid characters
   - Command-line argument parsing

## Display Formatting

### TUI Mode
- Word and phonetics at the top
- Part of speech in bold magenta
- Definitions numbered and indented
- Examples in italic grey
- Search history in a separate view

### CLI Mode
- Clean panel design
- Consistent spacing
- Color-coded output
- Error messages in panels

## Future Improvements

1. **Features**
   - Audio pronunciation playback
   - Offline dictionary support
   - Synonym lookup
   - Word of the day

2. **Technical**
   - Caching layer
   - Configuration file support
   - Custom theme support
   - Plugin system

3. **User Experience**
   - Fuzzy search
   - Autocomplete
   - Multiple dictionary sources
   - Export functionality

## Building and Distribution

### Build
```sh
go build -o dt
```

### Install
```sh
go install
```

### Dependencies
```sh
go mod download
```

## Contributing

Contributions are welcome! Please ensure:
1. Code follows Go standards
2. Tests are included
3. Documentation is updated
4. Commit messages are clear

## License

MIT License - Feel free to use and modify as needed.