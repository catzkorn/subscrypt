import React, { useState } from "react";
import BankModal from "../banking/bankModal";
import SubscriptionsModal from "./subscriptionsModal";
import SubscriptionsTable from "./subscriptionsTable";

function Dashboard(props) {
  const [showSubscriptionModal, setShowSubscriptionModal] = useState(false);
  const [showBankModal, setShowBankModal] = useState(false);
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
        onClick={() => setShowSubscriptionModal(true)}
      >
        Add new subscription
      </button>
      <button
        type="button"
        className="btn btn-primary"
        data-toggle="modal"
        data-target="#chooseBankAccountModal"
        onClick={() => setShowBankModal(true)}
      >
        Load from bank account
      </button>

      <SubscriptionsModal
        subscriptions={props.subscriptions}
        setSubscriptions={props.setSubscriptions}
        showSubscriptionModal={showSubscriptionModal}
        setShowSubscriptionModal={setShowSubscriptionModal}
      />

      <BankModal
        showBankModal={showBankModal}
        setShowBankModal={setShowBankModal}
      />
    </div>
  );
}

export default Dashboard;
