function importTransactionsToSubscriptions() {
    let xhttp = new XMLHttpRequest();
    let url = "/api/transactions/load-subscriptions"
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState === 4 && xhttp.status === 200) {
            loadSubscriptions();
        }
    }
    xhttp.open("POST", url, true);
    xhttp.send();
}