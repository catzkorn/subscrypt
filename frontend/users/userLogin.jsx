import React, { useEffect, useState } from "react";

const atSvg = (
  <svg
    className="w-6 h-6"
    fill="none"
    stroke="currentColor"
    viewBox="0 0 24 24"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"
    ></path>
  </svg>
);

const userSvg = (
  <svg
    className="w-6 h-6"
    fill="none"
    stroke="currentColor"
    viewBox="0 0 24 24"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
    ></path>
  </svg>
);

const pencilSvg = (
  <svg
    className="w-6 h-6"
    fill="currentColor"
    viewBox="0 0 20 20"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"></path>
  </svg>
);

function UserDetails(props) {
  const [userName, setUserName] = useState("");
  const [userEmail, setUserEmail] = useState("");
  const [editUser, setEditUser] = useState(false);

  useEffect(() => {
    const url = "/api/users";
    fetch(url)
      .then((response) => {
        if (response.status !== 200) {
          console.log("error getting users");
          return;
        }
        return response.json();
      })
      .then((payload) => {
        props.setUser(payload);
      });
  }, []);

  function handleUserNameChange(event) {
    setUserName(event.target.value);
  }

  function handleUserEmailChange(event) {
    setUserEmail(event.target.value);
  }

  function submitUserDetails(event, userName, userEmail) {
    event.preventDefault();
    const url = "/api/users";
    options = {
      method: "POST",
      headers: { "Content-type": "application/json" },
      body: JSON.stringify({
        name: userName,
        email: userEmail,
      }),
    };

    fetch(url, options).then((response) => {
      if (response.status !== 200) {
        console.log("There was an error submitting user details", response);
        return;
      }
      props.setUser({ Name: userName, Email: userEmail });
    });
  }

  function _newUserForm() {
    return (
      <div className="card mx-auto justify-content-center" id="new-user-form">
        <div className="card-body">
          <h5 className="card-title text-center new-user pb-3">
            Welcome! Enter your details
          </h5>
          <form className="row" id="user-form">
            <div className="form-group row justify-content-center new-user">
              <label
                htmlFor="user-name"
                className="col- col-form-label col-form-label-md ml-4"
              >
                <span className="icon icon-light">{userSvg}</span>
              </label>
              <div className="col">
                <input
                  type="text"
                  className="form-control form-control-md input-dark"
                  id="user-name"
                  placeholder="Name"
                  onChange={(event) => handleUserNameChange(event)}
                  value={userName}
                ></input>
              </div>
              <label
                htmlFor="email"
                className="col- col-form-label col-form-label-md ml-4"
              >
                <span className="icon icon-light">{atSvg}</span>
              </label>
              <div className="col">
                <input
                  type="email"
                  className="form-control form-control-md input-dark"
                  id="email"
                  placeholder="Email"
                  onChange={(event) => {
                    handleUserEmailChange(event);
                  }}
                  value={userEmail}
                ></input>
              </div>
              <div className="col-2">
                <button
                  type="button"
                  className="btn btn-primary mr-5"
                  id="create-user-button"
                  onClick={(event) =>
                    submitUserDetails(event, userName, userEmail)
                  }
                >
                  Submit
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    );
  }

  function _existingUser() {
    if (editUser) {
      return (
        <form className="form-inline" id="user-form">
          <div className="form-group">
            <label
              htmlFor="user-name"
              className="col-form-label col-form-label-md"
            >
              <span className="icon icon-light">{userSvg}</span>
            </label>
            <input
              type="text"
              className="form-control form-control-md input-active input-light"
              id="user-name"
              onChange={(event) => {
                props.setUser({
                  Name: event.target.value,
                  Email: props.user.Email,
                });
              }}
              value={props.user.Name}
            ></input>
            <label htmlFor="email" className="col-form-label col-form-label-md">
              <span className="icon icon-light">{atSvg}</span>
            </label>
            <input
              type="text"
              className="form-control form-control-md input-active input-light"
              id="email"
              onChange={(event) => {
                props.setUser({
                  Name: props.user.Name,
                  Email: event.target.value,
                });
              }}
              value={props.user.Email}
            ></input>
            <button
              type="button"
              className="btn btn-primary"
              id="create-user-button"
              onClick={(event) => {
                submitUserDetails(event, props.user.Name, props.user.Email);
                setEditUser(false);
              }}
            >
              Submit
            </button>
          </div>
        </form>
      );
    }

    return (
      <>
        <h3>Welcome back</h3>
        <form className="form-inline">
          <div className="form-group">
            <label
              htmlFor="user-name"
              className="col-form-label col-form-label-mdt"
            >
              <span className="icon icon-light">{userSvg}</span>
            </label>
            <input
              type="text"
              readOnly
              className="form-control-plaintext form-control-md input-light"
              id="user-name"
              value={props.user.Name}
            ></input>
            <label htmlFor="email" className="col-form-label col-form-label-md">
              <span className="icon icon-light">{atSvg}</span>
            </label>
            <input
              type="email"
              readOnly
              className="form-control-plaintext form-control-md input-light"
              id="email"
              value={props.user.Email}
            ></input>
            <button
              type="button"
              className="icon-button icon-light"
              id="edit-user-button"
              onClick={() => setEditUser(true)}
            >
              {pencilSvg}
            </button>
          </div>
        </form>
      </>
    );
  }

  let userForm;

  if (props.user === null) {
    userForm = _newUserForm();
  } else {
    userForm = _existingUser();
  }

  return <div className="user">{userForm}</div>;
}

export default UserDetails;
