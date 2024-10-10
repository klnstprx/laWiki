/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        if (!context.vms['containercrearwiki']) {
            context.top.active('containercrearwiki');
        }
        context.vms['containercrearwiki'].init();
    };
};
