// loadSubscriptions();

function loadSubscriptions() {
    _getSubscriptions(_showSubscriptions);
}

function createSubscription() {
    let name = document.getElementById('subscription-name').value;
    let amount = document.getElementById('subscription-amount').value;
    let dateDue = _formatDateForJSON(document.getElementById('subscription-date').value);

    if (_validateSubscriptionValues(name, amount, dateDue) !== false) {
        _postSubscription(name, amount, dateDue);
    }
}

function deleteSubscription(id) {
    let xhttp = new XMLHttpRequest();
    let url = "/api/subscriptions/" + id;
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState === 4 && xhttp.status === 200) {
            loadSubscriptions();
        }
    };
    xhttp.open("DELETE", url, true);
    xhttp.send();
}

function _getSubscriptions(callback) {
    let xhttp = new XMLHttpRequest();
    let path = '/api/subscriptions';
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState === 4 && xhttp.status === 200) {
            let subscriptions = _convertToSubscriptions(xhttp.responseText);
            callback(subscriptions);
        }
    };
    xhttp.open("GET", path, true);
    xhttp.send();
}

function _showSubscriptions(subscriptions) {
    let subscriptionsHTML = "";
    if (subscriptions.length > 0) {
        subscriptionsHTML = _formatSubscriptionsTable(subscriptions);
    } else {
        subscriptionsHTML = "You don't have any subscriptions";
    }
    document.getElementById("subscriptions-table").innerHTML = subscriptionsHTML;
}

function _formatSubscriptionsTable(subscriptions) {
    let tableHTML = `<table class="table" id=\"table-subscriptions\" style=\"width:100%\">
                        <thead>
                            <tr>
                                <th scope="col">Subscription Name</th>
                                <th scope="col">Amount</th>
                                <th scope="col">Payment Date</th>
                                <th scope="col">Frequency</th>
                            </tr>
                        </thead>
                        <tbody>`;

    subscriptions.forEach(function (subscription) {
        tableHTML += _formatSubscription(subscription);
    });

    tableHTML += "</tbody></table>";
    return tableHTML;
}

function _formatSubscription(subscription) {
    return `<tr>
            <th scope="row">${subscription.name}</th>
            <td>${_formatAmountTwoDecimals(subscription.amount)}</td>
            <td>${_formatDateAsDay(subscription.dateDue)}</td>
            <td>Monthly</td>
            <td><button type="button" class="btn btn-primary" id="reminder-${subscription.id}" onclick="sendReminder(${subscription.id})">Reminder</button></td>
            <td><button type="button" class="btn btn-primary" id="delete-${subscription.id}" onclick="deleteSubscription(${subscription.id})">Delete</button></td>
            </tr>`;
}

function _formatAmountTwoDecimals(amount) {
    return parseFloat(amount).toFixed(2);
}

function _formatDateAsDay(date) {
    let d = date.getDate();
    return d + _getOrdinal(d);
}

function _getOrdinal(number) {
    if (number > 3 && number < 21) return 'th';
    switch (number % 10) {
        case 1:
            return "st";
        case 2:
            return "nd";
        case 3:
            return "rd";
        default:
            return "th";
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
            document.getElementById("create-subscription-form").reset();
        }
    };
    let data = JSON.stringify({"name": name, "amount": amount, "dateDue": dateDue});
    xhttp.send(data);
}

function _validateSubscriptionValues(name, amount, dateDue) {
    if (name === "" || amount === "" || dateDue === "") {
        document.getElementById("subscription-error").innerHTML = "Please enter subscription details";
        return false;
    } else {
        return true;
    }
}

function _convertToSubscriptions(res) {
    let resSubscriptions = JSON.parse(res);
    let subscriptions = [];
    if (resSubscriptions === null) {
        return subscriptions;
    } else {
        resSubscriptions.forEach(function (subscription) {
            let subscriptionObj = new Subscription(subscription.id, subscription.name, subscription.amount, subscription.dateDue);
            subscriptions.push(subscriptionObj);
        });
        return subscriptions;
    }
}

function _formatDateForJSON(date) {
    return date + "T00:00:00Z";
}