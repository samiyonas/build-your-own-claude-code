package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func Read(filePath string) (string, error) {
	filePath = strings.TrimPrefix(filePath, "/")

	clean := filepath.Clean(filePath)
	if strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid file path")
	}

	info, err := os.Stat(clean)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "", fmt.Errorf("file path is a directory")
	}

	if strings.Contains(filePath, "..") {
		return "", fmt.Errorf("invalid file path")
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func Write(filePath string, content string) error {
	filePath = strings.TrimPrefix(filePath, "/")

	clean := filepath.Clean(filePath)
	if strings.HasPrefix(clean, "..") {
		return fmt.Errorf("invalid file path")
	}

	if strings.Contains(filePath, "..") {
		return fmt.Errorf("invalid file path")
	}

	err := os.WriteFile(clean, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadTool(args string) (string, error) {
	var readArgs ReadArgs
	err := json.Unmarshal([]byte(args), &readArgs)
	if err != nil {
		return "", fmt.Errorf("error parsing arguments: %v", err)
	}

	content, err := Read(readArgs.FilePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return content, nil
}

func WriteTool(args string) (string, error) {
	var writeArgs WriteArgs
	err := json.Unmarshal([]byte(args), &writeArgs)
	if err != nil {
		return "", fmt.Errorf("error parsing arguments: %v", err)
	}

	err = Write(writeArgs.FilePath, writeArgs.Content)
	if err != nil {
		return "", fmt.Errorf("error writing file: %v", err)
	}

	return "File written successfully", nil
}

type ToolFunc func(args string) (string, error)

var toolFuncMap = map[string]ToolFunc{
	"Read":  ReadTool,
	"Write": WriteTool,
}

/*
func main() {
	type ReadArgs struct {
		FilePath string `json:"file_path"`
	}

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
	resp, err := client.Chat.Completions.New(context.Background(),
		openai.ChatCompletionNewParams{
			Model: "anthropic/claude-haiku-4.5",
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String(prompt),
						},
					},
				},
			},
			Tools: []openai.ChatCompletionToolUnionParam{
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
			},
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if len(resp.Choices) == 0 {
		panic("No choices in response")
	}

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	messg := resp.Choices[0].Message
	if len(messg.ToolCalls) == 0 {
		fmt.Println(messg.Content)
		return
	}
	//name := resp.Choices[0].Message.ToolCalls[0].Function.Name

	args := messg.ToolCalls[0].Function.Arguments

	var readArgs ReadArgs

	err = json.Unmarshal([]byte(args), &readArgs)
	if err != nil {
		log.Fatal("Error parsing arguments: ", err)
	}

	content, err := Read(readArgs.FilePath)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}

	fmt.Println(content)
	//fmt.Print(resp.Choices[0].Message.Content)
}
*/

type ReadArgs struct {
	FilePath string `json:"file_path"`
}

type WriteArgs struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
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
