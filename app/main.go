package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/samiyonas/build-your-own-claude-code/app/llmtools"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type ToolFunc func(args string) (string, error)

var toolFuncMap = map[string]ToolFunc{
	"Read":  llmtools.ReadTool,
	"Write": llmtools.WriteTool,
	"Bash":  llmtools.BashTool,
}

func main() {

	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String(prompt),
				},
			},
		},
	}
	tools := []openai.ChatCompletionToolUnionParam{
		{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name: "Read",
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]any{
							"file_path": map[string]any{
								"type":        "string",
								"description": "The path to the file to read. The file is guaranteed to be less than 100KB in size.",
								"examples":    []string{"/tmp/file.txt"},
							},
						},
						"required": []string{"file_path"},
					},
				},
			},
		},
		{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        "Write",
					Description: openai.String("Write content to a file"),
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]any{
							"file_path": map[string]any{
								"type":        "string",
								"description": "The path to the file to write.",
							},
							"content": map[string]any{
								"type":        "string",
								"description": "The content to write to the file. The content will be less than 100KB in size.",
							},
						},
						"required": []string{"file_path", "content"},
					},
				},
			},
		},
		{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        "Bash",
					Description: openai.String("Execute a bash command and return the output. The command will not have access to any files and will be guaranteed to finish executing within 10 seconds."),
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]any{
							"command": map[string]any{
								"type":        "string",
								"description": "The command to execute",
							},
						},
						"required": []string{"command"},
					},
				},
			},
		},
	}

	for {
		resp, err := client.Chat.Completions.New(context.Background(),
			openai.ChatCompletionNewParams{
				Model:    "anthropic/claude-haiku-4.5",
				Messages: messages,
				Tools:    tools,
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(resp.Choices) == 0 {
			return
		}

		messg := resp.Choices[0].Message
		if len(messg.ToolCalls) == 0 {
			fmt.Println(messg.Content)
			return
		}

		for _, toolCall := range messg.ToolCalls {
			args := toolCall.Function.Arguments
			name := toolCall.Function.Name

			toolFunc, ok := toolFuncMap[name]
			if !ok {
				fmt.Fprintf(os.Stderr, "error: unknown tool %s\n", name)
				return
			}

			content, err := toolFunc(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error executing tool: %v\n", err)
				return
			}

			messages = append(messages,
				openai.ChatCompletionMessageParamUnion{
					OfAssistant: &openai.ChatCompletionAssistantMessageParam{
						ToolCalls: []openai.ChatCompletionMessageToolCallUnionParam{
							{
								OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
									ID: toolCall.ID,
									Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
										Name:      toolCall.Function.Name,
										Arguments: toolCall.Function.Arguments,
									},
								},
							},
						},
					},
				},
				openai.ChatCompletionMessageParamUnion{
					OfTool: &openai.ChatCompletionToolMessageParam{
						Role:       "tool",
						ToolCallID: toolCall.ID,
						Content: openai.ChatCompletionToolMessageParamContentUnion{
							OfString: openai.String(content),
						},
					},
				},
			)
		}
	}
}

/*
{
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": null,
        "tool_calls": [
          {
            "id": "call_abc123",
            "type": "function",
            "function": {
              "name": "Read",
              "arguments": "{\"file_path\": \"/path/to/file.txt\"}"
            }
          }
        ]
      },
      "finish_reason": "tool_calls"
    }
  ]
}
*/
