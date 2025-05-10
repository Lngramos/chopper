# Chopper

Chopper is an experimental, powerful, and extensible command-line assistant powered by local LLMs via [Ollama](https://ollama.com). It supports structured tool calls to interact with your system.

> 🧠 Chopper seeks to be a local-first alternative to tools like **Claude Code CLI** and **Amazon Q CLI**, but with all language model execution done on your machine using **self-hosted models via Ollama**.

⚠️ **Warning**: This tool is experimental and potentially dangerous. Tool execution is not sandboxed or restricted — use with caution.

## Prerequisites

- [Ollama](https://ollama.com) must be installed and running locally.
- At least one model (e.g. `gemma3`, `qwen3`, `deepseek-r1`, etc.) should be pulled and available via `ollama list`.

## Usage

### Chat (one-shot prompt)

```bash
$ go run main.go chat -m gemma3 -t 0.7 'Run the command `whoami` using the run tool'
luis
```

### REPL (interactive mode)

```bash
$ go run main.go repl

>> run the command `whoami` using the run tool
luis

>> read the contents of /etc/hosts using the read_file tool
127.0.0.1 localhost

```

## Supported Tools

- `run`: Execute a shell command
- `read_file`: Read a file from disk


## Roadmap

Chopper is evolving from a structured LLM wrapper into a more intelligent, agent-like assistant. Planned features:

### Completed
- Interactive REPL with message history
- One-shot chat command
- Basic tool calling via JSON (`tool_call`)
- Modular, extensible tool interface
- Tools: `run`, `read_file`

### Coming Soon
- Tool: `write_file`
- File summarisation, code linting
- Session persistence to disk
- Config file support (`~/.config/chopper`)
- More advanced prompt management
- Multiple tool call execution
- Safe mode / sandboxing
- Agent planning capabilities

## License

MIT
