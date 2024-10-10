/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        var promise = context.actions['fc018cbe-0455-4a1b-9d3f-deab97b44706']();
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
