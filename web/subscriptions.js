loadSubscriptions()

function loadSubscriptions() {
    _getSubscriptions(_showSubscriptions)
}

function _getSubscriptions(callback) {
    let xhttp = new XMLHttpRequest();
    let path = '/api/subscriptions';
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState === 4 && xhttp.status === 200) {
            let subscriptions = convertToSubscriptions(xhttp.responseText);
            callback(subscriptions);
        }
    };
    xhttp.open("GET", path, true);
    xhttp.send();
}

function _showSubscriptions(subscriptions) {
    let subscriptionsHTML = ""
    if (subscriptions.length > 0) {
        subscriptionsHTML = _formatSubscriptionsTable(subscriptions)
    } else {
        subscriptionsHTML = "You don't have any subscriptions"
    }
    document.getElementById("subscriptions").innerHTML = subscriptionsHTML
}

function _formatSubscriptionsTable(subscriptions) {
    let tableHTML = `<table id=\"table-subscriptions\" style=\"width:100%\">
                                <tr>
                                    <td>Subscription Name</td>
                                    <td>Amount</td>
                                    <td>Next Payment Date</td>
                                </tr>`;

    subscriptions.forEach(function (subscription) {
        tableHTML += _formatSubscription(subscription);
    })

    tableHTML += "</table>"
    return tableHTML
}

function _formatSubscription(subscription) {
    return `<tr>
            <td>${subscription.name}</td>
            <td>${subscription.amount}</td>
            <td>${subscription.date}</td>
            <td><button type="button" id="reminder-${subscription.id}" onclick="sendReminder(${subscription.id})">Reminder</button></td>
            <td><button type="button" id="delete-${subscription.id}" onclick="deleteSubscription(${subscription.id})">Delete</button></td>
            </tr>`
}

function createSubscription() {
    let name = document.getElementById('subscription-name').value;
    let amount = document.getElementById('subscription-amount').value;
    let dateDue = formatDateForJSON(document.getElementById('subscription-date').value);

    if (validateSubscriptionValues(name, amount, dateDue) !== false) {
        _postSubscription(name, amount, dateDue)
    }
}


function _postSubscription(name, amount, dateDue) {
    let xhttp = new XMLHttpRequest();
    let url = "/api/subscriptions";
    xhttp.open("POST", url, true);
    xhttp.setRequestHeader("Content-type", "application/json");
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState == 4 && xhttp.status == 200) {
            loadSubscriptions();
        }
    }
    let data = JSON.stringify({"name": name, "amount": amount, "dateDue": dateDue});
    xhttp.send(data);
}

function validateSubscriptionValues(name, amount, dateDue) {
    if (name === "" || amount === "" || dateDue === "") {
        document.getElementById("subscription-error").innerHTML = "Please enter subscription details";
        return false
    } else {
        return true
    }
}

function convertToSubscriptions(res) {
    let resSubscriptions = JSON.parse(res);
    let subscriptions = [];
    resSubscriptions.forEach(function (subscription) {
        let subscriptionObj = new Subscription(subscription.id, subscription.name, subscription.amount, subscription.dateDue);
        subscriptions.push(subscriptionObj);
    })
    return subscriptions;
}

function deleteSubscription(id) {
    let xhttp = new XMLHttpRequest();
    let url = "/api/subscriptions/" + id
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState === 4 && xhttp.status === 200) {
            loadSubscriptions();
        }
    }
    xhttp.open("DELETE", url, true);
    xhttp.send();
}

function formatDateForJSON(date) {
    return date + "T00:00:00Z"
}