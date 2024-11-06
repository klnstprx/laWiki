/*jslint node: true, nomen: true */
"use strict";

var ko = require('knockout');

exports.register = function () {
    require('./main-application').register();
    require('./c-containercrearentrada').register();
    require('./c-containercrearwiki').register();
    require('./c-containerentrada').register();
    require('./c-containerindice').register();
    require('./c-containeriniciarsesion').register();
    require('./c-containerwiki').register();
    require('./c-componentecrearentrada').register();
    require('./c-componenteformwiki').register();
    require('./c-componentedetallesentrada').register();
    require('./c-componentelistacomentarios').register();
    require('./c-componenteaadircomentario').register();
    require('./c-componentelistawikis').register();
    require('./c-componentewiki').register();
    require('./c-componentelistaentradas').register();
                                        };
