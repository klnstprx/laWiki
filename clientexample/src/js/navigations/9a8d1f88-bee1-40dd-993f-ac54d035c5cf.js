/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        var promise = context.actions['5e9b2fbc-eb5a-413f-8e16-4d8791b49dbc']();
        context.runningActionsByContainer['containerentrada'].push(promise);
        promise.then(function (result) {
            context.runningActionsByContainer['containerentrada'].splice(
                context.runningActionsByContainer['containerentrada'].indexOf(promise), 1
            );
            if (result.event) {
                context.navigations[result.event](context, result.data);
            }
        });
    };
};
