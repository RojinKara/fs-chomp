const axios = require('axios');

/*sets a default directory if first run*/
if (localStorage.getItem("directory") === null || localStorage.getItem("directory") === undefined) {
    let directory = encodeURIComponent("C:\\Users\\camde\\Documents\\fs-chomp");
    localStorage.setItem("directory", directory);
}

window.addEventListener('DOMContentLoaded', () => {
    const replaceText = (selector, text) => {
        const element = document.getElementById(selector)
        if (element) element.innerText = text
    }
    if (!window.location.href.includes("index")) { //must be search.html
        const search = localStorage.getItem("search");

        axios.get(`http://localhost:6969/search/${localStorage.getItem("directory")}/node_modules,.git,.idea,.vscode,__pycache__/.gitignore,.xml/${search}`).then((response) => {
            if (response == null) {
                document.querySelector('body').appendChild(document.createElement('h1').innerText = "No results found");
                return;
            }
            const table = document.createElement('table');
            table.classList.add("w-3/5")
            /*adds a row for each row of data*/
            for (let i = 0; i < response.data.length; i++) {
                let row = document.createElement('tr');
                // /*FULL PATHS OF RESULTS*/
                // let path = document.createElement('td');
                // path.innerText = response.data[i].file;
                // let lineNumber = document.createElement('td');
                // lineNumber.innerText = response.data[i].lineNumber;
                /*relative path to search*/
                let relativePath = document.createElement('td');
                let div = document.createElement("div")
                div.innerHTML = `<div class="block w-full p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">

<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">${response.data[i].file}</h5>
<h6 class="font-normal text-gray-700 dark:text-gray-400">${response.data[i].line}</h6>
</div>`;
                relativePath.appendChild(div)
                /*line containing match*/
                //let line = document.createElement('td');
                //line.innerText = response.data[i].line;
                row.appendChild(relativePath);
                //row.appendChild(line);
                table.appendChild(row);
            }
            document.querySelector('body').appendChild(table);
        })
    } else { //index.html
        axios.get(`http://localhost:6969/tree/${localStorage.getItem("directory")}/node_modules,.git,.idea,.vscode,__pycache__/.gitignore,.xml,.drawio/`).then((response) => {
            if (response == null) {
                document.createElement('h1').innerText = "No results found";
                return;
            }
            /*create a table*/
            let table = document.createElement('table');
            /*add a go to parent directory*/
            let row = document.createElement('tr');
            row.classList.add("my-4", "cursor-pointer")
            let icon = document.createElement('td');
            icon.innerHTML = '<img class = "icon mx-2" src = "./public/arrow-left.svg" alt="Go back">';
            let dot = document.createElement('td');
            dot.innerText = '..';
            row.addEventListener('click', () => {
                const currentDirectory = localStorage.getItem("directory");
                for (let i = currentDirectory.length - 1; i >= 0; i--) {
                    if (currentDirectory[i] === "C" && currentDirectory[i - 1] === "5" && currentDirectory[i - 2] === "%") { //look for backslash
                        localStorage.setItem("directory", currentDirectory.substring(0, i - 2));
                        window.location.reload();
                        break;
                    }
                }
            })
            row.appendChild(icon);
            row.appendChild(dot);
            table.appendChild(row);

            /*adds a row for each row of data*/
            for (let i = 0; i < response.data.length; i++) {
                let row = document.createElement('tr');
                /*current files available in current directory*/
                let path = document.createElement('td');
                path.innerText = response.data[i].name;
                /*is a file or a directory*/
                let isFile = document.createElement('td');
                /*icons for the current directory*/
                if (!response.data[i].isFile) {
                    /*adds a clickable to directories only*/
                    isFile.innerHTML = '<img class = "icon mx-2" src = "./public/folder.svg" alt="Folder">';
                    row.addEventListener('click', () => {
                        localStorage.setItem("directory", encodeURIComponent(response.data[i].fullPath));
                        window.location.reload();
                    })
                } else {
                    isFile.innerHTML = '<img class = "icon mx-2" src = "./public/file-text.svg" alt="File">';
                }
                row.appendChild(isFile);
                row.appendChild(path);
                table.appendChild(row);
            }
            document.querySelector('#load').appendChild(table);
        })

        // axios.get(`http://localhost:6969/search/C%3A%5CUsers%5Ccamde%5CDocuments%5Cfs-chomp/node_modules,.git/%3Cempty%3E/window`).then((response) => {
        //     replaceText('files', JSON.stringify(response.data));
        // })
    }
})