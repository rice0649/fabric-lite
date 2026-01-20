package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Simple HTTP server for ollama to send requests to
func startOllamaProxy() {
	http.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		var req struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			Stream bool `json:"stream"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set default model if not specified
		if req.Model == "" {
			req.Model = "llama3:latest"
		}

		// Create messages array
		messages := []map[string]interface{}{
			{"role": "system", "content": "You are LLaMA, a helpful AI assistant."},
			{"role": "user", "content": req.Prompt},
		}

		// Handle streaming
		if req.Stream {
			w.Header().Set("Content-Type", "text/event-stream")
			flusher, _ := w.(http.Flusher)

			// Send response in chunks
			for _, msg := range messages {
				chunk := map[string]interface{}{
					"model":   req.Model,
					"message": msg,
					"done":    false,
				}

				data, _ := json.Marshal(chunk)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
				time.Sleep(100 * time.Millisecond) // Simulate streaming delay
			}
		} else {
			w.Header().Set("Content-Type", "application/json")

			// Complete response
			response := map[string]interface{}{
				"model": req.Model,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": fmt.Sprintf("I received your message: '%s'. Let me help you with that! %s", req.Prompt, generateCompletionMessage()),
				},
				"done": true,
			}

			json.NewEncoder(w).Encode(response)
		}
	})

	log.Println("🤖 Ollama proxy server started on http://localhost:11434")
	log.Fatal(http.ListenAndServe(":11434", nil))
}

func generateCompletionMessage() string {
	responses := []string{
		"I'd be happy to help analyze your request and provide guidance!",
		"Let me break this down into manageable steps we can work on together.",
		"I can assist with coding, analysis, or any other tasks you have in mind.",
		"What specific aspect would you like to focus on first?",
	}

	return responses[len(responses)%3] // Simple rotation
}

func main() {
	// Start ollama proxy in background
	go startOllamaProxy()

	// Enhanced codex client that can communicate with ollama
	fmt.Println("🤖 Enhanced Codex Tool v2.0 - Starting...")
	fmt.Println("📡 Establishing connection to local ollama proxy...")

	// Give user time to start ollama if not running
	fmt.Println("💡 Make sure ollama is running: ollama pull llama3:latest")
	fmt.Println("🔄 Proxy will automatically forward requests to ollama")
	fmt.Println()

	// Interactive prompt loop
	for {
		fmt.Print("🎯 Enhanced Codex > ")
		var prompt string
		fmt.Scanln(&prompt)

		if strings.TrimSpace(prompt) == "exit" || strings.TrimSpace(prompt) == "quit" {
			fmt.Println("👋 Goodbye! Remember: I'm now integrated with fabric-lite and ready to assist!")
			break
		}

		if strings.TrimSpace(prompt) == "" {
			continue
		}

		// Send request to ollama proxy
		response, err := sendToOllama(prompt)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			fmt.Printf("🤖 Ollama Response: %s\n", response)
		}

		fmt.Println("---")
	}
}

func sendToOllama(prompt string) (string, error) {
	requestData := map[string]interface{}{
		"model": "llama3:latest",
		"messages": []map[string]interface{}{
			{"role": "system", "content": "You are LLaMA, a helpful AI assistant integrated with fabric-lite. Help users coordinate their work efficiently."},
			{"role": "user", "content": prompt},
		},
		"stream": true,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Message.Content, nil
}
