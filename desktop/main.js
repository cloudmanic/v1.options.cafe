const electron = require('electron');
const path = require('path');
const url = require('url');
const ipc = require('electron').ipcMain;

const app = electron.app;
const BrowserWindow = electron.BrowserWindow;

let mainWindow = null;

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', function () {
  createMainWindow();
});

// Quit when all windows are closed.
app.on('window-all-closed', function ()
{
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if(process.platform !== 'darwin')
  {
    app.quit();
  }
});

// On activate
app.on('activate', function ()
{
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if(mainWindow === null)
  {
    createMainWindow();
  }
  
});

// ----------------- Main Window ------------------ //

//
// Main app window.
//
function createMainWindow()
{
  // Create the browser window.
  mainWindow = new BrowserWindow({ width: 1200, height: 800, minWidth: 1200, minHeight: 700 });

  // and load the index.html of the app.
  mainWindow.loadURL(url.format({
    pathname: path.join(__dirname, 'index.html'),
    protocol: 'file:',
    slashes: true
  }))

  // Open the DevTools.
  mainWindow.webContents.openDevTools();

  // Emitted when the window is closed.
  mainWindow.on('closed', function ()
  {
    mainWindow = null;
  });
}

/* End File */