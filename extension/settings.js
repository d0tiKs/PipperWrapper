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
  chrome.storage.local.set({ SERVER_ADDRESS: e.target.value });
});

document.getElementById('serverPort').addEventListener('change', (e) => {
  chrome.storage.local.set({ SERVER_PORT: e.target.value });
}); 