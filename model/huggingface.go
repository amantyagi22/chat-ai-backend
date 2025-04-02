package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type HuggingFaceRequest struct {
	Inputs      string                 `json:"inputs"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// BlenderBot specific response format
type HuggingFaceResponse struct {
	GeneratedText string   `json:"generated_text,omitempty"`
	Conversation struct {
		Generated []string `json:"generated_responses,omitempty"`
		Past     []string `json:"past_user_inputs,omitempty"`
	} `json:"conversation,omitempty"`
}

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

func GetModelResponse(input string) (string, error) {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		response, err := makeRequest(input)
		if err != nil {
			log.Printf("Attempt %d failed: %v", i+1, err)
			if i < maxRetries-1 && isRetryableError(err) {
				delay := time.Second * time.Duration(2<<i)
				log.Printf("Waiting %v before next attempt...", delay)
				time.Sleep(delay)
				continue
			}
			return "", err
		}
		return response, nil
	}
	return "", fmt.Errorf("max retries exceeded after %d attempts, model may be temporarily unavailable", maxRetries)
}

func makeRequest(input string) (string, error) {
	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("HUGGINGFACE_API_KEY environment variable not set")
	}

	modelID := os.Getenv("MODEL_ID")
	if modelID == "" {
		modelID = "facebook/blenderbot-400M-distill"
	}

	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", modelID)
	log.Printf("Making request to model: %s", modelID)

	requestBody := HuggingFaceRequest{
		Inputs: input,
		Parameters: map[string]interface{}{
			"max_length": 128,
			"temperature": 0.8,
			"top_p": 0.95,
			"repetition_penalty": 1.2,
			"num_return_sequences": 1,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		return "", fmt.Errorf("service unavailable (503): model is loading")
	}

	if resp.StatusCode != http.StatusOK {
		// Try to parse error response
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != nil {
			return "", fmt.Errorf("API error: %v", errorResp.Error)
		}
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	log.Printf("Raw response: %s", string(body))

	// Try to parse as BlenderBot response first
	var blenderResponse HuggingFaceResponse
	if err := json.Unmarshal(body, &blenderResponse); err == nil {
		if len(blenderResponse.Conversation.Generated) > 0 {
			return blenderResponse.Conversation.Generated[0], nil
		}
		if blenderResponse.GeneratedText != "" {
			return blenderResponse.GeneratedText, nil
		}
	}

	// Fallback to array response format
	var arrayResponse []struct {
		GeneratedText string `json:"generated_text"`
	}
	if err := json.Unmarshal(body, &arrayResponse); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v, raw response: %s", err, string(body))
	}

	if len(arrayResponse) == 0 {
		return "", fmt.Errorf("empty response from model")
	}

	return arrayResponse[0].GeneratedText, nil
}

func isRetryableError(err error) bool {
	errStr := err.Error()
	return errStr == "service unavailable (503): model is loading" ||
		errStr == "error making request: context deadline exceeded" ||
		errStr == "empty response from model"
}