import React from "react";
import SubscriptionsModal from "./subscriptionsModal";
import SubscriptionsTable from "./subscriptionsTable";

function Dashboard(props) {
  return (
    <div className="container" id="subscriptions">
      <SubscriptionsTable
        subscriptions={props.subscriptions}
        setSubscriptions={props.setSubscriptions}
      />

      <button
        type="button"
        className="btn btn-primary"
        data-toggle="modal"
        data-target="#addSubscriptionModal"
      >
        Add new subscription
      </button>
      <button
        type="button"
        className="btn btn-primary"
        data-toggle="modal"
        data-target="#chooseBankAccountModal"
      >
        Load from bank account
      </button>

      <SubscriptionsModal
        subscriptions={props.subscriptions}
        setSubscriptions={props.setSubscriptions}
      />
    </div>
  );
}

export default Dashboard;
