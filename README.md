# Chopper

Chopper is an experimental command-line tool that uses local LLMs (via [Ollama](https://ollama.com)) to help with coding and system tasks.

⚠️ **Warning**: This tool is experimental and potentially dangerous. Tool execution is not sandboxed or restricted — use with caution.

## Prerequisites

- [Ollama](https://ollama.com) must be installed and running locally.
- At least one model (e.g. `qwen3:14b`, `mistral`) should be pulled and available via `ollama list`.

## Usage

### Chat (one-shot prompt)

```bash
$ chopper chat -m qwen3:14b 'Run the command `whoami` using the run tool'
luis
```

### REPL (interactive mode)

```
$ chopper repl -m qwen3:14b
>> run the command `whoami` using the run tool
luis

>> read the contents of /etc/hosts using the read_file tool
127.0.0.1 localhost
```

## Supported Tools

- `run`: Execute a shell command
- `read_file`: Read a file from disk
