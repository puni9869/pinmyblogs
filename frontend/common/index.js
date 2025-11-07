import {StringBuilder} from './stringBuilder.js';

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

function jsonToCsv(jsonData) {
	const csv = new StringBuilder();
	let headers = Object.keys(jsonData[0]);
	csv.write(headers.join(',') + '\n');
	jsonData.forEach(function (row) {
		let data = headers.map(header => row[header]).join(',');
		csv.write(data + '\n');
	});
	return csv.toString();
}

export function SaveFile(filename, dataObjToWrite, format) {
	const blob = new Blob([
			format === 'json' ? JSON.stringify(dataObjToWrite) : jsonToCsv(dataObjToWrite)
		],
		{
			type: `text/${format}`
		});

	const link = document.createElement("a");
	link.download = filename;
	link.href = window.URL.createObjectURL(blob);
	link.dataset.downloadurl = [`text/${format}`, link.download, link.href].join(':');

	const evt = new MouseEvent("click", {view: window, bubbles: true, cancelable: true,});
	link.dispatchEvent(evt);
	link.remove();
}
