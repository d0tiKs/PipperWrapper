document.addEventListener('DOMContentLoaded', () => {
  chrome.storage.local.get(['DEBUG_MODE', 'SERVER_ADDRESS', 'SERVER_PORT'], (result) => {
    document.getElementById('debugMode').checked = result.DEBUG_MODE || false;
    document.getElementById('serverAddress').value = result.SERVER_ADDRESS || '';
    document.getElementById('serverPort').value = result.SERVER_PORT || '';
  });
});

document.getElementById('debugMode').addEventListener('change', (e) => {
  const isDebugMode = e.target.checked;
  chrome.storage.local.set({ DEBUG_MODE: isDebugMode });
  chrome.runtime.sendMessage({
    type: 'DEBUG_MODE_CHANGED',
    value: isDebugMode
  });
});

document.getElementById('serverAddress').addEventListener('change', (e) => {
  const newAddress = e.target.value;
  chrome.storage.local.set({ SERVER_ADDRESS: newAddress });
  chrome.runtime.sendMessage({
    type: 'WEBSOCKET_CONFIG_CHANGED',
    serverAddress: newAddress,
    serverPort: document.getElementById('serverPort').value // Send current port
  });
});

document.getElementById('serverPort').addEventListener('change', (e) => {
  const newPort = e.target.value;
  chrome.storage.local.set({ SERVER_PORT: newPort });
  chrome.runtime.sendMessage({
    type: 'WEBSOCKET_CONFIG_CHANGED',
    serverAddress: document.getElementById('serverAddress').value, // Send current address
    serverPort: newPort
  });
}); 