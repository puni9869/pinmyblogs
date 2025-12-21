import {RedirectToLogin, SaveFile} from '../common/index.js';

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
	const url = '/setting/disable-my-account';
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
		console.info('Disable my account event');
		const res = await DisableAccount();
		if (res && res?.Status === 'OK') {
			console.info('Account is disabled until you logged in back.');
			RedirectToLogin();
		}
	});
}

export function DownloadMyData() {
	const download = document.querySelectorAll('#download');
	if (!download.length) {
		return;
	}
	download.forEach(d => {
		d.addEventListener('click', async (e) => {
			if (!e.target.dataset?.format.length) {
				return
			}
			const url = `/setting/download-my-data/${e.target.dataset?.format}`;
			try {
				const headers = new Headers();
				headers.append('Content-Type', 'application/json');

				const response = await fetch(url, {method: 'GET', headers: headers});
				if (!response.ok) {
					throw new Error(`Response status: ${response.status}`);
				}
				const jData = await response.json();
				SaveFile(`pinmyblogs.${e.target.dataset?.format}`, jData, e.target.dataset?.format);
			} catch (error) {
				console.error(error.message);
			}
		});
	});
}
