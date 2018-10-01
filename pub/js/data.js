var width = 420,
    barheight = 20;

var x = d3.scale.linear()
    .range([0, width]);

var chart = d3.select(".chart")
    .attr("width", width);


d3.tsv("data.tsv", type, function(error, data) {

    x.domain([0, d3.max(data, function(d) { return d.value; })]);
    chart.attr("height", barheight *data.length);

    var bar = chart.selectAll("g")
	.data(data)
	.enter().append("g")
	.attr("transform", function(d, i) { return "translate(0, " + i * barheight + ")"; });

    bar.append("rect")
	.attr("width", function(d) { return x(d.value); })
	.attr("height", barheight - 1);

    bar.append("text")
	.attr("x", function(d) { return x(d.value) -3; })
	.attr("y", barheight / 2)
	.attr("dy", ".35em")
	.text(function(d) { return d.value; });
});

function type(d) {
    d.value = +d.value;
    return d
}

