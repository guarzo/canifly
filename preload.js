// preload.js
const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electronAPI', {
    closeWindow: () => ipcRenderer.send('close-window'),
    chooseDirectory: (defaultPath) => ipcRenderer.invoke('choose-directory', defaultPath),
    openExternal: (url) => ipcRenderer.invoke('open-external', url),
});