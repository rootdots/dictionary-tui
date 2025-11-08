# Dictionary TUI 📚

A terminal-based dictionary application written in Go. It provides a simple and efficient way to look up word definitions directly from the command line.

The application features a clean, text-based user interface (TUI) built using the Bubble Tea framework, allowing for an interactive and user-friendly experience.

## ✨ Features

-   **Word Lookup:** Users can enter a word to fetch its definition.
-   **Definition Display:** The application displays the word, its phonetic transcription, and a list of meanings, including part of speech, definitions, and example sentences.
-   **Search History:** The application maintains a history of the last 10 searched words, which can be accessed by pressing `Ctrl+H`.
-   **Interactive UI:** The TUI allows for easy navigation and interaction, with clear instructions on how to use the application.
-   **CLI Mode:** The application can be used as a command-line tool to look up words directly from the terminal.

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

To start the interactive TUI mode, run the application without any arguments:

```bash
dt
```

-   Type a word and press `Enter` to search for its definition.
-   Press `Ctrl+H` to access your search history.
-   Press `Ctrl+C` to exit the application.

### Command-Line Mode (CLI)

You can also use the application as a standard command-line tool for quick lookups:

```bash
dt [word]
```

**Example:**

```bash
dt ephemeral
```

This will display the definition for the word "ephemeral" directly in your terminal.

## 🛠️ Building from Source

To build and run the project from source, you need to have Go installed on your system.

### Build

To build the application, run the following command in the project's root directory:

```sh
go build ./cmd/dt
```

This will create a binary executable named `dt` in the project directory.

### Run

To run the application directly without building a separate binary, use the following command:

```sh
go run ./cmd/dt
```

This will compile and run the application.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.