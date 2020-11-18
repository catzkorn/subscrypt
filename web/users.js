loadUser();

// OUTLINE ICONS
// const pencilSvg = `<svg class="w-6 h-6"fill="none" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
//                         <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z">
//                         </path>
//                     </svg>`

const atSvg = `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207">
                    </path>
                </svg>`

const userSvg =   `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z">
                    </path>
                  </svg>`


// SOLID ICONS

const pencilSvg = `<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"></path></svg>`

// const userSvg = `<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"></path></svg>`
//
// const atSvg = `<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M14.243 5.757a6 6 0 10-.986 9.284 1 1 0 111.087 1.678A8 8 0 1118 10a3 3 0 01-4.8 2.401A4 4 0 1114 10a1 1 0 102 0c0-1.537-.586-3.07-1.757-4.243zM12 10a2 2 0 10-4 0 2 2 0 004 0z" clip-rule="evenodd"></path></svg>`

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
    loadSubscriptions();
  } else {
    userHTML = _newUserForm();
  }
  document.getElementById("user").innerHTML = userHTML;
}

function _formatUser(user) {
  return `
      <form>
        <div class="form-group row">
            <label for="name" class="col- col-form-label col-form-label-sm"><span class="icon">${userSvg}</span></label>
            <div class="col-3">
                <input type="text" readonly class="form-control-plaintext form-control-sm" id="username" value="${user.Name}">
            </div>
              <label for="email" class="col- col-form-label col-form-label-sm"><span class="icon">${atSvg}</span></label>
              <div class="col-3">
                  <input type="text" readonly class="form-control-plaintext form-control-sm" id="email" value="${user.Email}">
              </div>
              <div class="col-5">
              ${_formatEditUserButton(user)}
</div>
        </div>
    </form>
  `;
}



function _formatEditUserButton(user) {
  return `<button type="button" class="icon-button" id="edit-user-button" onclick="showEditUserForm('${user.Name}', '${user.Email}')">` +
          pencilSvg +
          `</button>`;
}

function _newUserForm() {
  let newUserForm = `<div class="card w-110 mx-auto" id='new-user-form'>` +
                      `<div class="card-body">` +
                        `<h5 class="card-title text-center">Welcome! Enter your details</h5>` +
                        _formatUserForm() +
                      "</div>" +
                    "</div>"
  return newUserForm
}

function _formatUserForm(type) {
  let centerClass = ""
  if (type === "new") {
    centerClass = "justify-content-center"
  }
  let userForm = `
    <form>
        <div class="form-group row ${centerClass}">
            <label for="name" class="col- col-form-label col-form-label-sm"><span class="icon">${userSvg}</span></label>
            <div class="col-3">
                <input type="text" class="form-control form-control-sm" id="username" placeholder="Name">
            </div>
            <label for="email" class="col- col-form-label col-form-label-sm"><span class="icon">${atSvg}</span></label>
            <div class="col-3">
                <input type="text" class="form-control form-control-sm" id="email" placeholder="Email">
            </div>
            <div class="col-5">
                <button type="button" class="btn btn-primary" id="create-user-button" onclick="createUser()">Submit</button>
            </div>
        </div>
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
      loadUser();
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