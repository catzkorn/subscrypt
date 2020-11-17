function createSubscription() {
    let name = document.getElementById('subscription-name').value;
    let amount  = document.getElementById('subscription-amount').value;
    let dateDue  = formatDate(document.getElementById('subscription-date').value);

    let xhttp = new XMLHttpRequest();
    let url = "/api/subscriptions";
    xhttp.open("POST", url, true);
    xhttp.setRequestHeader("Content-type", "application/json");
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState == 4 && xhttp.status == 200) {
            window.location.href = "/";
        }
    }
    let data = JSON.stringify({"name": name, "amount": amount, "dateDue": dateDue});
    xhttp.send(data);
}

function getSubscriptionFormValues() {

}

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

function formatDate(date) {
    return date + "T00:00:00Z"
}