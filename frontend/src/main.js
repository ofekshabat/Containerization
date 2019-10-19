const { app, ipcMain, BrowserWindow } = require('electron');

let listWindow;

global.sharedObject = {
	selectedContainerName: null
}

ipcMain.on('container:create', (e, container) => {
	listWindow.webContents.send('container:create', container);
});

ipcMain.on('container:edit', (e, container) => {
	listWindow.webContents.send('container:edit', container);
});

function createWindow() {
	listWindow = new BrowserWindow({ width: 1000, height: 600, icon: 'assets/container.png' });
	listWindow.loadFile('src/containers/list.html');

	listWindow.on('closed', () => {
		// Dereference the window object, usually you would store windows
		// in an array if your app supports multi windows, this is the time
		// when you should delete the corresponding element.
		listWindow = null;
	});
}

app.on('ready', createWindow);

app.on('activate', () => {
	if (listWindow === null) {
		createWindow();
	}
});
