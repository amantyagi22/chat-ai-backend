# Claude Clone Backend

A robust Go backend service that powers the Claude Clone chat interface. This service handles chat interactions and integrates with Hugging Face's model API to provide AI responses.

## Tech Stack

- **Language**: [Go 1.21+](https://golang.org/) - For high performance and strong typing
- **Framework**: [Fiber](https://gofiber.io/) - Express-inspired web framework
- **AI Integration**: [Hugging Face API](https://huggingface.co/) - For accessing AI models
- **Environment**: [Docker](https://www.docker.com/) - For containerization
- **Deployment**: [Render](https://render.com/) - For cloud hosting

## Features

- 🚀 High-performance Go backend
- 🔒 Secure API endpoints
- 🤖 Hugging Face model integration
- 🔄 Real-time chat processing
- 🎯 RESTful API design
- 📦 Docker containerization
- 🛡️ CORS support
- 🔍 Request validation
- 📝 Structured logging

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerization)
- Hugging Face API key

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory:
```env
PORT=8080
HUGGINGFACE_API_KEY=your_api_key_here
MODEL_ID=facebook/blenderbot-400M-distill
```

4. Run the server:
```bash
go run main.go
```

The server will start on http://localhost:8080

### Docker Setup

1. Build the Docker image:
```bash
docker build -t claude-clone-backend .
```

2. Run the container:
```bash
docker run -p 8080:8080 --env-file .env claude-clone-backend
```

## API Endpoints

### POST /api/chat
Processes chat messages and returns AI responses.

Request:
```json
{
  "message": "Hello, how are you?",
  "conversation_id": "optional-conversation-id"
}
```

Response:
```json
{
  "response": "I'm doing well, thank you for asking! How can I help you today?",
  "conversation_id": "generated-conversation-id"
}
```

## Project Structure

```
backend/
├── main.go                 # Entry point
├── handlers/              # Request handlers
├── models/               # Data models
├── config/               # Configuration
├── middleware/           # Custom middleware
├── utils/               # Utility functions
├── Dockerfile           # Docker configuration
├── go.mod               # Go modules file
└── go.sum               # Go modules checksum
```

## Environment Variables

- `PORT`: Server port (default: 8080)
- `HUGGINGFACE_API_KEY`: Your Hugging Face API key
- `MODEL_ID`: Hugging Face model ID to use
- `ALLOWED_ORIGINS`: CORS allowed origins (comma-separated)

## Deployment

This project is configured for deployment on Render:

1. Push your code to GitHub
2. Create a new Web Service on Render
3. Connect your repository
4. Set environment variables
5. Deploy!

## Error Handling

The API uses standard HTTP status codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 500: Internal Server Error

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with Go's excellent standard library
- Uses Fiber for high-performance routing
- Integrates with Hugging Face's powerful AI models 