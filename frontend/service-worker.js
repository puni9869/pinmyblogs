const CACHE_NAME = "pinmyblogs-v1";
const OFFLINE_URL = "/offline.html";

/* Install */
self.addEventListener("install", (event) => {
	event.waitUntil(
		caches.open(CACHE_NAME).then((cache) => cache.add(OFFLINE_URL))
	);
	self.skipWaiting();
});

/* Activate */
self.addEventListener("activate", (event) => {
	event.waitUntil(self.clients.claim());
});

/* Fetch handling */
self.addEventListener("fetch", (event) => {
	event.respondWith(
		fetch(event.request).catch(() => {
			if (event.request.mode === "navigate") {
				return caches.match(OFFLINE_URL);
			}
		})
	);
});

/* Reload page when internet comes back */
self.addEventListener("message", (event) => {
	if (event.data === "ONLINE") {
		self.clients.matchAll({type: "window"}).then((clients) => {
			clients.forEach((client) => {
				client.navigate(client.url);
			});
		});
	}
});
