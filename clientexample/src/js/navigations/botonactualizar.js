/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        var promise = context.actions['bf2aebf4-2069-4167-a611-bfd4dc8460f6']();
        context.runningActionsByContainer['containerwiki'].push(promise);
        promise.then(function (result) {
            context.runningActionsByContainer['containerwiki'].splice(
                context.runningActionsByContainer['containerwiki'].indexOf(promise), 1
            );
            if (result.event) {
                context.navigations[result.event](context, result.data);
            }
        });
    };
};
