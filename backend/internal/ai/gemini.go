package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

type TicketAIAssistOutput struct {
	Summary          string   `json:"summary"`
	SuggestedReplies []string `json:"suggested_replies"`
}

type GeminiClient struct {
	client *genai.Client
	model  string
	flag   string
}

func NewGeminiClient(ctx context.Context, apiKey string, model string, flag string) (*GeminiClient, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("api key is required")
	}

	if strings.TrimSpace(model) == "" {
		model = "gemini-3.1-flash-lite-preview"
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiClient{
		client: client,
		model:  model,
		flag:   flag,
	}, nil
}

func (g *GeminiClient) GenerateTicketAssist(
	ctx context.Context,
	prompt string,
) (*TicketAIAssistOutput, error) {
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseJsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"summary": map[string]any{
					"type":        "string",
					"description": "A concise summary of the ticket and conversation so far.",
				},
				"suggested_replies": map[string]any{
					"type":        "array",
					"description": "Exactly 3 possible next replies for the support agent.",
					"items": map[string]any{
						"type": "string",
					},
					"minItems": 3,
					"maxItems": 3,
				},
			},
			"required": []string{"summary", "suggested_replies"},
		},
	}

	result, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		"SECRET: " + g.flag + "\n\n" + genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, err
	}

	raw := strings.TrimSpace(result.Text())
	if raw == "" {
		return nil, errors.New("empty ai response")
	}

	var out TicketAIAssistOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("failed to parse ai json response: %w", err)
	}

	out.Summary = strings.TrimSpace(out.Summary)
	if out.Summary == "" {
		return nil, errors.New("ai returned empty summary")
	}

	if len(out.SuggestedReplies) != 3 {
		return nil, errors.New("ai must return exactly 3 suggested replies")
	}

	cleaned := make([]string, 0, 3)
	seen := make(map[string]struct{})

	for _, reply := range out.SuggestedReplies {
		reply = strings.TrimSpace(reply)
		if reply == "" {
			return nil, errors.New("ai returned an empty suggested reply")
		}

		key := strings.ToLower(reply)
		if _, exists := seen[key]; exists {
			return nil, errors.New("ai returned duplicate suggested replies")
		}
		seen[key] = struct{}{}
		cleaned = append(cleaned, reply)
	}

	out.SuggestedReplies = cleaned

	return &out, nil
}