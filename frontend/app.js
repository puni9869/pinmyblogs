import {
	NavItemSelected,
	SideNavCollapse,
	AddNewWebLinkInit,
	WebLinkActionsInit,
	ShareLinkInit,
	BackToTopBtn,
	AvatarToggle
} from './home/index.js';

import {DeleteAccountModelInit, DisableMyAccountInit, DownloadMyData} from './setting/index.js';

document.addEventListener('DOMContentLoaded', function () {
	NavItemSelected();
	SideNavCollapse();
	DeleteAccountModelInit();
	DisableMyAccountInit();
	AddNewWebLinkInit();
	WebLinkActionsInit();
	ShareLinkInit();
	DownloadMyData();
	BackToTopBtn();
	AvatarToggle();
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

