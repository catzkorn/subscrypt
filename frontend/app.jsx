import React, { useState } from "react";
import ReactDOM from "react-dom";
import BankModal from "./banking/bankModal";
import Navbar from "./navbar/navbar";
import Dashboard from "./subscriptions/subscriptionsDashboard";
import SubscriptionsModal from "./subscriptions/subscriptionsModal";
import SubscriptionsTable from "./subscriptions/subscriptionsTable";
import UserDetails from "./users/userLogin";
import UserLogin from "./users/userLogin";

ReactDOM.render(<App />, document.getElementById("root"));

function App(props) {
  const [subscriptions, setSubscriptions] = useState([]);

  return (
    <>
      <Navbar />

      <UserDetails />

      <Dashboard
        subscriptions={subscriptions}
        setSubscriptions={setSubscriptions}
      />

      <BankModal />
    </>
  );
}
