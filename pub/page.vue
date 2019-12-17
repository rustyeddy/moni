<template>

    <class="card">
	<div class="card-body">
	    <h4>{{ url }}</h4>
	    <h6>{{ elapsed }}</h6>
	    <ul class="link-group">
		<li class="link-group-item" v-for="link in links">
		    <a clas="link-group-link" :href="link">{{ link }}</a>
		</li>
	    </ul>
	</div>
    </class>

</template>

<script>
 export default {
     el: '#page',
     data: {
	 url: "",
	 elapsed: "",
	 links: []
     },
     mounted() {
	 /*
	 axios.post('/api/crawl/', {
	     vm.url: "rustyeddy.com"
	 }).then((response) => {
	     vm.message = (response.data);
	     vm.page.url = (response.data.URL);
	     vm.page.elasped = response.data.Elapsed;
	     vm.page.links = response.data.Links;
	 }, (error) => {
	     console.log(error);
	 });
	 */
	 console.log("page has been mounted");
     }
 };
</script>

<style scoped>
 .row {
     margin: 2rem;
 }
</style>
