import {
	NavItemSelected,
	SideNavCollapse,
	AddNewWebLinkInit,
	WebLinkActionsInit,
	ShareLinkInit,
	BackToTopBtn,
	AvatarToggle,
	SearchTextToggle,
	UrlMenuToggle,
	SelectUrlRows,
	SelectAllCount
} from './home/index.js';

import {DeleteAccountModelInit, DisableMyAccountInit, DownloadMyData} from './setting/index.js';

document.addEventListener('DOMContentLoaded', async function () {
	NavItemSelected();
	await SideNavCollapse();
	DeleteAccountModelInit();
	DisableMyAccountInit();
	AddNewWebLinkInit();
	WebLinkActionsInit();
	ShareLinkInit();
	DownloadMyData();
	BackToTopBtn();
	AvatarToggle();
	SearchTextToggle();
	UrlMenuToggle();
	SelectUrlRows();
	SelectAllCount()
	console.info('App is loaded');
}, false);


if ("serviceWorker" in navigator) {
	window.addEventListener("load", () => {
		navigator.serviceWorker.register("/statics/service-worker.js");
	});
}

window.addEventListener("online", async () => {
	if (navigator.serviceWorker.controller) {
		navigator.serviceWorker.controller.postMessage("ONLINE");
	} else {
		location.reload();
	}
});

