Vue.component('site-list', {
    template: `
	<ul class="list-group">
		<site v-for="site in sites">{{ site }}</site>
	</ul>
	`,
    data() {
	axios.get('/api/sites')
	    .then((response) => {
		this.sites = response.data.Sites;
	    }, (error) => {
		console.log(error);
	    });

	return {
	    sites: [
		{ site: "rustyeddy.com", up: false },
		{ site: "sierrahydrographics.com", up: false },		
	    ]
	}
    },
})

Vue.component('site', {
    template: '<li class="list-group-item"><a href="<slot></slot>"><slot></slot></a></li>',
});

Vue.component('alert', {
    props: ['level'],
    template: '<div class="alert" :class="level" role="alert"><slot></slot></div>'
});

new Vue({
    el: '#app',
    data: {
	title: "The Title"
    }
});

