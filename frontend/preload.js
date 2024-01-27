const axios = require('axios');

window.addEventListener('DOMContentLoaded', () => {
    const replaceText = (selector, text) => {
        const element = document.getElementById(selector)
        if (element) element.innerText = text
    }

    for (const type of ['chrome', 'node', 'electron']) {
        replaceText(`${type}-version`, process.versions[type])
    }
    replaceText('path', process.cwd()); //current directory
    
    axios.get(`http://localhost:6969/${"."}/<empty>/<empty>/window`).then((response) => {
        replaceText('files', response.data);
    })
})