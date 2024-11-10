export function RedirectToLogin() {
	window.location = '/logout';
}

export function RefreshPage(path = '') {
	window.location = path;
}

export function CloseModal(modalId, toClose) {
	const m = document.querySelector(modalId);
	m.addEventListener('click', () => {
		toClose.classList.toggle('hidden');
	});
}

export function JsonToCsv(jsonData) {
	let csv = '';
	// Get the headers
	let headers = Object.keys(jsonData[0]);
	csv += headers.join(',') + '\n';

	// Add the data
	jsonData.forEach(function (row) {
		let data = headers.map(header => row[header]).join(',');
		csv += data + '\n';
	});
	return csv;
}

export function SaveFile(filename, dataObjToWrite, format) {
	const blob = new Blob([
			format === 'json' ? JSON.stringify(dataObjToWrite) : JsonToCsv(dataObjToWrite)
		],
		{
			type: `text/${format}`
		});
	const link = document.createElement("a");

	link.download = filename;
	link.href = window.URL.createObjectURL(blob);
	link.dataset.downloadurl = [`text/${format}`, link.download, link.href].join(':');

	const evt = new MouseEvent("click", {
		view: window,
		bubbles: true,
		cancelable: true,
	});

	link.dispatchEvent(evt);
	link.remove()
}
