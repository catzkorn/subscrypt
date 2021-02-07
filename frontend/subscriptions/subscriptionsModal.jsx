import React from "react";

function SubscriptionsModal(props) {
  return (
    <div
      className="modal fade"
      id="addSubscriptionModal"
      tabindex="-1"
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
            >
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div className="modal-body">
            <form>
              <div className="form-group">
                <label for="subscription-name" className="col-form-label">
                  Subscription name:
                </label>
                <input
                  type="text"
                  className="form-control"
                  id="subscription-name"
                ></input>
              </div>
              <div className="form-group">
                <label for="subscription-amount" className="col-form-label">
                  Price:
                </label>
                <input
                  type="text"
                  className="form-control"
                  id="subscription-amount"
                ></input>
              </div>
              <div className="form-group">
                <label for="subscription-date" className="col-form-label">
                  Next payment date:
                </label>
                <input
                  type="date"
                  className="form-control"
                  id="subscription-date"
                ></input>
              </div>
            </form>
          </div>
          <div className="modal-footer">
            <button
              type="button"
              className="btn btn-secondary"
              data-dismiss="modal"
            >
              Close
            </button>
            <button
              type="button"
              className="btn btn-primary"
              id="create-subscription-button"
              data-dismiss="modal"
              onclick="createSubscription()"
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
