export function downloadFile(userId, path, itemName) {
    window.location.href=`/api/v1/files/?userId=${userId}&path=${path}&filename=${encodeURIComponent(itemName)}`;
}

export function deleteFile(userId, path, itemName, itemType) {
    return new Promise((resolve, reject) => {
        if (itemType === "file") {
            fetch(`/api/v1/files/?userId=${userId}&path=${path}&filename=${encodeURIComponent(itemName)}`, {
                method: 'DELETE',
                credentials: 'include'
            })
                .then(() => resolve())
                .catch(error => reject(error));
        } else {
            fetch(`/api/v1/files/folder?userId=${userId}&path=${path}&folderName=${encodeURIComponent(itemName)}`, {
                method: 'DELETE',
                credentials: 'include'
            })
                .then(() => resolve())
                .catch(error => reject(error));
        }
    })
}

export function renameFile(userId, path, actualName, newName) {
    return new Promise((resolve, reject) => {
        if (newName) {
            fetch(`/api/v1/files/?userId=${userId}&path=${path}&filename=${encodeURIComponent(actualName)}&newName=${encodeURIComponent(newName)}`, {
                method: 'PATCH',
                credentials: 'include'
            })
                .then(() => resolve())
                .catch(error => console.error('Error renaming file:', error));
        } else {
            reject()
        }
    })
}

export function createFolder(userId, path, folderName) {
    return new Promise((resolve, reject) => {
        if (folderName) {
            fetch(`/api/v1/files/folder?userId=${userId}&path=${path}&folderName=${encodeURIComponent(folderName)}`, {
                method: 'POST',
                credentials: 'include'
            })
                .then(() => resolve())
                .catch(error => reject(error));
        } else {
            reject()
        }
    })
}

export function createShare(path, date, isPublic) {
    return new Promise((resolve, reject) => {
        fetch(`/api/v1/share/?path=${path}&date=${date}&public=${isPublic}`, {
            method: 'POST',
            credentials: 'include'
        })
            .then(response => response.text().then(shareId => resolve(shareId)))
            .catch(error => reject(error));
    })
}

export function getFolderContent(userId, path) {
    return new Promise((resolve, reject) => {
        fetch(`/api/v1/files/content?userId=${userId}&path=${path}`, {
            method: 'GET',
            credentials: 'include'
        })
            .then(response => response.json())
            .then(data => {
                resolve(data)
            })
            .catch(error => reject(error));
    })
}
