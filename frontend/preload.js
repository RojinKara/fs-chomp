const axios = require('axios');

window.addEventListener('DOMContentLoaded', () => {
    const replaceText = (selector, text) => {
        const element = document.getElementById(selector)
        if (element) element.innerText = text
    }

    for (const type of ['chrome', 'node', 'electron']) {
        replaceText(`${type}-version`, process.versions[type])
    }
    let directory = encodeURIComponent("C:\\Users\\camde\\Documents\\fs-chomp");
    axios.get(`http://localhost:6969/tree/${directory}/node_modules,.git/%3Cempty%3E/`).then((response) => {
        let table = document.createElement('table')
        for (let i = 0; i < response.data.length; i++) {
            let row = document.createElement('tr')
            let path = document.createElement('td');
            path.innerText = response.data[i].name;
            let file = document.createElement('td');
            file.innerText = response.data[i].isFile;
            row.appendChild(path);
            row.appendChild(file);
            table.appendChild(row);
        }
        document.querySelector('body').appendChild(table);
    })

    // axios.get(`http://localhost:6969/search/C%3A%5CUsers%5Ccamde%5CDocuments%5Cfs-chomp/node_modules,.git/%3Cempty%3E/window`).then((response) => {
    //     replaceText('files', JSON.stringify(response.data));
    // })
})