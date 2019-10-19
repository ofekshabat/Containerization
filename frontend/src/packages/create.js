const { remote, ipcRenderer } = require('electron');
const currentWindow = remote.getCurrentWindow();
const packageApi = require('./api');

const app = new Vue({
    el: "#app",
    data: {
        package: {}
    },
    methods: {
        createPackage: async function () {
            await packageApi.create(this.package);
            currentWindow.close();
        }
    }
});
