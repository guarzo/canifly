const { app, BrowserWindow, dialog, ipcMain, shell } = require('electron');
const path = require('path');
const http = require('http')

const isDev = !app.isPackaged;

let goProcess;
const maxRetries = 5;
const retryDelay = 1000;
let startTime;

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
        mainWindow.loadURL('http://localhost:5173');
        mainWindow.webContents.openDevTools();
    } else {
        mainWindow.loadFile(path.join(__dirname, 'renderer', 'dist', 'index.html'));
    }

    // Insert CSS to hide scrollbars once the content is loaded
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
        app.quit(); // Quit the app if backend not ready
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
