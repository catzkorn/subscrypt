package server

import "html/template"

const indexTemplate = `<!DOCTYPE html>
<html>
<style>
    table, td {
        border: 1px solid black;
    }
</style>

<body>
<h1>{{.PageTitle}}</h1>
<table style="width:100%">
    <tr>
        <td>Name Of Subscription</td>
        <td>Amount Due Monthly</td>
        <td>Date Due Monthly</td>
    </tr>
    {{range .Subscriptions}}
        <tr>
            <td>{{.Name}}</td>
            <td>{{.Amount}}</td>
            <td>{{.DateDue}}</td>
			<td><button type="button" id="delete-{{.ID}}" onclick="deleteSubscription({{.ID}})">Delete</button>
        </tr>
    {{end}}
</table>

<form action="/" method="post">
    <label for="name">Name Of Subscription:</label>
    <input type="text" name="name"><br>
    <label for="amount">Amount Due Monthly:</label>
    <input type="text" name="amount"><br>
    <label for="date">Date Due Monthly:</label>
    <input type="date" name="date"><br>
    <input type="submit" value="Submit">
</form>

<script>
	function deleteSubscription(id) {
		let xhttp = new XMLHttpRequest();
		let url = "/api/subscriptions/" + id
		xhttp.onreadystatechange = function () {
			if (xhttp.readyState === 4 && xhttp.status === 200) {
				window.location.href = "/";
			}
		}
		xhttp.open("DELETE", url, true);
		xhttp.send();
	}
</script>

</body>
</html>`

var ParsedIndexTemplate = template.Must(template.New("index").Parse(indexTemplate))
