<template>
    <div class="row">
	<div class="col">
	    <h1>Page Vue</h1>
	</div>
    </div>
</template>

<script>
 export default {
     el: '#app',
     data: {
	 message: "bye"
     },
     methods: {
	 handleClick() {
	     this.count++;
	     console.log("clicked", this.count);
	 },
	 message() {
	     return "Yeppers!";
	 }
     },
 };
</script>
<style scoped>
 .row {
     margin: 2rem;
 }
</style>
