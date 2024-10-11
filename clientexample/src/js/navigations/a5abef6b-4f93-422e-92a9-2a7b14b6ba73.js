/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        var promise = context.actions['5d73955f-0d01-44d3-8dbd-5e84189d93be']();
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
