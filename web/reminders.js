function sendReminder(id) {
  let xhttp = new XMLHttpRequest();
  let url = "/api/reminders";
  let data = JSON.stringify({ "id": id });

  xhttp.onreadystatechange = function() {
    if (xhttp.readyState === 4 && xhttp.status === 200) {
      window.location.href = "/";
    }
  };
  xhttp.open("POST", url, true);
  xhttp.send(data);
}