// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

let isDebugMode = false;
let websocket = null;

let serverAddress;
let serverPort;

function HandleWebSocketConnection(serverAddress, serverPort, message) {

  if (websocket) {
    closeWebSocket();
  }
  
  try {
    websocket = new WebSocket(`ws://${serverAddress}:${serverPort}/tts`);
    
    websocket.onopen = () => {
      if (isDebugMode) console.log('[Debug] WebSocket connected to /tts');
      websocket.send(message)
      if (isDebugMode) console.log('[Debug] Sent to TTS:', message);
      closeWebSocket();
    };
    
    websocket.onclose = () => {
      if (isDebugMode) console.log('[Debug] WebSocket disconnected');
    };
    
    websocket.onerror = (error) => {
      if (isDebugMode) console.error('[Debug] WebSocket error:', error);
    };
  } catch (error) {
    if (isDebugMode) console.error('[Debug] WebSocket connection error:', error);
  }
}
// Function to close the WebSocket connection
function closeWebSocket(statusCode = 1000, reason = '') {
  if (websocket) {
    websocket.close(statusCode, reason); // Close with status code and reason
    websocket = null; // Reset the websocket variable
    if (isDebugMode) console.log('[Debug] WebSocket connection closed with status:', statusCode, 'Reason:', reason);
  }
}
function sendTextToTTS(text) {
  const message = JSON.stringify({ text: text });
  
  // Open the WebSocket connection only if it's not already open
  if (!websocket || websocket.readyState === WebSocket.CLOSED) {
    HandleWebSocketConnection(serverAddress, serverPort, message);
  }

}

// Init debug mode and context menu
chrome.runtime.onInstalled.addListener(() => {
  if (isDebugMode) console.log('[Debug] Adding listener for context menu.')

  chrome.contextMenus.create({
    id: 'sendToServer',
    title: 'Send to server',
    contexts: ['selection']
  });
});

// Update debug mode when changed
chrome.runtime.onMessage.addListener((message) => {
  if (message.type === 'DEBUG_MODE_CHANGED') {
    isDebugMode = message.value;
    if (isDebugMode) console.log('[Debug] Debug mode enabled');
  } else if (message.type === 'WEBSOCKET_CONFIG_CHANGED') {
    // Save the new configuration
    serverAddress = message.serverAddress;
    serverPort = message.serverPort;
    if (isDebugMode) console.log('[Debug] WebSocket configuration updated:', serverAddress, serverPort);
  }
});

// Handle right-click menu selection
chrome.contextMenus.onClicked.addListener((info, tab) => {
  if (isDebugMode) console.log('Context menu item clicked:', info.menuItemId);
  if (info.menuItemId === 'sendToServer') {
    if (isDebugMode) console.log('[Debug] Selected text:', info.selectionText);
    // Use the sendTextToTTS function to send the selected text
    sendTextToTTS(info.selectionText);
  }
});
// Initialize WebSocket connection on startup
chrome.storage.local.get(['DEBUG_MODE', 'SERVER_ADDRESS', 'SERVER_PORT'], (result) => {
  isDebugMode = result.DEBUG_MODE || false;
  serverAddress = result.SERVER_ADDRESS || "127.0.0.1"
  serverPort = result.SERVER_PORT || "8080"
});
