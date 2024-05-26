import {NavItemSelected, SideNavCollapse} from './home/index.js';
import {DeleteAccountModelInit} from './setting/index.js';

document.addEventListener('DOMContentLoaded', function () {
    NavItemSelected();
    SideNavCollapse();
    DeleteAccountModelInit();
    console.log("App is loaded");
}, false);
