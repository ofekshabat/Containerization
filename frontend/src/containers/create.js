const { remote, ipcRenderer } = require('electron');
const currentWindow = remote.getCurrentWindow();
const containerApi = require('./api');
const imageApi = require('../images/api');

const app = new Vue({
    el: "#app",
    data: {
        container: {
            baseImageName: "ubuntu",
            cmdLine: "bash",
            maxCpu: 100,
            maxMemory: 1024,
            maxPids: 100,
            address: "10.10.10.2/24"
        },
        images: []
    },
    methods: {
        createContainer: async function () {
            await containerApi.create(this.container);
            ipcRenderer.send('container:create', this.container);
            currentWindow.close();
        }
    }
});

async function loadData() {
    app.images = await imageApi.getAll();
}
loadData();