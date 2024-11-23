# Piper TTS WebSocket System

A complete text-to-speech system combining a WebSocket server powered by Piper TTS and a Chrome extension for easy text selection and playback.

## Components

### 1. WebSocket Server
A Go-based server that handles text-to-speech conversion using Piper TTS. [Learn more](server/README.md)

Key features:
- WebSocket communication
- Piper TTS integration
- Real-time audio processing
- Easy deployment

### 2. Chrome Extension
A browser extension that allows sending selected text to the TTS server. [Learn more](extension/README.md)

Key features:
- Right-click text selection
- Configurable server settings
- Debug mode
- Simple installation

## Quick Start

### Server Setup
```bash
cd server
go mod tidy
go build -o piper-websocket-tts-server
./piper-websocket-tts-server
```

### Extension Setup
1. Open Chrome and go to `chrome://extensions/`
2. Enable "Developer mode"
3. Click "Load unpacked" and select the `extension` directory
4. Configure server address and port in extension settings

## System Requirements

### Server
- Go 1.22.9+
- Piper TTS
- aplay (for audio playback)

### Extension
- Google Chrome browser
- WebSocket-compatible environment

## Architecture
```
┌─────────────────┐      WebSocket       ┌──────────────┐
│ Chrome Extension│ ─────────────────────>│   Go Server  │
└─────────────────┘    {"text": "..."}   └──────────────┘
                                               │
                                               │ Piper TTS
                                               ▼
                                         Audio Playback
```

## Development Status
- [x] Basic WebSocket server implementation
- [x] Chrome extension with text selection
- [x] WebSocket communication
- [ ] Configuration file support
- [ ] Multiple voice model support
- [ ] Swagger documentation
- [ ] License selection

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
TODO: Choose a license

## Acknowledgments
- [Piper TTS](https://github.com/rhasspy/piper) for the text-to-speech engine
- [Gorilla WebSocket](https://github.com/gorilla/websocket) for WebSocket implementation
