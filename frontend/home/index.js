const navUrl = {
  '/': 'home',
  '/home': 'home',
  '/archived': 'archived',
  '/favourite': 'favourite',
  '/setting': 'setting',
  '/trash': 'trash',
  '/login': 'login'
};

export function NavItemSelected() {
  if (!navUrl[window.location.pathname]) {
    return;
  }

  const selector = window.location.pathname === '/'
    ? document.querySelector('[data-home]')
    : document.querySelector(`[data-${navUrl[window.location.pathname]}]`);
  if (!selector) {
    return;
  }
  selector.classList.add('text-orange-500');
  selector.querySelector('svg').classList.add('text-orange-500');
  console.log('NavBar loaded...', navUrl[window.location.pathname]);
  console.log('NavBar selected item ', selector);
}

export function SideNavCollapse() {
  const sideNav = document.querySelector('#side-navbar');
  const sideNavOpener = document.querySelector('#nav-bar-open');
  if (!sideNav || !sideNavOpener) {
    return;
  }
  sideNavOpener.addEventListener('click', (e) => {
    sideNav.classList.toggle('hidden');
    sideNavOpener.classList.toggle('text-orange-500');
  });
}
