// Copyright 2017 The vizmon-demo Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

const indexTmpl = `
<html>
	<head>
		<title>SoLiD sensors monitoring</title>
		<script type="text/javascript">
		var sock = null;

		function update(data) {
			var p = null;
			
			p = document.getElementById("update-message");
			p.innerHTML = "Last Update: <code>"+data.update+"</code>";

			p = document.getElementById("mon-plot");
			p.innerHTML = data.plot;

			p = document.getElementById("mon-data");
			p.innerHTML = "<pre>"+data.data+"</pre>";
		};

		window.onload = function() {
			sock = new WebSocket("ws://"+location.host+"/data");
			sock.onmessage = function(event) {
				var data = JSON.parse(event.data);
				update(data);
			};
		};
		</script>

		<style>
		.vizmon-plot-style {
			font-size: 14px;
			line-height: 1.2em;
		}
		</style>
	</head>

	<body>
		<h2>VizMon demo monitoring plots ({{.Freq}} Hz)</h2>

		<div id="mon-plots">
			<div id="mon-plot" class="vizmon-plot-style"></div>
		</div>

		<br>
		<div id="update-message">Last Update: N/A
		</div>
		<div id="mon-data"></div>
	</body>
</html>
`
