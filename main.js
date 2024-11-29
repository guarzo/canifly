// main.js

const { app, BrowserWindow } = require('electron');
const path = require('path');
const { spawn } = require('child_process');
const http = require('http');

let goProcess;
const maxRetries = 5;
const retryDelay = 1000;
let startTime;

// Determine if the app is in development mode
const isDev = !app.isPackaged;

function createWindow() {
    console.log("Creating Electron window...");
    const mainWindow = new BrowserWindow({
        width: 1000,
        height: 800,
        webPreferences: {
            nodeIntegration: true,
            contextIsolation: false, // Important for communication
        }
    });

    // Disable the menu bar
    mainWindow.setMenuBarVisibility(false);

    if (isDev) {
        mainWindow.loadURL('http://localhost:5173'); // Vite dev server
        mainWindow.webContents.openDevTools();
    } else {
        mainWindow.loadURL(`file://${path.join(__dirname, 'renderer', 'dist', 'index.html')}`);
    }

    mainWindow.on('closed', function () {
        if (goProcess) goProcess.kill();
    });
}

function checkBackendReady(retries) {
    if (retries === 0) {
        console.error("Go backend did not start in time.");
        return;
    }

    // Ensure the port matches your Go backend's listening port (e.g., 8713)
    http.get('http://localhost:8713', (res) => {
        if (res.statusCode === 200) {
            const totalDelay = Date.now() - startTime;
            console.log(`Go backend is ready, launching Electron window after ${totalDelay / 1000} seconds.`);
            createWindow();
        } else {
            retryCheck(retries);
        }
    }).on('error', () => {
        retryCheck(retries);
    });
}

function retryCheck(retries) {
    setTimeout(() => {
        checkBackendReady(retries - 1);
    }, retryDelay);
}

app.on('ready', () => {
    startTime = Date.now();

    if (!isDev) {
        // Define executable location based on environment
        const goExecutable = process.platform === 'win32'
            ? path.join(app.isPackaged ? path.join(process.resourcesPath) : path.join(__dirname, 'dist'), 'canifly-backend.exe')
            : path.join(app.isPackaged ? path.join(process.resourcesPath) : path.join(__dirname, 'dist'), 'canifly-backend');

        console.log(`Starting Go backend from: ${goExecutable}`);
        goProcess = spawn(goExecutable, {
            stdio: 'inherit',
            shell: process.platform === 'win32' // Necessary for Windows
        });

        goProcess.on('error', (err) => {
            console.error('Failed to start Go backend:', err);
        });

        goProcess.on('exit', (code, signal) => {
            if (!isDev) { // Only in production
                console.log(`Go backend exited with code ${code} and signal ${signal}`);
                app.quit();
            }
        });

        // Wait for Go backend to be ready before creating the window
        checkBackendReady(maxRetries);
    } else {
        // In development, assume Go backend is already running via "npm run dev:go"
        // Wait briefly to ensure backend is up before creating the window
        setTimeout(() => {
            createWindow();
        }, 2000); // Adjust delay as necessary
    }
});

app.on('window-all-closed', function () {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});

app.on('activate', function () {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
});
