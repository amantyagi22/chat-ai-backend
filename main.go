package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"claude-clone/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, you should configure this properly
	},
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Message string `json:"message"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/api/chat", handleChat).Methods("POST")
	r.HandleFunc("/ws", handleWebSocket)

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://claude-clone-frontend.vercel.app", // Add your Vercel domain
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	var message ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get response from the model
	modelResponse, err := model.GetModelResponse(message.Content)
	if err != nil {
		log.Printf("Error getting model response: %v", err)
		http.Error(w, "Error getting model response", http.StatusInternalServerError)
		return
	}

	response := ChatResponse{
		Message: modelResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Get response from the model
		modelResponse, err := model.GetModelResponse(string(p))
		if err != nil {
			log.Printf("Error getting model response: %v", err)
			if err := conn.WriteMessage(messageType, []byte("Error getting model response")); err != nil {
				log.Println(err)
				return
			}
			continue
		}

		if err := conn.WriteMessage(messageType, []byte(modelResponse)); err != nil {
			log.Println(err)
			return
		}
	}
} 