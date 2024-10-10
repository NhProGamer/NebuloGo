import * as jwt from './jwt.js'
import * as api from './api.js'

export let ActualPath = new URLSearchParams(window.location.search).get('path') ?? '';

const jsonWebToken = jwt.getCookie('token');
export const UserId = jwt.parseJwt(jsonWebToken).user_id;

export let SelectedItem = new Map();
SelectedItem.set('Name', "");

export function loadFolderContent() {
    api.getShares().then(data => document.getElementById('share-list').innerHTML = renderFolderContent(data))
}

// Generate folder and file HTML
function renderFolderContent(data) {
    if (!data || data.length === 0) {
        return '';  // Ne rien faire si data est vide
    }
    return data.map(item => {
        const parsedDate = new Date(item.Expiration);

        // Extraire les composants de la date
        const day = String(parsedDate.getDate()).padStart(2, '0'); // Jour
        const month = String(parsedDate.getMonth() + 1).padStart(2, '0'); // Mois (0-11, donc +1)
        const year = String(parsedDate.getFullYear()); // Année (deux derniers chiffres)
        const hours = String(parsedDate.getHours()).padStart(2, '0'); // Heures
        const minutes = String(parsedDate.getMinutes()).padStart(2, '0'); // Minutes

            return `
                <tr class="hover:bg-gray-100 context-menu-target" data-info="${item.InternalID}">
                    <td class="px-4 py-2">${item.FilePath}</td>
                    <td class="px-4 py-2">${item.Owner}</td>
                    <td class="px-4 py-2">${item.Public}</td>
                    <td class="px-4 py-2">${day}/${month}/${year} ${hours}:${minutes}</td>
                </tr>`
    }).join('');
}

// Handle browser navigation using the back button
window.addEventListener('popstate', event => {
    ActualPath = (event.state && event.state.path !== null) ? event.state.path : "";
    loadFolderContent();
});


// Variables pour gérer le menu contextuel
export const contextMenu = document.getElementById('contextMenu');
export let targetElement = null;

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
        SelectedItem.set('Name', name);

    } else {
        contextMenu.style.display = 'none'; // Cacher si clic droit en dehors
        targetElement = null
    }
});

// Fonction pour cacher le menu contextuel quand on clique ailleurs
document.addEventListener('click', function () {
    contextMenu.style.display = 'none';
    targetElement = null
});