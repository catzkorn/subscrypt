function sendReminder(id) {

  showReminderToast();

  let name = document.getElementById('user-name').value;
  let email = document.getElementById('email').value;

  if (validateUsersInformation(name, email) == false) {
    return;
  }

  let xhttp = new XMLHttpRequest();
  let url = "/api/reminders";
  let data = JSON.stringify({ "id": id });

  xhttp.onreadystatechange = function() {
    if (xhttp.readyState === 4 && xhttp.status === 200) {
      // showReminderToast()
    }
  };
  xhttp.open("POST", url, true);
  xhttp.send(data);
}




function validateUsersInformation(name, email) {
  if (name == "" || email == "") {
    document.getElementById("reminder-error").innerHTML = "Please enter user details to receive a reminder";
    return false;
  } else {
    return true;
  }
}