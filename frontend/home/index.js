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
	document.addEventListener('click', async (e) => {
		const action = e.target.closest('.action-btn');
		if (!action) return;

		const {id, url, isFav, isArchived, isDeleted} = action.dataset;
		/* COPY */
		if (action.classList.contains('copy-weblink')) {
			const svg = action.querySelector('svg');
			svg?.setAttribute('fill', 'black');
			setTimeout(() => svg?.setAttribute('fill', 'none'), 700);
			url && await Copy(url);
			return;
		}
		/* FAV / ARCHIVE / DELETE */
		const res = await WebLinkActions({id, isFav, isArchived, isDeleted});
		if (!res) {
			return;
		}
		if (Object.keys(action?.dataset)?.includes("isFav")) {
			const svg = action.querySelector('svg');
			if (isFav === "true") {
				svg?.setAttribute('fill', 'none');
				action.classList.remove("bg-red-50");
			} else {
				svg?.setAttribute('fill', 'currentColor');
				action.classList.add("bg-red-50");
			}

			const current = action.dataset["isFav"] === "true";
			const next = !current;
			action.dataset["isFav"] = next.toString();
		}
		const isArchivedOrDeleted = Object.keys(action?.dataset)?.includes("isArchived") || Object.keys(action?.dataset)?.includes("isDeleted")
		if (isArchivedOrDeleted) {
			const row = action.closest(".url-row");
			if (!row) return;
			// Lock current width (important for smooth animation)
			row.style.width = row.offsetWidth + "px";
			row.offsetHeight; // force reflow
			// Animate
			row.style.width = "0px";
			row.style.opacity = "0";
			// Remove after animation
			setTimeout(() => row.remove(), 300);
		}
		SelectAllCount();
	});
}

async function WebLinkActions(data) {
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
			return true;
		}
	} catch (error) {
		console.error(error.message);
		return false;
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

export function SelectUrlRows() {
	const rows = document.querySelectorAll("[data-select-id]");
	if (!rows.length) {
		return;
	}
	rows.forEach(row => {
		row.addEventListener("click", (e) => {
			// If click happened inside actions → ignore row click
			if (e.target.closest('#weblink-actions') || e.target.closest("#web-link")) {
				return;
			}
			row.classList.toggle('bg-blue-50');
			row.classList.toggle('bg-white');
			const checkBox = row.querySelector("#row-selected-checkbox");
			checkBox.checked = !checkBox.checked;
			SelectAllCount();
		});
	});
}

export function SelectAllCount() {
	const selectAll = document.querySelector("#select-all-label");
	if (!selectAll) {
		return;
	}

	const checkBoxes = document.querySelectorAll('#row-selected-checkbox:checked');
	selectAll.innerText = "Select all";
	if (!checkBoxes?.length) {
		document.querySelector("#url-menu").classList.add("hidden");
		return;
	}
	selectAll.innerText = `${checkBoxes.length} selected`;
	document.querySelector("#url-menu").classList.remove("hidden");
}

async function bulkActionApi(data) {
	const url = '/actions/bulk';
	try {
		const headers = new Headers();
		headers.append('Content-Type', 'application/json');
		const response = await fetch(url, {method: 'PATCH', headers: headers, body: JSON.stringify(data)});
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
		const resp = await response.json();
		if (resp?.Status === "OK") {
			console.info(resp?.Message);
			return true;
		}
	} catch (error) {
		return false;
	}
}

export function BulkAction() {
	["bulk-delete", "bulk-archive", "bulk-favourite"].forEach((a) => {
		const ele = document.querySelector(`#${a}`);
		if (!ele) {
			return;
		}
		ele.addEventListener("click", async () => {
			const items = document.querySelectorAll('#row-selected-checkbox:checked');
			const action = ele?.dataset?.action;
			if (!items || !items.length || !action) {
				return;
			}
			const ids = [];
			items.forEach((item) => {
				item?.dataset?.id && ids.push(item?.dataset?.id);
			});
			document.querySelectorAll("details[open]").forEach((details) => {
				details.removeAttribute("open");
			});
			const resp = await bulkActionApi({ids, action});
			if (resp) {
				console.log("Bulk action resp", resp);
				RefreshPage(window?.location.href);
			}
		});
	});
}

export function SelectAllBtn() {
	const btn = document.querySelector("#select-all-btn");
	if (!btn) {
		return;
	}
	btn.addEventListener("change", () => {
		const rows = document.querySelectorAll("[data-select-id]");
		if (!rows) {
			return;
		}
		rows.forEach(row => {
			row.classList.add('bg-blue-50');
			row.classList.add('bg-white');
			const checkBox = row.querySelector("#row-selected-checkbox");
			checkBox.checked = !checkBox.checked;
		});
		SelectAllCount();
	})
}



