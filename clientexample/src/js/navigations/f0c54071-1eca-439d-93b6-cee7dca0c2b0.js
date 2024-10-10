/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        if (!context.vms['containerwiki']) {
            context.top.active('containerwiki');
        }
        context.vms['containerwiki'].init();
    };
};
