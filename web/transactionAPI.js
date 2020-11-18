function getTransactions() {
    let xhttp = new XMLHttpRequest();
    let url = "/api/transactions/"
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState === 4 && xhttp.status === 200) {
            window.location.href = "/transactions";
        }
    }
    xhttp.open("GET", url, true);
    xhttp.send();
}