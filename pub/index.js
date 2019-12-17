Vue.component('site-list', {
    template: `
	<div class="list-group">
		<site v-for="site in sites">{{ site.site }}</site>
	</div>
	`,
    data() {
	return {
	    sites: [
		{ site: "oclowvision.com", up: false },
		{ site: "rustyeddy.com", up: false },
		{ site: "sierrahydrographics.com", up: false },		
	    ]
	}
    }
})

Vue.component('site', {
    template: '<li class="list-item"><slot></slot></li>',
});

Vue.component('alert', {
    props: ['level'],
    template: '<div class="alert" :class="level" role="alert"><slot></slot></div>'
});

new Vue({
    el: '#app',
});

