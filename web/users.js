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
    userHTML += _formatEditUserButton(user);
    loadSubscriptions();
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
  return `<button type="button" class="btn btn-default" id="edit-user-button" onclick="showEditUserForm('${user.Name}', '${user.Email}')">
            <i class="fas fa-pen"></i>
          </button>`;
}

function _formatUserForm() {

  let userForm = `<form>
    <div class="row">
        <div class="col">
            <input type="text" class="form-control" id="username" placeholder="Name">
        </div>
        <div class="col">
            <input type="text" class="form-control" id="email" placeholder="Email">
        </div>
        <div class="col">
            <button type="button" class="btn btn-primary" id="create-user-button" onclick="createUser()"><i class="fas fa-pen"></i></button>
        </div>
    </div>
<!--  <label for="name">Name:</label>-->
<!--  <input type="text" id="username"><br>-->
<!--  <label for="email">Email:</label>-->
<!--  <input type="text" id="email"><br>-->
<!--  <button type="button" id="create-user-button" onclick="createUser()">Submit</button>-->
</form>`;

  return userForm;
}

function showEditUserForm(name, email) {
  let editUserFormHTML = _formatUserForm();
  document.getElementById("user").innerHTML = editUserFormHTML;
  document.getElementById("username").value = name;
  document.getElementById("email").value = email;
}


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