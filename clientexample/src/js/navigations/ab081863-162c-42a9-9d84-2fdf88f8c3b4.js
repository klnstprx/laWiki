/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context, data) {
        if (!context.vms['containerentrada']) {
            context.top.active('containerentrada');
            context.vms['containerentrada'].init({mask: 'componentedetallesentrada'});
        }
        data = data || {};
        var packet = {
            'id' : data['id']
        };
        context.vms['componentedetallesentrada'].init({input: packet});
    };
};
