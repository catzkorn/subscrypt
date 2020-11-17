function createSubscription() {
    let name = document.getElementById('subscription-name').value;
    let amount = document.getElementById('subscription-amount').value;
    let dateDue = formatDate(document.getElementById('subscription-date').value);

    if (validateSubscriptionValues(name, amount, dateDue) == false) {
        return
    }

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

function validateSubscriptionValues(name, amount, dateDue) {
    if (name == "" || amount == "" || dateDue == "") {
        document.getElementById("subscription-error").innerHTML="Please enter subscription details";
        return false
    } else {
        return true
    }
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