export async function Prefrences(data) {
	console.info("Prefrences:- ", data);
	const url = '/setting/prefs';
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
	} catch (error) {
		console.error(error.message);
	}
}
