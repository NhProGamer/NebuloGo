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

export function openPopupInfo(title, content) {
    const popup = document.getElementById('warningPopupOverlay');
    popup.querySelector('#infoPopupTitle').innerHTML = title;
    popup.querySelector('#infoPopupContent').innerHTML = content;
    popup.classList.remove('hidden');

    document.getElementById('warningOkPopupButton').onclick = function() {
        closePopupWarning()
    };
}

export function showSharePopup() {
    const sharePopupOverlay = document.getElementById("sharePopupOverlay");
    return new Promise((resolve, reject) => {
        sharePopupOverlay.classList.remove("hidden");
        document.getElementById('confirmSharePopupButton').onclick = function() {
            const shareDate = document.getElementById("shareDate").value;
            const isPublic = document.getElementById("isPublic").value;
            hideSharePopup()
            resolve([shareDate, isPublic])
        };
        document.getElementById('cancelSharePopupButton').onclick = function() {
            hideSharePopup();
            reject()
        };
    })
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

function closePopupInfo() {
    const popup = document.getElementById('infoPopupOverlay');
    popup.classList.add('hidden');
}

function hideSharePopup() {
    const popup = document.getElementById("sharePopupOverlay");
    popup.classList.add("hidden");
    document.getElementById("shareDate").value = ""
    document.getElementById("isPublic").value = ""
}