/*jslint node: true, nomen: true */
"use strict";

exports.createNavigation = function () { // add "options" parameter if needed
    return function (context) {
        var promise = context.actions['appendanuncio']();
        context.runningActionsByContainer['containercrearwiki'].push(promise);
        promise.then(function (result) {
            context.runningActionsByContainer['containercrearwiki'].splice(
                context.runningActionsByContainer['containercrearwiki'].indexOf(promise), 1
            );
            if (result.event) {
                context.navigations[result.event](context, result.data);
            }
        });
    };
};
