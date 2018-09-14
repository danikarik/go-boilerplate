require('./bootstrap');

window.Vue = require('vue');

Vue.component('home-component', require('./components/HomeComponent.vue'));

const app = new Vue({
    el: '#app'
});
