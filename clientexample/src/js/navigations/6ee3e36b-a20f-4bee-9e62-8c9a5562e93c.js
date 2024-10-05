/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        if (!context.vms['containeriniciarsesion']) {
            context.top.active('containeriniciarsesion');
        }
        context.vms['containeriniciarsesion'].init();
    };
};
