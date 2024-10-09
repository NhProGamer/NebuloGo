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