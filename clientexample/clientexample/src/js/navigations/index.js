/*jslint node: true, nomen: true */
"use strict";

exports.createNavigations = function (options) {
    return {
        'eventactualizarentrada': require('./eventactualizarentrada').createNavigation(options)
        ,'eventcrearentrada': require('./eventcrearentrada').createNavigation(options)
        ,'eventeliminarentrada': require('./eventeliminarentrada').createNavigation(options)
        ,'eventappendcomentario': require('./eventappendcomentario').createNavigation(options)
        ,'eventappendwiki': require('./eventappendwiki').createNavigation(options)
        ,'eventactualizarwiki': require('./eventactualizarwiki').createNavigation(options)
        ,'buttoncrearentrada': require('./buttoncrearentrada').createNavigation(options)
        ,'buttoncrear': require('./buttoncrear').createNavigation(options)
        ,'buttonactualizarentrada': require('./buttonactualizarentrada').createNavigation(options)
        ,'buttoneliminarentrada': require('./buttoneliminarentrada').createNavigation(options)
        ,'buttonnuevawiki': require('./buttonnuevawiki').createNavigation(options)
        ,'eventiniciarsesion': require('./eventiniciarsesion').createNavigation(options)
        ,'eventsesioniniciada': require('./eventsesioniniciada').createNavigation(options)
        ,'buttonactualizarwiki': require('./buttonactualizarwiki').createNavigation(options)
        ,'buttoneliminarwiki': require('./buttoneliminarwiki').createNavigation(options)
        ,'botonnuevaentrada': require('./botonnuevaentrada').createNavigation(options)
        ,'eventeliminarwiki': require('./eventeliminarwiki').createNavigation(options)
        ,'buttonnuevocomentario': require('./buttonnuevocomentario').createNavigation(options)
        ,'eventwikiselected': require('./eventwikiselected').createNavigation(options)
        ,'evententradaselected': require('./evententradaselected').createNavigation(options)
    };
};
