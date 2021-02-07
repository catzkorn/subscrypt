import React from "react";

function Navbar(props) {
  return (
    <nav className="navbar navbar-expand-sm navbar-dark bg-transparent">
      <span className="navbar-brand">Subscrypt</span>
      <ul className="navbar-nav mr-auto">
        <li className="nav-item active">
          <a className="nav-link" href="/">
            Home<span className="sr-only">(current)</span>
          </a>
        </li>
        <li className="nav-item">
          <a className="nav-link" href="/transactions">
            Transactions
          </a>
        </li>
      </ul>
    </nav>
  );
}

export default Navbar;
