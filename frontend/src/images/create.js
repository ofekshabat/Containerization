const { remote } = require('electron');
const currentWindow = remote.getCurrentWindow();
const imageApi = require('./api');

const app = new Vue({
    el: "#app",
    data: {
        image: {
            baseImageName: ""
        },
        images: []
    },
    methods: {
        createImage: async function () {
            await imageApi.create(this.image);
            currentWindow.close();
        }
    }
});

async function loadImages() {
    app.images = await imageApi.getAll();
}
loadImages();