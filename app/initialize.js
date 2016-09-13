import Vue from 'vue';
import VueRouter from 'vue-router';
import VueResource from 'vue-resource';

import { home, info } from './components';


Vue.use(VueRouter);
Vue.use(VueResource);

const environment = process.env.NODE_ENV;

Vue.config.debug = (environment === 'development');
Vue.config.devtools = (environment === 'development');

let app = {};

let router = new VueRouter({
    mode: 'hash'
});

router.map({
//    '*':{
//        name: 'error',
//        component: components.error
//    },

    '/': {
        component: home
    },

    '/info/:url': {
        name: 'info',
        component: info
    }
});

router.start(app, '#content', () => {
    if (window) {
        window.App = router.app;
    }
});
