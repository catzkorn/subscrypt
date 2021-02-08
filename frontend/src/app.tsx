import React, { useState } from "react";
import ReactDOM from "react-dom";
import { Route, BrowserRouter as Router, Switch } from "react-router-dom";
import BankModal from "./banking/bankModal";
import Navbar from "./navbar/navbar";
import Dashboard from "./subscriptions/subscriptionsDashboard";
import SubscriptionsModal from "./subscriptions/subscriptionsModal";
import SubscriptionsTable from "./subscriptions/subscriptionsTable";
import Transactions from "./transactions/transactions";
import UserDetails from "./users/userLogin";
import UserLogin from "./users/userLogin";
import Subscription from "./subscriptions/subscriptionType";
import User from "./users/userTypes";

ReactDOM.render(<App />, document.getElementById("root"));

function App() {
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);
  const [user, setUser] = useState<User | null>(null);

  if (user === null) {
    return (
      <Router>
        <Navbar />

        <UserDetails user={user} setUser={setUser} />
      </Router>
    );
  }

  return (
    <Router>
      <Navbar />
      <Switch>
        <Route path="/transactions">
          <Transactions />
        </Route>

        <Route path="/">
          <UserDetails user={user} setUser={setUser} />

          <Dashboard
            subscriptions={subscriptions}
            setSubscriptions={setSubscriptions}
          />
        </Route>
      </Switch>
    </Router>
  );
}
