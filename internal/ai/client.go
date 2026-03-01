package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const groqURL = "https://api.groq.com/openai/v1/chat/completions"
const model = "llama-3.3-70b-versatile"
const visionModel = "meta-llama/llama-4-scout-17b-16e-instruct"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		apiKey: os.Getenv("GROQ_API_KEY"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) IsConfigured() bool {
	return c.apiKey != ""
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"` // string or []contentPart
}

type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL string `json:"url"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *Client) GenerateResponse(prompt string) (string, error) {
	reqBody := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{
				Role: "system",
				Content: `You are FitTalk AI — a knowledgeable fitness and health assistant on a forum.
You provide helpful, evidence-based advice about training, nutrition, health, recovery, and motivation.
Always respond in the same language as the user's post.
Keep responses concise (2-4 paragraphs). Use markdown formatting for readability.
If the topic is medical, remind the user to consult a professional.
IMPORTANT: Never mix languages or scripts. If the post is in Russian, use only Cyrillic characters. Never insert Chinese, Japanese, or other non-relevant characters.`,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponse marshal: %w", err)
	}

	req, err := http.NewRequest("POST", groqURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponse request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponse do: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponse read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ai.Client.GenerateResponse status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponse unmarshal: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("ai.Client.GenerateResponse: empty choices")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// GenerateResponseWithImage sends a prompt with an image to the vision model.
// imageBase64 should be the full data URL, e.g. "data:image/jpeg;base64,..."
func (c *Client) GenerateResponseWithImage(prompt string, imageBase64 string) (string, error) {
	userContent := []contentPart{
		{Type: "text", Text: prompt},
		{Type: "image_url", ImageURL: &imageURL{URL: imageBase64}},
	}

	reqBody := chatRequest{
		Model: visionModel,
		Messages: []chatMessage{
			{
				Role: "system",
				Content: `You are FitTalk AI — a knowledgeable fitness and health assistant on a forum.
You provide helpful, evidence-based advice about training, nutrition, health, recovery, and motivation.
Always respond in the same language as the user's post.
Keep responses concise (2-4 paragraphs). Use markdown formatting for readability.
If the topic is medical, remind the user to consult a professional.
When an image is attached, analyze its contents and incorporate your observations into the response.
IMPORTANT: Never mix languages or scripts. If the post is in Russian, use only Cyrillic characters. Never insert Chinese, Japanese, or other non-relevant characters.`,
			},
			{
				Role:    "user",
				Content: userContent,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponseWithImage marshal: %w", err)
	}

	req, err := http.NewRequest("POST", groqURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponseWithImage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponseWithImage do: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponseWithImage read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ai.Client.GenerateResponseWithImage status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("ai.Client.GenerateResponseWithImage unmarshal: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("ai.Client.GenerateResponseWithImage: empty choices")
	}

	return chatResp.Choices[0].Message.Content, nil
}
