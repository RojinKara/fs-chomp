const { app, BrowserWindow } = require('electron/main');
const path = require('node:path');
const { generateFileTreeObject } = require('./file.js');
const electronReload = require('electron-reload');

function createWindow () {
    const win = new BrowserWindow({
        width: 800,
        height: 600,
        webPreferences: {
            nodeIntegration: true, //allows node.js in the browser
            preload: path.join(__dirname, 'preload.js')
        }

    })
    // generateFileTreeObject(process.cwd()).then((files) => {
    //     console.log(files);
    // })
    win.loadFile('index.html');
    // generateFileTreeObject(process.cwd()).then((files) => {
    //     let temp = "";
    //     for (let file in files) {
    //         temp += file + "\n";
    //     }
    //     win.webContents.send('files', temp);
    //
    // });
}

app.whenReady().then(() => {
    createWindow()

    app.on('activate', () => {
        if (BrowserWindow.getAllWindows().length === 0) {
            createWindow()
        }
    })
})

app.on('window-all-closed', () => {
        if (process.platform !== 'darwin') {
        app.quit()
    }
})