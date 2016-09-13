import Vue from 'vue';
import html from './home.jade';


export default Vue.component('home', {
    name: 'home',
    template: html(),

    data() {
        return {
            longUrl: ''
        }
    },

    methods: {
        shorten() {
            this.$http.post('/create', {url: this.longUrl})
                .then(resp => resp.json())
                .then(json => this.$router.go(`/info/${json.Url}`));
        }
    }
});
