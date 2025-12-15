/* =========================================================
   pinmyblogs.com - Service Worker (MV3)
   ========================================================= */

const LOGIN_URL = "http://127.0.0.1/login";
const SAVE_URL = "http://127.0.0.1/new";

/**
 * Get active tab
 */
async function getActiveTab() {
	const tabs = await chrome.tabs.query({
		active: true,
		currentWindow: true
	});
	return tabs[0] || null;
}

/**
 * Open login page in new tab
 */
async function openLoginTab() {
	await chrome.tabs.create({
		url: LOGIN_URL,
		active: true
	});
}

/**
 * Save URL to pinmyblogs.com
 */
async function saveToPinMyBlogs(url) {
	const response = await fetch(SAVE_URL, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			"Accept": "application/json"
		},
		body: JSON.stringify({
			url,
			tag: "black"
		}),
		credentials: "include"
	});

	// 🚨 Not logged in
	if (response.status === 401 || response.status === 403) {
		await openLoginTab();
		throw new Error("NOT_AUTHENTICATED");
	}

	if (!response.ok) {
		const text = await response.text();
		throw new Error(`HTTP ${response.status}: ${text}`);
	}

	const data = await response.json();

	// Optional JSON-based auth check
	if (data?.authenticated === false) {
		await openLoginTab();
		throw new Error("NOT_AUTHENTICATED");
	}

	return data;
}

/**
 * Click on extension icon
 */
chrome.action.onClicked.addListener(async (tab) => {
	try {
		const activeTab = tab?.url ? tab : await getActiveTab();

		if (!activeTab?.url) {
			console.error("❌ No active tab URL found");
			return;
		}

		const url = activeTab.url;

		// Block internal browser pages
		if (
			url.startsWith("chrome://") ||
			url.startsWith("edge://") ||
			url.startsWith("about:")
		) {
			console.warn("⚠️ Internal page cannot be saved:", url);
			return;
		}

		console.log("📌 Attempting to save:", url);

		const result = await saveToPinMyBlogs(url);

		console.log("✅ Saved successfully:", result);

	} catch (error) {
		if (error.message === "NOT_AUTHENTICATED") {
			console.warn("🔐 User not logged in → redirecting to login");
		} else {
			console.error("❌ Save failed:", error.message);
		}
	}
});

/**
 * Install hook
 */
chrome.runtime.onInstalled.addListener(() => {
	console.log("📌 pinmyblogs extension installed");
});
