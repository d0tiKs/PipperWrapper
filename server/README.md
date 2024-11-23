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
  "text": "Hello, world!"
}
```

## Testing
You can test the server using `websocat`:
```bash
echo '{"text": "Hello world"}' | websocat -t ws://localhost:8080/tts
```

## Dependencies
- [github.com/gorilla/websocket](https://github.com/gorilla/websocket)

## Limitations
- Requires a local Piper TTS installation
- Currently supports a single voice model

## ToDo
- [ ] Implement configuration file support
- [ ] Use `gorilla/mux` instead of `net/http`
- [ ] Generate Swagger documentation with `summerfish-swagger`
- [ ] Choose a license

## License
[Specify your license]
