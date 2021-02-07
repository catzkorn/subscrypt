import React, { useState } from "react";

function SubscriptionsModal(props) {
  const [subscriptionDate, setSubscriptionDate] = useState("");
  const [subscriptionName, setSubscriptionName] = useState("");
  const [subscriptionAmount, setSubscriptionAmount] = useState(0);
  const [validationErrorMessage, setValidationErrorMessage] = useState("");

  function handleSubscriptionSubmit(event) {
    event.preventDefault();

    if (
      !_subscriptionValuesAreValid(
        subscriptionName,
        subscriptionAmount,
        subscriptionDate
      )
    ) {
      return;
    }

    const url = "/api/subscriptions";
    const formatDate = _formatDateForJSON(subscriptionDate);
    const options = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: subscriptionName,
        amount: subscriptionAmount,
        dateDue: formatDate,
      }),
    };

    fetch(url, options)
      .then((response) => {
        if (response.status !== 200) {
          console.log("There was an error with the submitted data", response);
          return;
        }
        return response.json();
      })
      .then((payload) => {
        subscriptions = props.subscriptions.concat([payload]);
        subscriptions.sort(sortDates);
        props.setSubscriptions(subscriptions);
        props.setShowModal(false);
      });
  }

  function _subscriptionValuesAreValid(
    subscriptionName,
    subscriptionAmount,
    subscriptionDate
  ) {
    if (
      subscriptionName === "" ||
      subscriptionAmount === "" ||
      subscriptionDate === ""
    ) {
      setValidationErrorMessage("One or more fields is empty");
      return false;
    }

    if (parseInt(subscriptionAmount) <= 0) {
      setValidationErrorMessage("Amount can not be 0 or a negative number");
      return false;
    }
    return true;
  }

  function _formatDateForJSON(date) {
    return date + "T00:00:00Z";
  }

  function handleSubscriptionDateChange(event) {
    setSubscriptionDate(event.target.value);
  }

  function handleSubscriptionNameChange(event) {
    setSubscriptionName(event.target.value);
  }

  function handleSubscriptionAmountChange(event) {
    setSubscriptionAmount(event.target.value);
  }
  if (!props.showModal) {
    return null;
  }
  return (
    <div
      className="modal fade"
      id="addSubscriptionModal"
      tabIndex="-1"
      role="dialog"
      aria-labelledby="addSubscriptionModalLabel"
      aria-hidden="true"
    >
      <div className="modal-dialog modal-dialog-centered" role="document">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title" id="addSubscriptionModalLabel">
              Add a subscription
            </h5>
            <button
              type="button"
              className="close"
              data-dismiss="modal"
              aria-label="Close"
              onClick={() => props.setShowModal(false)}
            >
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div className="modal-body">
            <form>
              <div className="form-group">
                <label htmlFor="subscription-name" className="col-form-label">
                  Subscription name:
                </label>
                <input
                  type="text"
                  className="form-control"
                  onChange={(event) => handleSubscriptionNameChange(event)}
                  value={subscriptionName}
                  id="subscription-name"
                ></input>
              </div>
              <div className="form-group">
                <label htmlFor="subscription-amount" className="col-form-label">
                  Price:
                </label>
                <input
                  type="text"
                  className="form-control"
                  onChange={(event) => handleSubscriptionAmountChange(event)}
                  value={subscriptionAmount}
                  id="subscription-amount"
                ></input>
              </div>
              <div className="form-group">
                <label htmlFor="subscription-date" className="col-form-label">
                  Next payment date:
                </label>
                <input
                  type="date"
                  className="form-control"
                  onChange={(event) => handleSubscriptionDateChange(event)}
                  value={subscriptionDate}
                  id="subscription-date"
                ></input>
              </div>
            </form>
          </div>
          <div className="subscription-validation-error">
            {validationErrorMessage}
          </div>
          <div className="modal-footer">
            <button
              type="button"
              className="btn btn-secondary"
              data-dismiss="modal"
              onClick={() => props.setShowModal(false)}
            >
              Close
            </button>
            <button
              type="button"
              className="btn btn-primary"
              id="create-subscription-button"
              data-dismiss="modal"
              onClick={(event) => {
                setValidationErrorMessage("");
                handleSubscriptionSubmit(event);
              }}
            >
              Add subscription
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default SubscriptionsModal;
