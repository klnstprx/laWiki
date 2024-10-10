/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        var promise = context.actions['3fd91a41-e222-4d04-a0e1-e3e6883b5b1f']();
        context.runningActionsByContainer['containercrearentrada'].push(promise);
        promise.then(function (result) {
            context.runningActionsByContainer['containercrearentrada'].splice(
                context.runningActionsByContainer['containercrearentrada'].indexOf(promise), 1
            );
            if (result.event) {
                context.navigations[result.event](context, result.data);
            }
        });
    };
};
