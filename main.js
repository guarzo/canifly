const { app, BrowserWindow, ipcMain, dialog } = require('electron');
const fs = require('fs').promises;
const path = require('path');
const axios = require('axios');
const tar = require('tar');
const os = require('os');
const { exec } = require('child_process');

// Determine if the app is in development mode
const isDev = !app.isPackaged;

let settingsDir;
let associations = [];
let failedESICache = {};
let characterCache = {};
let failedHTTPRequests = 0;

// Paths
const settingsFilePath = path.join(app.getPath('userData'), 'settings.json');
const associationsPath = path.join(app.getPath('userData'), 'associations.json');
const failedESICachePath = path.join(app.getPath('userData'), 'failed_esi_requests.json');
const characterCachePath = path.join(app.getPath('userData'), 'character_cache.json');
const defaultPaths = {
    win32: path.join(os.homedir(), 'AppData', 'Local', 'CCP', 'EVE', 'c_ccp_eve_online_tq_tranquility'),
    darwin: path.join(os.homedir(), 'Library', 'Application Support', 'CCP', 'EVE', 'c_ccp_eve_online_tq_tranquility'),
    linux: path.join(os.homedir(), '.local', 'share', 'CCP', 'EVE', 'c_ccp_eve_online_tq_tranquility'),
};

// Load and Save Helpers
async function loadJSON(filePath, defaultValue) {
    try {
        const data = await fs.readFile(filePath, 'utf8');
        return JSON.parse(data);
    } catch {
        return defaultValue;
    }
}

async function saveJSON(filePath, data) {
    try {
        await fs.writeFile(filePath, JSON.stringify(data, null, 2), 'utf8');
        console.log(`Saved data to ${filePath}`);
    } catch (error) {
        console.error(`Failed to save data to ${filePath}:`, error);
    }
}

// Caches and Associations
async function loadCachesAndAssociations() {
    failedESICache = await loadJSON(failedESICachePath, {});
    characterCache = await loadJSON(characterCachePath, {});
    associations = await loadJSON(associationsPath, []);
}

// Get Default Directory Based on Platform
async function getDefaultDirectory() {
    const platform = process.platform;
    if (await isWSL()) {
        const windowsHome = await getWindowsHomeInWSL();
        return path.join(windowsHome, 'AppData', 'Local', 'CCP', 'EVE', 'c_ccp_eve_online_tq_tranquility');
    }
    return defaultPaths[platform];
}

// Function to Detect WSL
async function isWSL() {
    if (os.platform() === 'linux') {
        try {
            const version = await fs.readFile('/proc/version', 'utf8');
            return version.includes('microsoft');
        } catch {
            return false;
        }
    }
    return false;
}

// Get Windows Home Directory in WSL
async function getWindowsHomeInWSL() {
    return new Promise((resolve, reject) => {
        exec('cmd.exe /C "echo %USERPROFILE%"', (error, stdout) => {
            if (error) return reject('Failed to retrieve Windows home directory in WSL.');
            resolve(stdout.trim().replace(/\\/g, '/'));
        });
    });
}

// Electron Window Creation
function createWindow() {
    console.log('Creating Electron window...');
    const mainWindow = new BrowserWindow({
        width: 1000,
        height: 800,
        webPreferences: {
            preload: path.join(__dirname, 'preload.js'),
            contextIsolation: true,
            nodeIntegration: false,
            enableRemoteModule: false,
        },
        frame: false,
    });

    mainWindow.setMenuBarVisibility(false);

    if (isDev) {
        mainWindow.loadURL('http://localhost:5173');
        mainWindow.webContents.openDevTools();
    } else {
        mainWindow.loadFile(path.join(__dirname, 'renderer', 'dist', 'index.html'));
    }

    mainWindow.webContents.on('did-finish-load', () => {
        console.log('Renderer process loaded successfully.');
    });

    mainWindow.webContents.on('did-fail-load', (event, errorCode, errorDescription) => {
        console.error(`Renderer process failed to load: ${errorDescription}`);
    });

    mainWindow.on('closed', () => {
        app.quit();
    });

    return mainWindow;
}

// IPC Handlers
function registerIpcHandlers() {
    ipcMain.handle('load-settings', async () => {
        console.log('Loading settings...');
        return await loadCachesAndAssociations();
    });

    ipcMain.handle('choose-settings-dir', async () => {
        const result = await dialog.showOpenDialog({ properties: ['openDirectory'] });
        if (!result.canceled && result.filePaths.length > 0) {
            settingsDir = result.filePaths[0];
            await saveJSON(settingsFilePath, { settingsDir });
            return { success: true, settingsDir };
        }
        return { success: false };
    });

    ipcMain.handle('reset-to-default-directory', async () => {
        const defaultDir = await getDefaultDirectory();
        const dirExists = await fs.access(defaultDir).then(() => true).catch(() => false);
        if (dirExists) {
            settingsDir = defaultDir;
            await saveJSON(settingsFilePath, { settingsDir });
            return { success: true, settingsDir };
        }
        return { success: false, message: 'Default directory does not exist.' };
    });
}

// App Initialization
app.whenReady().then(async () => {
    console.log('App is ready. Initializing...');
    settingsDir = await loadJSON(settingsFilePath, {}).settingsDir || (await getDefaultDirectory());
    await loadCachesAndAssociations();
    createWindow();

    app.on('activate', () => {
        if (BrowserWindow.getAllWindows().length === 0) createWindow();
    });
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});

// Register IPC Handlers
registerIpcHandlers();
