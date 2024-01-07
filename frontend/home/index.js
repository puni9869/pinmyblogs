const navUrl = {
    '/': 'home',
    '/home': 'home',
    '/archived': 'archived',
    '/favourite': 'favourite'
};

export default function NavItemSelected() {
    let selector;
    if (!navUrl[window.location.pathname]) {
        return;
    }

    selector = window.location.pathname === '/'
        ? document.querySelector(`[data-home]`) :
        document.querySelector(`[data-${navUrl[window.location.pathname]}]`);
    if (!selector) {
        return;
    }
    selector.classList.add('text-orange-500');
    selector.querySelector('svg').classList.add('text-orange-500');

    console.log('NavBar loaded', navUrl[window.location.pathname]);
    console.log('NavBar selected item ', selector);
}
