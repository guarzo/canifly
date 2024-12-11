const { app, BrowserWindow, dialog, ipcMain, shell } = require('electron');
const path = require('path');
const http = require('http');
const { spawn } = require('child_process'); // Import spawn

const isDev = !app.isPackaged;

let goProcess;
const maxRetries = 5;
const retryDelay = 1000;
let startTime;

function startGoBackend() {
    if (isDev) {
        // Start Go backend in dev mode
        // Make sure 'go' is in PATH and you can run `go run .` from project root
        goProcess = spawn('go', ['run', '.'], { cwd: process.cwd() });
    } else {
        // Start the packaged backend binary
        // The backend was placed via extraResources as "canifly-backend.exe" (for Windows)
        // The resources path = process.resourcesPath
        const backendPath = process.platform === 'win32' ?
            path.join(process.resourcesPath, 'canifly-backend.exe') :
            path.join(process.resourcesPath, 'canifly-backend');

        goProcess = spawn(backendPath, [], {
            cwd: process.resourcesPath
        });
    }

    goProcess.stdout.on('data', (data) => {
        console.log(`Go backend: ${data}`);
    });

    goProcess.stderr.on('data', (data) => {
        console.error(`Go backend error: ${data}`);
    });

    goProcess.on('close', (code) => {
        console.log(`Go backend exited with code ${code}`);
    });
}

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

    mainWindow.setMenuBarVisibility(false);

    if (isDev) {
        mainWindow.loadURL('http://localhost:5173');
        mainWindow.webContents.openDevTools();
    } else {
        mainWindow.loadFile(path.join(__dirname, 'renderer', 'dist', 'index.html'));
    }

    mainWindow.webContents.on('did-finish-load', () => {
        const customCSS = `
            html, body {
              overflow: auto !important;
            }
            ::-webkit-scrollbar {
              width: 8px;
            }
            ::-webkit-scrollbar-track {
              background: #1f2937;
            }
            ::-webkit-scrollbar-thumb {
              background-color: #14b8a6;
              border-radius: 4px;
              border: 2px solid #1f2937;
            }
            ::-webkit-scrollbar-thumb:hover {
              background-color: #0d9488;
            }
          `;
        mainWindow.webContents.insertCSS(customCSS);
    });

    ipcMain.on('close-window', () => {
        mainWindow.close();
    });

    ipcMain.handle('choose-directory', async (event, defaultPath) => {
        const options = {
            properties: ['openDirectory']
        };
        if (defaultPath) {
            options.defaultPath = defaultPath;
        }
        const result = await dialog.showOpenDialog(options);
        if (result.canceled || result.filePaths.length === 0) {
            return null;
        }
        return result.filePaths[0];
    });

    ipcMain.handle('open-external', (event, url) => {
        if (url && typeof url === 'string') {
            shell.openExternal(url);
        }
    });

    mainWindow.on('closed', function () {
        if (goProcess) {
            goProcess.kill();
            console.log('Go process terminated');
        }
        app.quit();
        process.exit();
    });
}

app.whenReady().then(async () => {
    startTime = Date.now();
    startGoBackend(); // Start the Go backend first
    checkBackendReady(maxRetries);

    app.on('activate', function () {
        if (BrowserWindow.getAllWindows().length === 0) createWindow();
    });
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') app.quit();
});

function checkBackendReady(retries) {
    if (retries === 0) {
        console.error("Go backend did not start in time. Shutting down.");
        if (goProcess) {
            goProcess.kill();
        }
        app.quit();
        return;
    }

    http.get('http://localhost:8713/static/', (res) => {
        if (res.statusCode === 200) {
            const totalDelay = Date.now() - startTime;
            console.log(`Go backend is ready, launching Electron window after ${totalDelay / 1000} seconds.`);
            createWindow();
        } else {
            console.log(res.statusCode);
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
