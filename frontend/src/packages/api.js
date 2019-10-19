const axios = require('../../lib/axios.min');
const apiUrl = 'http://localhost:8000/packages';

exports.get = async function (name) {
    const response = await axios.get(`${apiUrl}/${name}`);
    return response.data;
}

exports.getAll = async function () {
    const response = await axios.get(apiUrl);
    return response.data;
}

exports.create = async function (packageInfo) {
    const response = await axios.post(`${apiUrl}/create`, packageInfo);
    return response.data;
}

exports.edit = async function (name, packageInfo) {
    const response = await axios.put(`${apiUrl}/${name}`, packageInfo);
    return response.data;
}

exports.delete = async function (name) {
    const response = await axios.delete(`${apiUrl}/${name}`);
    return response.data;
}
