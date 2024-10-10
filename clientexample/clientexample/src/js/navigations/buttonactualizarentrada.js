/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        var promise = context.actions['1b4a9ceb-a03a-4d00-9d08-e6a765e2485e']();
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
