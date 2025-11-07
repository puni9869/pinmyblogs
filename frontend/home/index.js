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
	const topNavBar = document.querySelector('#top-nav-bar');
	if (!sideNav || !sideNavOpener || !topNavBar) {
		return;
	}
	sideNavOpener.addEventListener('click', (e) => {
		sideNav.classList.toggle('hidden');
		topNavBar.classList.toggle('ml-40');
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
		const errEl = document.getElementById('weblink-err');
		if (errEl) {
			errEl.classList.toggle("hidden");
			errEl.innerText = error.message;
		}
		console.error(error.message);
	}
}

export function WebLinkActionsInit() {
	const actions = document.querySelector('#weblink-actions');
	if (!actions) {
		return;
	}

	['#move-to-favourite', '#move-to-archive', '#move-to-trash'].forEach((actionSelector) => {
		document.querySelectorAll(actionSelector).forEach((elm) => {
			elm.addEventListener('click', async (e) => {
				await WebLinkActions(e.target.dataset);
				RefreshPage();
			});
		});
	});
}


async function WebLinkActions(data) {
	console.info("Action is called", data);
	const url = '/actions';
	try {
		const headers = new Headers();
		headers.append('Content-Type', 'application/json');

		const response = await fetch(url, {
			method: 'PUT',
			headers: headers,
			body: JSON.stringify(data)
		});
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
		const resp = await response.json();
		if (resp?.Status === "OK") {
			console.info(resp?.Message);
		}
	} catch (error) {
		console.error(error.message);
	}
}

export function ShareLinkInit() {
	document.querySelectorAll('#share').forEach((elm) => {
		elm.addEventListener('click', async (e) => (await ShareLink(e.target.dataset)));
	});
}

async function ShareLink(data) {
	if (!data?.id) {
		return;
	}
	console.info("ShareLink is called", data);
	const url = `/share/${data.id}`;
	try {
		const response = await fetch(url, {method: 'GET'});
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
		const resp = await response.text();
		document.getElementById('share-content').innerHTML = resp;
	} catch (error) {
		console.error(error.message);
	}
}

