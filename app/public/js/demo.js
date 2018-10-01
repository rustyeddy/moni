<script type="application/javascript"></script>
    var dotIndex = 0;
var graphviz = d3.select("#graph").graphviz()
    .transition(function () {
	return d3.transition("main")
	    .ease(d3.easeLinear)
	    .delay(500)
	    .duration(1500);
    })

    .logEvents(true)
    .on("initEnd", render);

function render() {
    var dotLines = dots[dotIndex];
    var dot = dotLines.join('');
    graphviz
	.renderDot(dot)
	.on("end", function () {
	    dotIndex = (dotIndex + 1) % dots.length;
	    render();
	});
}

var dots = [
    [
	'digraph  {',
	'    node [style="filled"]',
	'    a [fillcolor="#d62728"]',
	'    b [fillcolor="#1f77b4"]',
	'    a -> b',
	'}'
    ],
    [
	'digraph  {',
	'    node [style="filled"]',
	'    a [fillcolor="#d62728"]',
	'    c [fillcolor="#2ca02c"]',
	'    b [fillcolor="#1f77b4"]',
	'    a -> b',
	'    a -> c',
	'}'
    ],
    [
	'digraph  {',
	'    node [style="filled"]',
	'    a [fillcolor="#d62728"]',
	'    b [fillcolor="#1f77b4"]',
	'    c [fillcolor="#2ca02c"]',
	'    a -> b',
	'    a -> c',
	'}'
    ],
    [
	'digraph  {',
	'    node [style="filled"]',
	'    a [fillcolor="#d62728", shape="box"]',
	'    b [fillcolor="#1f77b4", shape="parallelogram"]',
	'    c [fillcolor="#2ca02c", shape="pentagon"]',
	'    a -> b',
	'    a -> c',
	'    b -> c',
	'}'
    ],
    [
	'digraph  {',
	'    node [style="filled"]',
	'    a [fillcolor="yellow", shape="star"]',
	'    b [fillcolor="yellow", shape="star"]',
	'    c [fillcolor="yellow", shape="star"]',
	'    a -> b',
	'    a -> c',
	'    b -> c',
	'}'
    ],
    [
	'digraph  {',
	'    node [style="filled"]',
	'    a [fillcolor="#d62728", shape="triangle"]',
	'    b [fillcolor="#1f77b4", shape="diamond"]',
	'    c [fillcolor="#2ca02c", shape="trapezium"]',
	'    a -> b',
	'    a -> c',
	'    b -> c',
	'}'
    ],
    [
	'digraph  {',
	'    node [style="filled"]',
	'    a [fillcolor="#d62728", shape="box"]',
	'    b [fillcolor="#1f77b4", shape="parallelogram"]',
	'    c [fillcolor="#2ca02c", shape="pentagon"]',
	'    a -> b',
	'    a -> c',
	'    b -> c',
	'}'
    ],
    [
	'digraph  {',
	'    node [style="filled"]',
	'    a [fillcolor="#d62728"]',
	'    b [fillcolor="#1f77b4"]',
	'    c [fillcolor="#2ca02c"]',
	'    a -> b',
	'    a -> c',
	'    c -> b',
	'}'
    ],
    [
	'digraph  {',
	'    node [style="filled"]',
	'    b [fillcolor="#1f77b4"]',
	'    c [fillcolor="#2ca02c"]',
	'    c -> b',
	'}'
    ],
    [
	'digraph  {',
	'    node [style="filled"]',
	'    b [fillcolor="#1f77b4"]',
	'}'
    ],
];

</script>
