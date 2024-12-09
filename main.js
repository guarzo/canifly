// main.js

const { app, BrowserWindow, ipcMain } = require('electron');
const path = require('path');
const http = require('http')

// Determine if the app is in development mode
const isDev = !app.isPackaged;

let goProcess;
const maxRetries = 5;
const retryDelay = 1000;
let startTime;

// Create the main Electron window
function createWindow() {
    console.log("Creating Electron window...");
    const mainWindow = new BrowserWindow({
        width: 1000,
        height: 800,
        webPreferences: {
            nodeIntegration: false,
            contextIsolation: true,
            preload: path.join(__dirname, 'preload.js'),
            enableRemoteModule: false,
        },
        frame: false,
    });

    // Disable the menu bar
    mainWindow.setMenuBarVisibility(false);

    if (isDev) {
        mainWindow.loadURL('http://localhost:5173'); // Adjust to your dev server
        mainWindow.webContents.openDevTools();
    } else {
        mainWindow.loadFile(path.join(__dirname, 'renderer', 'dist', 'index.html'));
    }

    // Listen for close window request
    ipcMain.on('close-window', () => {
        mainWindow.close(); // Close the window
    });

    mainWindow.on('closed', function () {
        if (goProcess) {
            goProcess.kill();
            console.log('Go process terminated');
        }
        app.quit(); // Quit the Electron app completely
        process.exit(); // Force exit if necessary
    });
}

// Initialize and create the window
app.whenReady().then(async () => {

    startTime = Date.now();
    checkBackendReady(maxRetries);

    app.on('activate', function () {
        if (BrowserWindow.getAllWindows().length === 0) createWindow();
    });
});

// Quit when all windows are closed, except on macOS
app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') app.quit();
});


function checkBackendReady(retries) {
    if (retries === 0) {
        console.error("Go backend did not start in time.");
        return;
    }

    http.get('http://localhost:8713/static/', (res) => {
        if (res.statusCode === 200) {
            const totalDelay = Date.now() - startTime;
            console.log(`Go backend is ready, launching Electron window after ${totalDelay / 1000} seconds.`);
            createWindow();
        } else {
            console.log(res.statusCode)
            retryCheck(retries);
        }
    }).on('error', () => {
        console.log("retrying...............................")
        retryCheck(retries);
    });
}

function retryCheck(retries) {
    setTimeout(() => {
        checkBackendReady(retries - 1);
    }, retryDelay);
}
