/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        if (!context.vms['containercrearentrada']) {
            context.top.active('containercrearentrada');
        }
        context.vms['containercrearentrada'].init();
    };
};
