import React from "react";
import ReactDOM from "react-dom";
import BankModal from "./banking/bankModal";
import Navbar from "./navbar/navbar";
import SubscriptionsModal from "./subscriptions/subscriptionsModal";
import SubscriptionsTable from "./subscriptions/subscriptionsTable";
import UserLogin from "./users/userLogin";

ReactDOM.render(<App />, document.getElementById("root"));

function App(props) {
  return (
    <>
      <Navbar />

      <UserLogin />

      <SubscriptionsTable />

      <SubscriptionsModal />

      <BankModal />
    </>
  );
}
