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
