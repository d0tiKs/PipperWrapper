# Text-to-Speech Sender Chrome Extension

A Chrome extension that sends selected text to a local TTS (Text-to-Speech) server via WebSocket for audio playback.

## Features

- Right-click any selected text to send it to your TTS server
- Configurable server address and port
- Debug mode for troubleshooting
- WebSocket communication over `/tts` endpoint

## Installation

To install the extension, follow these steps:
1. Clone this repository or download the source code
2. Open Chrome and navigate to `chrome://extensions/`
3. Enable "Developer mode" in the top right
4. Click "Load unpacked" and select the extension directory

## Configuration

1. Click the extension's settings icon in Chrome
2. Set your TTS server address (e.g., `192.168.1.100`)
3. Set your TTS server port (e.g., `18080`)
4. Optionally enable Debug mode for detailed logging

## Usage

1. Select any text on a webpage
2. Right-click and choose "Send to server"
3. The selected text will be sent to your TTS server for processing

## Server Requirements

Your TTS server should:
- Accept WebSocket connections on the specified address and port
- Listen on the `/tts` endpoint
- Accept JSON messages in the format:
  ```json
  {
    "text": "The selected text"
  }
  ```

## Debugging

To view logs:
1. Go to `chrome://extensions/`
2. Find "Text-to-Speech Sender"
3. Click on "service worker" under "Inspect views"
4. Check the Console tab for debug information

## TODO

- [ ] choose language
- [ ] play/pause/resume
- [ ] small display

## License

TODO Define license
