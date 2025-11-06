# dt 📚

A terminal-based dictionary with both TUI and CLI modes.

## 🚀 Installation

### Binary Release
Download the latest release from [GitHub Releases](https://github.com/rootdots/dictionary-tui/releases)

```bash
# Example for Linux x86_64
curl -L -o dt.tar.gz https://github.com/rootdots/dictionary-tui/releases/download/vX.Y.Z/dictionary-tui_Linux_x86_64.tar.gz
tar xzf dt.tar.gz
chmod +x dt
```

### Docker 🐳
```bash
docker run --rm ghcr.io/rootdots/dt:latest -w hello
```

### From Source
```bash
go install github.com/rootdots/dictionary-tui/cmd/dt@latest
```

## 📖 Usage

### Interactive Mode (TUI)
```bash
dt
```
- Type a word and press Enter to search
- `Ctrl+H` to view search history
- `Ctrl+C` to exit

### CLI Mode
```bash
dt -w hello     # Look up a specific word
dt --version    # Show version
```

## ✨ Features
- 🖥️ Interactive TUI with search history
- 📝 Direct CLI word lookup
- 🎨 Styled output in both modes
- 🌐 Uses [Free Dictionary API](https://dictionaryapi.dev/)

## 📜 License
MIT