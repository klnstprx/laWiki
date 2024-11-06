/*jslint node: true, nomen: true */
"use strict";

var ko = require('knockout'),
    Promise = require('bluebird');

function ViewModel(params) {
    var self = this;
    self.context = params.context;
    self.status = ko.observable('');
    self.fields = ko.observable({});
    self.errors = ko.observable({});

    self.trigger = function (id) {
        self.context.navigations[id](self.context, self.output);
    };
}

ViewModel.prototype.id = 'componenteformwiki';

ViewModel.prototype.waitForStatusChange = function () {
    return this._initializing ||
           Promise.resolve();
};

ViewModel.prototype._compute = function () {
    this.output = {
        'desc': this.input['desc'],
        'genero': this.input['genero'],
        'nombre': this.input['nombre'],
        'url': this.input['url'],
    }
    var self = this,
        fields = {
            'desc': ko.observable(this.input['desc']),
            'genero': ko.observable(this.input['genero']),
            'nombre': ko.observable(this.input['nombre']),
            'url': ko.observable(this.input['url']),
        },
        errors = {
            'desc': ko.observable(this.input['desc-error']),
            'genero': ko.observable(this.input['genero-error']),
            'nombre': ko.observable(this.input['nombre-error']),
            'url': ko.observable(this.input['url-error']),
        };
    fields['desc'].subscribe(function (value) {
        self.output['desc'] = value;
        self.errors()['desc'](undefined);
    });
    fields['genero'].subscribe(function (value) {
        self.output['genero'] = value;
        self.errors()['genero'](undefined);
    });
    fields['nombre'].subscribe(function (value) {
        self.output['nombre'] = value;
        self.errors()['nombre'](undefined);
    });
    fields['url'].subscribe(function (value) {
        self.output['url'] = value;
        self.errors()['url'](undefined);
    });
    this.fields(fields);
    this.errors(errors);
    this.status('computed');
};


ViewModel.prototype.init = function (options) {
    options = options || {};
    this.output = undefined;
    this.fields({});
    this.errors({});
    this.input = options.input || {};
    this.status('ready');
    var self = this;
    this._initializing = new Promise(function (resolve) {
        setTimeout(function () {
            self._compute();
            resolve();
            self._initializing = undefined;
        }, 1);
    });
};

exports.register = function () {
    ko.components.register('c-componenteformwiki', {
        viewModel: {
            createViewModel: function (params, componentInfo) {
                var vm = new ViewModel(params);
                params.context.vms[vm.id] = vm;
                ko.utils.domNodeDisposal.addDisposeCallback(componentInfo.element, function () { delete params.context.vms[vm.id]; });
                return vm;
            }
        },
        template: require('./index.html'),
        synchronous: true
    });
};
