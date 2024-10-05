/*jslint node: true, nomen: true */
"use strict";

exports.createRepositories = function (options) {
    var repositories = {}
    repositories['Entrada'] = require('./Entrada').createRepository(options);
    repositories['Comentario'] = require('./Comentario').createRepository(options);
    repositories['Wiki'] = require('./Wiki').createRepository(options);
    return repositories;
};
