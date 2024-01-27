const fs = require("node:fs/promises");

const generateFileTreeObject = directoryString => {
    return fs.readdir(directoryString)
        .then(arrayOfFileNameStrings => {
            const fileDataPromises = arrayOfFileNameStrings.map(fileNameString => {
                const fullPath = `${directoryString}/${fileNameString}`;
                return fs.stat(fullPath)
                    .then(fileData => {
                        const file = {};
                        file.filePath = fullPath;
                        file.isFile = fileData.isFile();
                        return file;
                    });
            });

            return Promise.all(fileDataPromises);
        });
};

module.exports = {generateFileTreeObject}