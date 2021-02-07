import React from "react";

function SubscriptionsTable(props) {
  return (
    <div className="container" id="subscriptions">
      <div id="subscriptions-table"></div>
      <span id="reminder-error"></span>
      <div className="d-flex justify-content-center text-secondary">
        <div className="spinner-border" role="status" id="loading-spinner">
          <span className="sr-only">Loading...</span>
        </div>
      </div>
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
    </div>
  );
}

export default SubscriptionsTable;
