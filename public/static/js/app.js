import mime from 'https://cdn.jsdelivr.net/npm/mime@4.0.4/+esm'

export let ActualPath = new URLSearchParams(window.location.search).get('path') ?? '';
const jwt = getCookie('token');
export let UserId = parseJwt(jwt).user_id;
export let SelectedItem = new Map();
SelectedItem.set('Name', "");
SelectedItem.set('Type', "");
const MAX_PARALLEL_UPLOADS = 3;

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

function parseJwt(token) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

function formatBytes(bytes) {
    const units = ['o', 'Ko', 'Mo', 'Go', 'To', 'Po', 'Exa']; // Ajoutez d'autres unités si nécessaire
    let index = 0;

    while (bytes >= 1024 && index < units.length - 1) {
        bytes /= 1024;
        index++;
    }

    return `${bytes.toFixed(2)} ${units[index]}`; // Format avec 2 décimales
}

export function openPopup(title, content, placeholder) {
    return new Promise((resolve, reject) => {
        const popup = document.getElementById('inputPopupOverlay');
        document.getElementById('InputPopupInput').value = content
        popup.querySelector('#InputPopupTitle').innerHTML = title;
        popup.querySelector('#InputPopupInput').placeholder = placeholder;
        popup.classList.remove('hidden');

        document.getElementById('submitPopupButton').onclick = function() {
            let value = document.getElementById('InputPopupInput').value
            if (value) {
                closePopupInput();
                resolve(value);
            } else {
                closePopupInput();
                reject(new Error('Popup submission aborted'));
            }
        };

        document.getElementById('abortPopupButton').onclick = function() {
            closePopupInput();
            reject(new Error('Popup submission aborted'));
        };
    });
}

export function openPopupWarning(title, content) {
    return new Promise((resolve, reject) => {
        const popup = document.getElementById('warningPopupOverlay');
        popup.querySelector('#warningPopupTitle').innerHTML = title;
        popup.querySelector('#warningPopupContent').innerHTML = content;
        popup.classList.remove('hidden');

        document.getElementById('warningYesPopupButton').onclick = function() {
            closePopupWarning()
            resolve(true);
        };

        document.getElementById('warningNoPopupButton').onclick = function() {
            closePopupWarning()
            resolve(false);
        };
    });
}

function closePopupInput() {
    const popup = document.getElementById('inputPopupOverlay');
    popup.classList.add('hidden');
    document.getElementById('InputPopupInput').value = ""
}

function closePopupWarning() {
    const popup = document.getElementById('warningPopupOverlay');
    popup.classList.add('hidden');
}


// Utility function to update the URL without reloading
export function changeFolder(folderId) {
    ActualPath += folderId + "/";
    history.pushState({ path: ActualPath }, null, '?path=' + ActualPath);
    loadFolderContent();
}

// Load folder content dynamically
export function loadFolderContent() {
    fetch(`/api/v1/files/content?userId=${UserId}&path=${ActualPath}`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            document.getElementById('file-list').innerHTML = renderFolderContent(data);
        })
        .catch(error => console.error('Error loading content:', error));
}

// Generate folder and file HTML
function renderFolderContent(data) {
    if (!data || data.length === 0) {
        return '';  // Ne rien faire si data est vide
    }
    return data.map(item => {
        const parsedDate = new Date(item.Time);

        // Extraire les composants de la date
        const day = String(parsedDate.getDate()).padStart(2, '0'); // Jour
        const month = String(parsedDate.getMonth() + 1).padStart(2, '0'); // Mois (0-11, donc +1)
        const year = String(parsedDate.getFullYear()); // Année (deux derniers chiffres)
        const hours = String(parsedDate.getHours()).padStart(2, '0'); // Heures
        const minutes = String(parsedDate.getMinutes()).padStart(2, '0'); // Minutes

        // Construire la chaîne formatée

        if (item.Type === 'file') {
            return `
                <tr class="hover:bg-gray-100 context-menu-target" data-info="${item.Name}" data-type="file">
                    <td class="px-4 py-2 flex items-center">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="24px" height="24px"><path fill="#c9dcf4" d="m71.99,88.11c-14.66,1.19-29.31,1.19-43.97,0-3.67-.3-6.52-3.29-6.71-6.96-1.04-20.43-1.04-40.86,0-61.3.19-3.68,3.05-6.66,6.72-6.96,9.87-.8,19.74-1.06,29.62-.78,2.13.06,4.14.98,5.62,2.52,4.46,4.66,9.01,9.34,13.6,13.94,1.47,1.47,2.33,3.45,2.39,5.53.43,15.69.24,31.37-.55,47.06-.19,3.67-3.04,6.66-6.71,6.96Z"/><path fill="#4a254b" d="m50,73c2.57,0,4.68-1.94,4.97-4.43.03-.3-.19-.57-.49-.57-1.73,0-7.22,0-8.95,0-.3,0-.53.27-.49.57.28,2.49,2.4,4.43,4.97,4.43Z"/><circle cx="39.5" cy="61.5" r="5.5" fill="#fff"/><circle cx="39.5" cy="61.5" r="2.5" fill="#4a254b"/><circle cx="60.5" cy="61.5" r="5.5" fill="#fff"/><circle cx="60.5" cy="61.5" r="2.5" fill="#4a254b"/><path fill="#6b96d6" d="m68.63,31.42h10.08c-.41-1.06-1.03-2.04-1.85-2.87-4.58-4.6-9.13-9.28-13.6-13.94-.9-.94-2-1.64-3.21-2.07l.02,10.32c.01,4.72,3.84,8.54,8.56,8.54Z"/></svg>
                        <span>${item.Name}</span>
                    </td>
                    <td class="px-4 py-2">${formatBytes(item.Size)}</td>
                    <td class="px-4 py-2">${mime.getType(item.Name)}</td>
                    <td class="px-4 py-2">${day}/${month}/${year} ${hours}:${minutes}</td>
                </tr>`
        } else if (item.Type === 'directory') {
            return `
                <tr class="hover:bg-gray-100 context-menu-target" data-info="${item.Name}" data-type="directory" onclick="changeFolder('${item.Name}')">
                    <td class="px-4 py-2 flex items-center">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="24px" height="24px"><path fill="#ff931e" d="m19.2,29.71c-6.62-1.93-6.83,3.92-7.06,9.31-.4,9.57,2.28,18.79,2.35,28.31.02,3.14-.97,7.87.54,10.77,1.23,2.36,3.68,2.04,6,2.04,10.13,0,20.25,0,30.38,0s19.22.83,28.87.86c1.9,0,5.28.08,6.69-1.55s1.14-5.74,1.28-7.76c.64-9.01.19-17.85-1.05-26.79-.84-6.04-.05-12.58-7.29-13.93-9.64-1.79-19.97.28-29.71-.54-9.57-.8-19.38-2.56-28.98-1.24-1.9.26-1.09,3.15.8,2.89,8.76-1.21,18.09.27,26.81,1.18s18.06.26,26.92.41c1.56.03,3.5-.14,4.94.54,3.31,1.56,2.54,3.95,2.91,6.79.68,5.18,1.76,10.37,1.89,15.6.09,3.65.22,7.38,0,11.02s1.06,9.98-3.81,10.38c-8.45.7-17.06-.92-25.52-.86-9.79.07-19.59-.03-29.38-.02-1.77,0-6.15.82-7.68.02-2.13-1.11-1.65-2.99-1.64-4.87.01-4.28.14-8.57-.39-12.83-.63-5.07-1.72-10.03-1.92-15.15-.11-2.68-1.19-12.98,3.24-11.69,1.86.54,2.65-2.35.8-2.89h0Z"/><path fill="#fec933" d="m88.27,31.39c-.42-3.57-3.83-6.61-7.38-6.77-14.07-.57-28.14-.74-42.2-.53l-1.16-2.9c-.72-1.79-2.45-2.96-4.38-2.96h-16.32c-2.47,0-4.52,1.91-4.7,4.38l-.82,11.44h.15c-.76,8.27-.67,16.55.26,24.83.42,3.57,3.83,6.61,7.38,6.77,20.59.83,41.19.83,61.78,0,3.56-.15,6.97-3.2,7.38-6.77,1.04-9.16,1.04-18.32,0-27.49Z"/><path fill="#fec933" d="m82.98,81.05c-21.99.96-43.98.96-65.97,0-2.49-.11-4.56-1.94-5.01-4.39-2.42-13.23-3.46-26.39-2.8-39.47.14-2.73,2.33-4.91,5.06-5.04,23.82-1.13,47.65-1.13,71.47,0,2.73.13,4.92,2.32,5.06,5.04.66,13.08-.38,26.24-2.8,39.47-.45,2.45-2.52,4.29-5.01,4.39Z"/><path fill="#4a254b" d="m50,70c2.57,0,4.68-1.94,4.97-4.43.03-.3-.18-.57-.49-.57-1.72,0-7.23,0-8.96,0-.3,0-.52.27-.49.57.28,2.49,2.4,4.43,4.97,4.43Z"/><circle cx="37" cy="58" r="6" fill="#fff"/><circle cx="37" cy="58" r="2.75" fill="#4a254b"/><path fill="#4a254b" d="m88.82,33.19c-.11,0-.22-.04-.31-.11-.8-.65-1.77-1.02-2.8-1.07-23.7-1.12-47.73-1.12-71.42,0-.83.04-1.63.29-2.33.73-.23.15-.54.08-.69-.16-.15-.23-.08-.54.16-.69.84-.53,1.82-.84,2.81-.88,23.73-1.12,47.79-1.12,71.52,0,1.24.06,2.41.5,3.38,1.29.21.17.25.49.07.7-.1.12-.24.19-.39.19Z"/><circle cx="63" cy="58" r="6" fill="#fff"/><circle cx="63" cy="58" r="2.75" fill="#4a254b"/></svg>
                        <span>${item.Name}</span>
                    </td>
                    <td class="px-4 py-2">_</td>
                    <td class="px-4 py-2">_</td>
                    <td class="px-4 py-2">${day}/${month}/${year} ${hours}:${minutes}</td>
                </tr>`
        }
    }).join('');
}

// Handle browser navigation using the back button
window.addEventListener('popstate', event => {
    ActualPath = (event.state && event.state.path !== null) ? event.state.path : "";
    loadFolderContent();
});

// ---------------------- UPLOAD FILES ----------------------

const dropZone = document.getElementById('file-container');

// Prevent default behavior for drag-and-drop events
['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
    dropZone.addEventListener(eventName, e => e.preventDefault());
});

dropZone.addEventListener('drop', async (event) => {
    event.preventDefault();
    const files = Array.from(event.dataTransfer.files);
    const url = `/api/v1/files/?UserId=${UserId}&path=${ActualPath}`; // URL de votre API pour l'upload

    // Fonction pour uploader un fichier avec suivi de la progression
    const uploadFile = (file) => {
        return new Promise((resolve, reject) => {
            const xhr = new XMLHttpRequest();
            const formData = new FormData();
            formData.append('file', file);

            xhr.open('POST', url, true);

            // Affichage de la progression
            xhr.upload.onprogress = (event) => {
                if (event.lengthComputable) {
                    const percentComplete = (event.loaded / event.total) * 100;
                    console.log(`Upload progress for ${file.name}: ${percentComplete.toFixed(2)}%`);
                    // Ici, vous pouvez mettre à jour une barre de progression ou un autre indicateur dans votre UI
                }
            };

            xhr.onload = () => {
                if (xhr.status === 200) {
                    console.log(`File uploaded successfully: ${file.name}`);
                    resolve(xhr.response);
                } else {
                    console.error(`Upload failed for ${file.name}`);
                    reject(new Error(`Upload failed with status ${xhr.status}`));
                }
            };

            xhr.onerror = () => {
                console.error('Error during upload');
                reject(new Error('Upload error'));
            };

            xhr.send(formData);
        });
    };

    // Limiter les uploads à 5 en parallèle
    const uploadFilesInBatches = async (fileList, batchSize) => {
        const results = [];
        for (let i = 0; i < fileList.length; i += batchSize) {
            const batch = fileList.slice(i, i + batchSize);
            const uploadPromises = batch.map(uploadFile);

            try {
                const batchResults = await Promise.all(uploadPromises);
                results.push(...batchResults);
            } catch (error) {
                console.error('Error uploading files:', error);
            }
        }
        return results;
    };

    try {
        const uploadResults = await uploadFilesInBatches(files, MAX_PARALLEL_UPLOADS);
        console.log('All files uploaded:', uploadResults);
        loadFolderContent(); // Appelez votre fonction pour rafraîchir le contenu du dossier
    } catch (error) {
        console.error('Upload error:', error);
    }
});

// ---------------------- CONTEXT MENU ----------------------

// Variables pour gérer le menu contextuel
const contextMenu = document.getElementById('contextMenu');
let targetElement = null;

// Fonction pour afficher le menu contextuel
document.addEventListener('contextmenu', function (e) {
    // Si le clic droit est effectué sur un fichier
    const clickedElement = e.target.closest('.context-menu-target');
    if (clickedElement) {
        e.preventDefault();
        targetElement = e.target.closest('.context-menu-target');
        contextMenu.style.display = 'block';
        contextMenu.style.left = `${e.pageX}px`;
        contextMenu.style.top = `${e.pageY}px`;
        const name = clickedElement.getAttribute('data-info');
        const type = clickedElement.getAttribute('data-type');
        SelectedItem.set('Name', name);
        SelectedItem.set('Type', type);
        if (type !== "file") {
            document.getElementById('download').style.display = 'none'
        } else  {
            document.getElementById('download').style.display = ''
        }
    } else {
        contextMenu.style.display = 'none'; // Cacher si clic droit en dehors
    }
});

// Fonction pour cacher le menu contextuel quand on clique ailleurs
document.addEventListener('click', function () {
    contextMenu.style.display = 'none';
});

// ---------------------- FILE OPERATIONS ----------------------

export function downloadFile() {
    window.location.href=`/api/v1/files/?userId=${UserId}&path=${ActualPath}&filename=${encodeURIComponent(SelectedItem.get('Name'))}`;
}

export function deleteFile() {
    if (SelectedItem.get("Type") === "file") {
        fetch(`/api/v1/files/?userId=${UserId}&path=${ActualPath}&filename=${encodeURIComponent(SelectedItem.get('Name'))}`, {
            method: 'DELETE',
            credentials: 'include'
        })
            .then(() => loadFolderContent())
            .catch(error => console.error('Error deleting file:', error));
    } else {
        fetch(`/api/v1/files/folder?userId=${UserId}&path=${ActualPath}&folderName=${encodeURIComponent(SelectedItem.get('Name'))}`, {
            method: 'DELETE',
            credentials: 'include'
        })
            .then(() => loadFolderContent())
            .catch(error => console.error('Error deleting file:', error));
    }
}

export function renameFile(newName) {
    //let newName = prompt("Enter new name:", SelectedItem.get('Name'));
    if (newName) {
        fetch(`/api/v1/files/?userId=${UserId}&path=${ActualPath}&filename=${encodeURIComponent(SelectedItem.get('Name'))}&newName=${encodeURIComponent(newName)}`, {
            method: 'PATCH',
            credentials: 'include'
        })
            .then(() => loadFolderContent())
            .catch(error => console.error('Error renaming file:', error));
    }
}

export function createFolder(folderName) {
    if (folderName) {
        fetch(`/api/v1/files/folder?userId=${UserId}&path=${ActualPath}&folderName=${encodeURIComponent(folderName)}`, {
            method: 'POST',
            credentials: 'include'
        })
            .then(() => loadFolderContent())
            .catch(error => console.error('Error creating folder:', error));
    }
}

export function createShare(folderName) {
    if (folderName) {
        fetch(`/api/v1/share/folder?path=${ActualPath}&folderName=${encodeURIComponent(SelectedItem.get('Name'))}`, {
            method: 'POST',
            credentials: 'include'
        })
            .then(() => loadFolderContent())
            .catch(error => console.error('Error creating share:', error));
    }
}