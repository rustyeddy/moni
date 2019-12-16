var vm = new Vue({
    el: '#app',
    data: {
	title: "Page",
	message: "Hello",
	page: {
	    url: "",
	    links: [],
	    elapsed: "",
	}
	sites: [],
	domains: []
    },
    mounted() {
	axios.post('/api/crawl/', {
	    url: vm.url
	}).then((response) => {
	    vm.message = (response.data);
	    vm.page.url = (response.data.URL);
	    vm.page.elasped = response.data.Elapsed;
	    vm.page.links = response.data.Links;
	}, (error) => {
	    console.log(error);
	});
    }
})

