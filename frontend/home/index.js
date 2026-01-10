import {CloseModal, RefreshPage} from '../common/index.js';
import {Copy} from '../common/copy.js';
import {Prefrences} from './prefrences.js';

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
	selector.classList.add('bg-indigo-50');
	selector.querySelector('svg').classList.add('text-indigo-500');

	console.info('NavBar loaded...', navUrl[window.location.pathname]);
}

export async function SideNavCollapse() {
	const sideNav = document.querySelector('#side-navbar');
	const sideNavOpener = document.querySelector('#nav-bar-open');
	const topNavBar = document.querySelector('#top-nav-bar');
	if (!sideNav || !sideNavOpener || !topNavBar) {
		return;
	}
	const uiPrefs = JSON.parse(localStorage.getItem("ui-prefs") || "{}");
	if (uiPrefs?.storePrefsSideNav && uiPrefs?.storePrefsTopNavBar) {
		sideNav.className = uiPrefs.storePrefsSideNav;
		topNavBar.className = uiPrefs.storePrefsTopNavBar;
	}
	sideNavOpener.addEventListener('click', async (e) => {
		sideNav.classList.toggle('hidden');
		topNavBar.classList.toggle('ml-40');
		sideNavOpener.classList.toggle('text-indigo-500');

		localStorage.setItem("ui-prefs", JSON.stringify({
			storePrefsSideNav: sideNav.className,
			storePrefsTopNavBar: topNavBar.className
		}));

		const sideNavPref = {
			action: 'sideNav', value: sideNav.classList.contains('hidden') ? 'hide' : 'show'
		};
		await Prefrences(sideNavPref);
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
					console.log(selectedTag);
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
		RefreshPage('/');
	});

	CloseModal('#close-modal', linkModal);
}

async function AddNewLink(webLink, selectedTag) {
	const url = '/new';
	try {
		const headers = new Headers();
		headers.append('Content-Type', 'application/json');
		console.log(selectedTag);

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
			// RefreshPage();
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

	document.querySelectorAll('#copy-weblink').forEach((elm) => {
		elm.addEventListener('click', async (e) => {
			const svg = elm.querySelector('svg');
			const url = e.target.dataset?.url;
			if (svg) {
				svg.setAttribute('fill', 'black');
				setTimeout(() => svg.setAttribute('fill', 'none'), 700);
			}
			url && await Copy(e.target.dataset.url);
		});
	});

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
		document.getElementById('share-content').innerHTML = await response.text();
	} catch (error) {
		console.error(error.message);
	}
}

export function BackToTopBtn() {
	const btn = document.getElementById('backToTop');
	const footer = document.querySelector('footer');
	const pagination = document.querySelector('[data-pagination]');

	if (!btn || !footer || !pagination) {
		return;
	}

	const SHOW_AFTER = 200;
	const SAFE_OFFSET = 120;

	window.addEventListener('scroll', () => {
		const scrollY = window.scrollY;
		const windowH = window.innerHeight;

		let hide = scrollY < SHOW_AFTER;

		if (footer) {
			const footerTop = footer.getBoundingClientRect().top;
			if (footerTop < windowH + SAFE_OFFSET) hide = true;
		}

		if (pagination) {
			const paginationTop = pagination.getBoundingClientRect().top;
			if (paginationTop < windowH + SAFE_OFFSET) hide = true;
		}

		btn.classList.toggle('opacity-0', hide);
		btn.classList.toggle('pointer-events-none', hide);
	});

	btn.addEventListener('click', () =>
		window.scrollTo({top: 0, behavior: 'smooth'})
	);
}

export function AvatarToggle() {
	const menu = document.getElementById('avatar-menu');
	if (!menu) return;

	const btn = menu.querySelector('#avatar-btn');
	const dropdown = menu.querySelector('#avatar-dropdown');

	const open = () => {
		dropdown.classList.remove('opacity-0', 'scale-95', 'translate-y-1', 'pointer-events-none');
		dropdown.classList.add('opacity-100', 'scale-100', 'translate-y-0');
		btn.setAttribute('aria-expanded', 'true');
	};

	const close = () => {
		dropdown.classList.add('opacity-0', 'scale-95', 'translate-y-1', 'pointer-events-none');
		dropdown.classList.remove('opacity-100', 'scale-100', 'translate-y-0');
		btn.setAttribute('aria-expanded', 'false');
	};

	btn.addEventListener('click', (e) => {
		e.stopPropagation();
		dropdown.classList.contains('opacity-0') ? open() : close();
	});

	document.addEventListener('click', close);
	document.addEventListener('keydown', (e) => {
		if (e.key === 'Escape') close();
	});
}

export function SearchTextToggle() {
	const text = document.getElementById("search-text");
	const btn = document.getElementById("search-toggle");
	if (!text || !btn) return;
	let expanded = false;
	btn.addEventListener("click", () => {
		expanded = !expanded;
		if (expanded) {
			text.classList.remove("truncate");
			btn.textContent = "less";
		} else {
			text.classList.add("truncate");
			btn.textContent = "more";
		}
	});
}

export function UrlMenuToggle() {
	const menu = document.getElementById("url-menu");
	if (!menu) {
		return;
	}
	document.addEventListener('click', (e) => {
		document.querySelectorAll("details[open]").forEach((details) => {
			if (!details.contains(e.target)) {
				details.removeAttribute("open");
			}
		});
	});
}


