import React, { useState } from "react";
import SubscriptionsModal from "./subscriptionsModal";
import SubscriptionsTable from "./subscriptionsTable";

function Dashboard(props) {
  const [showModal, setShowModal] = useState(false);
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
        onClick={() => setShowModal(true)}
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
        showModal={showModal}
        setShowModal={setShowModal}
      />
    </div>
  );
}

export default Dashboard;
