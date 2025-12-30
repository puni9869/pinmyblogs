const CACHE_NAME = "pinmyblogs-v1";
const OFFLINE_URL = "/offline.html";

/* Install */
self.addEventListener("install", async (event) => {
	event.waitUntil(
		caches.open(CACHE_NAME).then((cache) =>
			cache.addAll([OFFLINE_URL]).catch(() => {
			})
		)
	);
	await self.skipWaiting();
});

/* Activate */
self.addEventListener("activate", (event) => {
	event.waitUntil(self.clients.claim());
});

/* Fetch */
self.addEventListener("fetch", (event) => {
	if (event.request.mode !== "navigate") return;

	event.respondWith(
		fetch(event.request).catch(async () => {
			const cached = await caches.match(OFFLINE_URL);
			return cached || new Response("Offline", {status: 503});
		})
	);
});
