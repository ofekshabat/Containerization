const { remote, ipcRenderer } = require('electron');
const containerApi = require('./api');
const imageApi = require('../images/api');

const containerName = remote.getGlobal('sharedObject').selectedContainerName;

const app = new Vue({
    el: "#app",
    data: {
        container: {},
        images: []
    },
    methods: {
        editContainer: async function() {
            await containerApi.edit(containerName, this.container);
            ipcRenderer.send('container:edit', this.container);
            remote.getCurrentWindow().close();
        }
    }
});

async function loadData() {
    app.container = await containerApi.get(containerName);
    app.images = await imageApi.getAll();
}
loadData();