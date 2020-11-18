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
    userHTML = _newUserForm();
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
  return `
<button type="button" id="edit-user-button" onclick="showEditUserForm('${user.Name}', '${user.Email}')">
    <svg class="w-6 h-6"fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z">
</path>
</svg>
          </button>
`;
}

function _newUserForm() {
  let newUserForm = "<div class=\"card w-75\">" +
                      "<div class=\"card-body\">" +
                        "<h5 class=\"card-title\">Welcome! Enter your details</h5>" +
                        _formatUserForm() +
                      "</div>" +
                    "</div>"
  return newUserForm
}

function _formatUserForm() {

  let userForm = `
    <form>
        <div class="row">
            <div class="col">
                <input type="text" class="form-control" id="username" placeholder="Name">
            </div>
            <div class="col">
                <input type="text" class="form-control" id="email" placeholder="Email">
            </div>
            <div class="col">
                <button type="button" class="btn btn-primary" id="create-user-button" onclick="createUser()">Submit</button>
            </div>
        </div>
    <!--  <label for="name">Name:</label>-->
    <!--  <input type="text" id="username"><br>-->
    <!--  <label for="email">Email:</label>-->
    <!--  <input type="text" id="email"><br>-->
    <!--  <button type="button" id="create-user-button" onclick="createUser()">Submit</button>-->
    </form>
`;

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