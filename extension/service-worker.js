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

function connectWebSocket(address, port) {
  if (websocket) {
    websocket.close();
  }
  
  try {
    websocket = new WebSocket(`ws://${address}:${port}/tts`);
    
    websocket.onopen = () => {
      if (isDebugMode) console.log('[Debug] WebSocket connected to /tts');
    };
    
    websocket.onclose = () => {
      if (isDebugMode) console.log('[Debug] WebSocket disconnected');
      websocket = null;
    };
    
    websocket.onerror = (error) => {
      if (isDebugMode) console.error('[Debug] WebSocket error:', error);
    };
  } catch (error) {
    if (isDebugMode) console.error('[Debug] WebSocket connection error:', error);
  }
}

// Init debug mode and context menu
chrome.runtime.onInstalled.addListener(() => {
  chrome.storage.local.get('DEBUG_MODE', (result) => {
    isDebugMode = result.DEBUG_MODE || false;
  });

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
    connectWebSocket(message.address, message.port);
  }
});

// Handle right-click menu selection
chrome.contextMenus.onClicked.addListener((info, tab) => {
  if (info.menuItemId === 'sendToServer') {
    if (isDebugMode) console.log('[Debug] Selected text:', info.selectionText);
    
    // Send text through WebSocket if connected
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      const message = {
        //type: 'selected_text',
        //timestamp: new Date().toISOString(),
        text: info.selectionText
      };
      if (isDebugMode) console.log('[Debug] Sending message:', message);
      websocket.send(JSON.stringify(message));
    } else {
      if (isDebugMode) console.log('[Debug] WebSocket not ready. State:', websocket ? websocket.readyState : 'null');
    }
  }
});


// Initialize WebSocket connection on startup
chrome.storage.local.get(['SERVER_ADDRESS', 'SERVER_PORT'], (result) => {
  if (result.SERVER_ADDRESS && result.SERVER_PORT) {
    connectWebSocket(result.SERVER_ADDRESS, result.SERVER_PORT);
  }
});
