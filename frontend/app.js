import { NavItemSelected, SideNavCollapse, AddNewWebLinkInit } from './home/index.js';
import { DeleteAccountModelInit, DisableMyAccountInit } from './setting/index.js';

document.addEventListener('DOMContentLoaded', function () {
  NavItemSelected();
  SideNavCollapse();
  DeleteAccountModelInit();
  DisableMyAccountInit();
  AddNewWebLinkInit();
  console.info('App is loaded');
}, false);
