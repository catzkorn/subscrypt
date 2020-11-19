// FILLED ICONS
const calendarSvg = `<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clip-rule="evenodd"></path></svg>`
const binSvg = `<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd"></path></svg>`

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

function showSpinner() {
    let spinner = document.getElementById("loading-spinner");
    spinner.style.display = "block";
}

function hideSpinner() {
    let spinner = document.getElementById("loading-spinner");
    spinner.style.display = "none";
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
    let subscriptionsHTML = `<h4>Subscriptions</h4>`;
    if (subscriptions.length > 0) {
        subscriptionsHTML += _formatSubscriptionsTable(subscriptions);
    } else {
        subscriptionsHTML += "<p>You don't have any subscriptions</p>";
    }
    document.getElementById("subscriptions-table").innerHTML = subscriptionsHTML;
    let subscriptionDiv = document.getElementById("subscriptions")
    subscriptionDiv.style.display = "block";
}

function _formatSubscriptionsTable(subscriptions) {
    let tableHTML = `<table class="table" id="table-subscriptions">
                        <thead>
                            <tr>
                                <th scope="col">Subscription Name</th>
                                <th scope="col">Amount</th>
                                <th scope="col">Payment Date</th>
                                <th scope="col">Frequency</th>
                                <th scope="col">Actions</th>
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
            <td><button type="button" class="icon-button" id="reminder-button" onclick="sendReminder(${subscription.id})">${calendarSvg}</button>
           <button type="button" class="icon-button" id="delete-${subscription.id}" onclick="deleteSubscription(${subscription.id})">${binSvg}</button></td>
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