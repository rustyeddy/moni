var vm = new Vue({
    el: '#app',
    data: {
	message: "Hello",
	title: "HelloPage"
    },
    computed: {
	messageNow: function(val) {
	    axios.post('/api/crawl/rustyeddy.com', {
		url: 'rustyeddy.com'
	    }).then((response) => {
		vm.message = (response.data);
	    }, (error) => {
		console.log(error);
	    });
	},
    }
})

