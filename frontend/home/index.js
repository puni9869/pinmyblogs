import {CloseModal, RefreshPage} from '../common/index.js';

const navUrl = {
	'/': 'home',
	'/home': 'home',
	'/archived': 'archived',
	'/favourite': 'favourite',
	'/setting': 'setting',
	'/trash': 'trash',
	'/login': 'login'
};

export function NavItemSelected() {
	if (!navUrl[window.location.pathname]) {
		return;
	}

	const selector = window.location.pathname === '/'
		? document.querySelector('[data-home]')
		: document.querySelector(`[data-${navUrl[window.location.pathname]}]`);
	if (!selector) {
		return;
	}
	selector.classList.add('text-indigo-500');
	selector.querySelector('svg').classList.add('text-indigo-500');
	console.info('NavBar loaded...', navUrl[window.location.pathname]);
	console.info('NavBar selected item ', selector);
}

export function SideNavCollapse() {
	const sideNav = document.querySelector('#side-navbar');
	const sideNavOpener = document.querySelector('#nav-bar-open');
	if (!sideNav || !sideNavOpener) {
		return;
	}
	sideNavOpener.addEventListener('click', (e) => {
		sideNav.classList.toggle('hidden');
		sideNavOpener.classList.toggle('text-indigo-500');
	});
}

export function AddNewWebLinkInit() {
	const link = document.querySelector('#add-weblink');
	const linkModal = document.querySelector('#add-weblink-modal');
	const saveBtn = document.querySelector('#add-link-btn');
	if (!link || !linkModal || !saveBtn) {
		return;
	}
	const tags = document.querySelectorAll("#tag");
	let selectedTag;
	link.addEventListener('click', async () => {
		linkModal.classList.toggle('hidden');
		if (tags.length) {
			tags.forEach((tag) => {
				tag.addEventListener('change', (e) => {
					selectedTag = e.target.value;
				});
			});
		}
	});
	saveBtn.addEventListener('click', async () => {
		let url = document.querySelector('#add-weblink-input-box');
		if (!url || url?.value.trim().length === 0) {
			return;
		}
		await AddNewLink(url?.value.trim(), selectedTag ? selectedTag : 'black');
		linkModal.classList.toggle('hidden');
		url.value = '';
	});

	CloseModal('#close-modal', linkModal);
}

async function AddNewLink(webLink, selectedTag) {
	const url = '/new';
	try {
		const headers = new Headers();
		headers.append('Content-Type', 'application/json');

		const response = await fetch(url, {
			method: 'POST',
			headers: headers,
			body: JSON.stringify({url: webLink, tag: selectedTag})
		});
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
		const resp = await response.json();
		if (resp?.Status === "OK") {
			console.info(resp?.Message);
			RefreshPage();
		}
	} catch (error) {
		console.error(error.message);
	}
}
