loadUser();

function loadUser() {
  _getUser(_showUser);
}

function _getUser(callback) {
  let xhttp = new XMLHttpRequest();
  let path = '/api/users';
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState === 4 && xhttp.status === 200) {
      let user = JSON.parse(xhttp.responseText);
      console.log(xhttp.responseText);
      callback(user);
    }
  };
  xhttp.open("GET", path, true);
  xhttp.send();
}

function _showUser(user) {
  let userHTML = "";
  if (user != null) {
    userHTML = _formatUser(user);
    userHTML += _formatEditUserButton(user)

  } else {
    userHTML = _formatUserForm();
  }
  document.getElementById("user").innerHTML = userHTML;
}

function _formatUser(user) {
  return `
  Name: ${user.Name}<br>
  Email: ${user.Email}<br>
  `;
}

function _formatEditUserButton(user) {
  return `<button type="button" id="edit-user-button" onclick="showEditUserForm('${user.Name}', '${user.Email}')">Edit</button>`;
}

function _formatUserForm() {

  let userForm = `<form>
  <label for="name">Name:</label>
  <input type="text" id="username"><br>
  <label for="email">Email:</label>
  <input type="text" id="email"><br>
  <button type="button" id="create-user-button" onclick="createUser()">Submit</button>
</form>`;

  return userForm;
}

function showEditUserForm(name, email) {
  let editUserFormHTML = _formatUserForm();
  document.getElementById("user").innerHTML = editUserFormHTML;
  document.getElementById("username").value = name;
  document.getElementById("email").value = email;
}


// THIS WORKS BELLOW LEAVE IT ALONE! VVVV


function createUser() {
  let name = document.getElementById('username').value;
  let email = document.getElementById('email').value;

  if (validateUserValues(name, email) == false) {
    return;
  }

  let xhttp = new XMLHttpRequest();
  let url = "/api/users";
  xhttp.open("POST", url, true);
  xhttp.setRequestHeader("Content-type", "application/json");
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState == 4 && xhttp.status == 200) {
      window.location.href = "/";
    }
  };
  let data = JSON.stringify({ "name": name, "email": email });
  xhttp.send(data);
}

function validateUserValues(name, email) {
  if (name == "" || email == "") {
    document.getElementById("user-error").innerHTML = "Please enter user details";
    return false;
  } else {
    return true;
  }
}