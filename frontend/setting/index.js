import {RedirectToLogin} from '../common/index.js';

export function DeleteAccountModelInit() {
    const deleteBtn = document.querySelector('#delete-my-account');
    if (!deleteBtn) {
        return;
    }
    const deleteModal = document.querySelector('#delete-account-modal');
    deleteBtn.addEventListener('click', (e) => {
        deleteModal.classList.toggle('hidden');
    });
}

async function DisableAccount() {
    const url = '/setting/disablemyaccount';
    try {
        const response = await fetch(url, {method: 'PUT'});
        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error(error.message);
    }
}

export function DisableMyAccountInit() {
    const disableBtn = document.querySelector('#disable-my-account');
    if (!disableBtn) {
        return;
    }
    disableBtn.addEventListener('click', async (e) => {
        console.log('Disable my account event');
        const res = await DisableAccount();
        if (res && res['Status'] === 'OK') {
            console.log('Account is disabled until you logged in back.');
            RedirectToLogin();
        }
    });
}
