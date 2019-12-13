var vm = new Vue({
    el: '#app',
    data: {
	title: "Page",
	message: "Hello",
	links: [],
	url: "",
	elapsed: "",
    },
    computed: {
	messageNow: function(val) {
	    axios.post('/api/crawl/rustyeddy.com', {
		url: 'rustyeddy.com'
	    }).then((response) => {
		vm.message = (response.data);
		vm.url = (response.data.URL);
		vm.elasped = response.data.Elapsed;
		vm.links = response.data.Links;
	    }, (error) => {
		console.log(error);
	    });
	},
    }
})

