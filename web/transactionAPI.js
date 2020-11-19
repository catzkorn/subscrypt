function importTransactionsToSubscriptions() {
    showSpinner();
    let xhttp = new XMLHttpRequest();
    let url = "/api/transactions/load-subscriptions"
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState === 4 && xhttp.status === 200) {
            loadSubscriptions();
            hideSpinner()
        }
    }
    xhttp.open("POST", url, true);
    xhttp.send();
}