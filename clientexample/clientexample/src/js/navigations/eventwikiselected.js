/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context, data) {
        if (!context.vms['containerwiki']) {
            context.top.active('containerwiki');
            context.vms['containerwiki'].init({mask: 'componentewiki'});
        }
        data = data || {};
        var packet = {
            'id' : data['id']
        };
        context.vms['componentewiki'].init({input: packet});
    };
};
