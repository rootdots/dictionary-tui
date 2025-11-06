# Gemini Code Assistant Context

## Project Overview

This project is a terminal-based dictionary application written in Go. It provides a simple and efficient way to look up word definitions directly from the command line.

The application features a clean, text-based user interface (TUI) built using the Bubble Tea framework, allowing for an interactive and user-friendly experience.

### Core Functionalities:

-   **Word Lookup:** Users can enter a word to fetch its definition.
-   **Definition Display:** The application displays the word, its phonetic transcription, and a list of meanings, including part of speech, definitions, and example sentences.
-   **Search History:** The application maintains a history of the last 10 searched words, which can be accessed by pressing `Ctrl+H`.
-   **Interactive UI:** The TUI allows for easy navigation and interaction, with clear instructions on how to use the application.

### Technical Details:

-   **Language:** Go
-   **Framework:** Bubble Tea for the TUI
-   **API:** The application uses the free Dictionary API (`https://dictionaryapi.dev/`) to fetch word definitions.
-   **Dependencies:** The project's dependencies are managed using Go Modules and are listed in the `go.mod` file. Key libraries include:
    -   `github.com/charmbracelet/bubbletea`
    -   `github.com/charmbracelet/bubbles`
    -   `github.com/charmbracelet/lipgloss`

## Building and Running

To build and run the project, you need to have Go installed on your system.

### Build

To build the application, run the following command in the project's root directory:

```sh
go build
```

This will create a binary executable named `dictionary-tui` in the project directory.

### Run

To run the application directly without building a separate binary, use the following command:

```sh
go run .
```

This will compile and run the application.

## Development Conventions

The project follows standard Go development practices and leverages the Model-View-Update (MVU) architecture provided by the Bubble Tea framework.

### Code Structure

-   **`main.go`**: This is the main entry point of the application. It contains the entire logic for the TUI, including the model, view, and update functions, as well as the API integration.
-   **`go.mod` and `go.sum`**: These files manage the project's dependencies.

### Coding Style

The code is well-structured and includes comments to explain the different parts of the application. It uses the `lipgloss` library for styling the TUI, with a clear separation of styles and components.

### Error Handling

The application includes error handling for API requests and displays informative messages to the user in case of an error.
