# Piper TTS WebSocket Server

## Overview
A WebSocket server that converts text to speech using Piper TTS, allowing remote TTS requests via WebSocket protocol.

## Prerequisites
- Go 1.22.9+
- Piper TTS installed and available in the path
- aplay (for audio playback)

## Installation
```bash
go mod tidy
go build -o piper-websocket-tts-server
```

## Configuration
Modify server settings in source code:
- WebSocket port (default: 8080)
- Piper voice model path

## Usage
### Start Server
```bash
./piper-websocket-tts-server
```

### Client Request Format
Send JSON:
```json
{
  "text": "Hello, world!",
}
```

## Testing
Use `websocat`:
```bash
echo '{"text": "Hello world"}' | websocat -t ws://localhost:8080/tts
```

## Dependencies
- github.com/gorilla/websocket

## Limitations
- Requires local Piper TTS installation
- Single voice model currently supported

## ToDo
[] Using a configuration file instead of static path
[] use `gorilla/mux` instead of `net/http`
- [] generate swagger documentation with `summerfish-swager`
[] Choose a license 

## License
[Specify your license]
