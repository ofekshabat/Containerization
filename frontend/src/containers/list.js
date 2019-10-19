const { remote, ipcRenderer } = require('electron');
const { dialog } = remote;
const containerApi = require('./api');
const currentWindow = remote.getCurrentWindow();

ipcRenderer.on('container:create', (e, container) => {
	app.containers.push(container);
});

ipcRenderer.on('container:edit', (e, container) => {
    loadContainers();
});

const app = new Vue({
    el: "#app",
    data: {
        containers: []
    },
    methods: {
        showWindow: function (filePath) {
            const win = new remote.BrowserWindow({ parent: currentWindow, modal:true, width: 700, height: 500 });
            win.loadFile(filePath);
            return win;
        },
        showCreateWindow: function () {
            this.showWindow('src/containers/create.html');
        },
        showCreateImageWindow: function () {
            this.showWindow('src/images/create.html');
        },
        showCreatePackageWindow: function () {
            this.showWindow('src/packages/create.html');
        },
        showImportWindow: function() {
            this.showWindow('src/containers/import.html');
        },
        showEditWindow: function (container) {
            remote.getGlobal('sharedObject').selectedContainerName = container.containerName;
            this.showWindow('src/containers/edit.html');
        },
        showExportWindow: async function (container) {
            const path = dialog.showSaveDialog(currentWindow, {
                title: 'Export',
                message: 'Export a container',
                defaultPath: `${container.containerName}.tar.gz`
            });
            if (path) {
                await containerApi.export(container.containerName, path)
            }
        },
        showDeleteDialog: async function (container) {
            const result = dialog.showMessageBox(currentWindow, {
                type: 'warning',
                buttons: [ "Delete", "Cancel" ],
                title: 'Delete',
                message: `Delete container ${container.containerName}?`
            });
            if (result == 0) {
                await containerApi.delete(container.containerName);
                const index = this.containers.indexOf(container);
                this.containers.splice(index, 1);
            }
        },
        startContainer: async function (container) {
            await containerApi.start(container.containerName);
            container.state = "running";
            // this.showWindow('src/terminal.html');
            
        },
        restartContainer: async function (container) {
            await containerApi.restart(container.containerName);
            
        },
        stopContainer: async function (container) {
            await containerApi.stop(container.containerName);
            container.state = "stopped";
        }
    }
});

async function loadContainers() {
    app.containers = await containerApi.getAll();
}
loadContainers();