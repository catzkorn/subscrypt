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