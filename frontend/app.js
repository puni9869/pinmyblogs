import {NavItemSelected, SideNavCollapse} from './home/index.js';
import {DeleteAccountModelInit, DisableMyAccountInit} from './setting/index.js';

document.addEventListener('DOMContentLoaded', function () {
    NavItemSelected();
    SideNavCollapse();
    DeleteAccountModelInit();
    DisableMyAccountInit();
    console.log('App is loaded');
}, false);
