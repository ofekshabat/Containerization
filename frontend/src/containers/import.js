const { remote, ipcRenderer } = require('electron');
const { dialog } = remote;
const currentWindow = remote.getCurrentWindow();
const containerApi = require('./api');

const app = new Vue({
    el: "#app",
    data: {
        path: "",
        containerName: ""
    },
    methods: {
        browseFile: function () {
            const paths = dialog.showOpenDialog(currentWindow, {
                title: 'Import',
                message: 'Import a container',
                filters: [
                    { name: 'Containers (.tar.gz)', extensions: [ 'gz' ] },
                    { name: 'All Files', extensions: [ '*' ] }
                ]
            });
            if (paths.length > 0) {
                this.path = paths[0];
            }
        },
        importContainer: async function () {
            const containerInfo = await containerApi.import(this.path, this.containerName);
            ipcRenderer.send('container:create', containerInfo);
            remote.getCurrentWindow().close();
        }
    }
});