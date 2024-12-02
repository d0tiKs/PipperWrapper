# Piper TTS WebSocket Server

## Overview
A WebSocket server that converts text to speech using Piper TTS, allowing remote TTS requests via the WebSocket protocol.

## Prerequisites
- Go 1.22.9+
- Piper TTS installed and available in the path
- aplay (for audio playback)

## Installation
To set up the server, run the following commands:
```bash
go mod tidy
go build -o piper-websocket-tts-server
```

## Configuration
Modify server settings in the source code:
- WebSocket port (default: 8080)
- Piper voice model path

## Usage
### Start Server
To start the server, execute:
```bash
./piper-websocket-tts-server
```

### Client Request Format
Send JSON requests in the following format:
```json
{
  "text": "Hello, world!",
  "language": "en"
}
```

## Testing
You can test the server using `websocat`:
```bash
echo '{"text": "Hello world", "language" : "en" }' | websocat -t ws://localhost:8080/tts
```

## Dependencies
- [github.com/rhasspy/piper](https://github.com/rhasspy/piper)
- [github.com/gorilla/websocket](https://github.com/gorilla/websocket)

## Limitations
- Requires a local Piper TTS installation

## In Progress
- Send voice over the websocket
- Command Line arguments + usage
- Config file

## ToDo
- [ ] Implement configuration file support
- [ ] Use `gorilla/mux` instead of `net/http`
- [ ] Generate Swagger documentation with `summerfish-swagger`
- [ ] Choose a license
- [ ] Install script that setup env and user
- [ ] Add service user to service
- [ ] Detect the text's language automaticaly



## License
[Specify your license]
