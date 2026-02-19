# Build Your Own Claude Code (Go)

A small CLI agent written in Go that can call LLMs through OpenRouter and use local tools (like reading files) in an agent loop.

The program sends a prompt to an LLM, allows the model to call tools, executes the tool locally, and feeds the result back to the model until a final answer is produced.

---

## Requirements

- Go 1.21+ (or newer)
- An OpenRouter API key

---

## Environment Setup (.env)

Create a `.env` file in the root of the project.

```
OPENROUTER_API_KEY=your_openrouter_api_key_here
OPENROUTER_BASE_URL=https://openrouter.ai/api/v1
```

---

## Install Dependencies

From the project root:

```bash
go mod tidy
```

## Run the Agent
```
./your_program -p "Determine in how many months the chemical expires by reading README.md. Respond with only a number."
```

["Build Your own Claude Code" Challenge](https://codecrafters.io/challenges/claude-code).

Along the way you'll learn about HTTP RESTful APIs, OpenAI-compatible tool
calling, agent loop, and how to integrate multiple tools into an AI assistant.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.