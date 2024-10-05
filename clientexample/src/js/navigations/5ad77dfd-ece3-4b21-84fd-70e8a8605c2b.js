/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        if (!context.vms['containerentrada']) {
            context.top.active('containerentrada');
        }
        context.vms['containerentrada'].init();
    };
};
