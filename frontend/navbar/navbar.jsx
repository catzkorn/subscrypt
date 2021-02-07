import React from "react";
import { Link } from "react-router-dom";

function Navbar(props) {
  return (
    <nav className="navbar navbar-expand-sm navbar-dark bg-transparent">
      <span className="navbar-brand">Subscrypt</span>
      <ul className="navbar-nav mr-auto">
        <li className="nav-item active">
          <Link className="nav-link" to="/">
            Home<span className="sr-only">(current)</span>
          </Link>
        </li>
        <li className="nav-item">
          <Link className="nav-link" to="/transactions">
            Transactions
          </Link>
        </li>
      </ul>
    </nav>
  );
}

export default Navbar;
