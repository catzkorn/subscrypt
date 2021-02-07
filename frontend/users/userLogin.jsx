import React from "react";

function UserLogin(props) {
  return (
    <div className="container user-header pl-4">
      <div className="row new-user" id="new-user">
        {/* <!-- New user form gets set here --> */}
      </div>
      <div className="row">
        <div className="col-10">
          <div id="existing-user"></div>
        </div>
      </div>
    </div>
  );
}

export default UserLogin;
