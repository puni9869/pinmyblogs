export function DeleteAccountModelInit() {
    const deleteBtn = document.querySelector('#delete-my-account');
    if (!deleteBtn) {
        return;
    }
    const deleteMadal = document.querySelector('#delete-account-modal');
    deleteBtn.addEventListener('click', (e) => {
        deleteMadal.classList.toggle('hidden');
    });
}