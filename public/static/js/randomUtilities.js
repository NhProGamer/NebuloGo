export function formatBytes(bytes) {
    const units = ['o', 'Ko', 'Mo', 'Go', 'To', 'Po', 'Exa']; // Ajoutez d'autres unités si nécessaire
    let index = 0;

    while (bytes >= 1024 && index < units.length - 1) {
        bytes /= 1024;
        index++;
    }

    return `${bytes.toFixed(2)} ${units[index]}`; // Format avec 2 décimales
}