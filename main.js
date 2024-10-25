const { app, BrowserWindow } = require('electron');
const path = require('path');
const { spawn } = require('child_process');
const http = require('http');

let goProcess;
const maxRetries = 5;
const retryDelay = 1000;
let startTime;

function createWindow() {
    console.log("Creating Electron window...");
    const mainWindow = new BrowserWindow({
        width: 1000,
        height: 800
    });

    // Disable the menu bar
    mainWindow.setMenuBarVisibility(false);

    mainWindow.loadURL(`http://localhost:8080`);


    mainWindow.on('closed', function () {
        if (goProcess) goProcess.kill();
    });
}



function checkBackendReady(retries) {
    if (retries === 0) {
        console.error("Go backend did not start in time.");
        return;
    }

    http.get('http://localhost:8080', (res) => {
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

    // Define executable location based on environment
    const goExecutable = process.platform === 'win32'
        ? path.join(app.isPackaged ? path.join(process.resourcesPath) : path.join(__dirname, 'dist'), 'canifly-backend.exe')
        : path.join(app.isPackaged ? path.join(process.resourcesPath) : path.join(__dirname, 'dist'), 'canifly-backend');
    
    console.log(`Starting Go backend from: ${goExecutable}`);
    goProcess = spawn(goExecutable);

    goProcess.stdout.on('data', (data) => {
        console.log(`Go output: ${data}`);
    });

    goProcess.stderr.on('data', (data) => {
        console.error(`Go error: ${data}`);
    });

    checkBackendReady(maxRetries);
});

app.on('window-all-closed', function () {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});

app.on('activate', function () {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
});
