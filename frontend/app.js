import { NavItemSelected, SideNavCollapse, AddNewWebLinkInit,WebLinkActionsInit } from './home/index.js';
import { DeleteAccountModelInit, DisableMyAccountInit } from './setting/index.js';

document.addEventListener('DOMContentLoaded', function () {
  NavItemSelected();
  SideNavCollapse();
  DeleteAccountModelInit();
  DisableMyAccountInit();
  AddNewWebLinkInit();
  WebLinkActionsInit();
  console.info('App is loaded');
}, false);
