/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        if (!context.vms['containerindice']) {
            context.top.active('containerindice');
        }
        context.vms['containerindice'].init();
    };
};
