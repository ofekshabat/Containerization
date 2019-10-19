const axios = require('../../lib/axios.min');
const apiUrl = 'http://localhost:8000/containers';

exports.get = async function (name) {
    const response = await axios.get(`${apiUrl}/${name}`);
    return response.data;
}

exports.getAll = async function () {
    const response = await axios.get(apiUrl);
    return response.data;
}

exports.create = async function (containerInfo) {
    containerInfo.state = "stopped";
    const response = await axios.post(`${apiUrl}/create`, containerInfo);
    return response.data;
}

exports.edit = async function (name, containerInfo) {
    const response = await axios.put(`${apiUrl}/${name}`, containerInfo);
    return response.data;
}

exports.delete = async function (name) {
    const response = await axios.delete(`${apiUrl}/${name}`);
    return response.data;
}

exports.start = async function (name) {
    const response = await axios.post(`${apiUrl}/${name}/start`);
    return response.data;
}

exports.restart = async function (name) {
    const response = await axios.post(`${apiUrl}/${name}/restart`);
    return response.data;
}

exports.stop = async function (name) {
    const response = await axios.post(`${apiUrl}/${name}/stop`);
    return response.data;
}

exports.import = async function (path, name) {
    const response = await axios.post(`${apiUrl}/import`, { path: path, containerName: name });
    return response.data;
}

exports.export = async function (name, path) {
    const response = await axios.post(`${apiUrl}/${name}/export`, { path: path });
    return response.data;
}