import React, { useEffect, useState } from "react";

function SubscriptionsTable(props) {
  const [subscriptions, setSubscriptions] = useState([]);

  useEffect(() => {
    const url = "/api/subscriptions";
    fetch(url)
      .then((response) => {
        return response.json();
      })
      .then((payload) => {
        setSubscriptions(payload);
      });
  }, []);

  function renderSubscription(subscription, index) {
    return (
      <tr key={subscription.id}>
        <th scope="row">{subscription.name}</th>
        <td>{_formatAmountTwoDecimals(subscription.amount)}</td>
        <td>{_formatDateAsDay(subscription.dateDue)}</td>
        <td>Monthly</td>
        <td>
          <button type="button" className="icon-button" id="reminder-button">
            {/* {calendarSvg} */}
          </button>
          <button
            type="button"
            className="icon-button"
            id="delete-${subscription.id}"
          >
            {/* ${binSvg} */}
          </button>
        </td>
      </tr>
    );
  }

  function _formatAmountTwoDecimals(amount) {
    return parseFloat(amount).toFixed(2);
  }

  function _formatDateAsDay(dateString) {
    const date = new Date(dateString);
    let d = date.getDate();
    return d + _getOrdinal(d);
  }

  function _getOrdinal(number) {
    if (number > 3 && number < 21) return "th";
    switch (number % 10) {
      case 1:
        return "st";
      case 2:
        return "nd";
      case 3:
        return "rd";
      default:
        return "th";
    }
  }

  return (
    <div className="container" id="subscriptions">
      <div id="subscriptions-table">
        <table className="table" id="table-subscriptions">
          <thead>
            <tr>
              <th scope="col">Subscription Name</th>
              <th scope="col">Amount</th>
              <th scope="col">Payment Date</th>
              <th scope="col">Frequency</th>
              <th scope="col">Actions</th>
            </tr>
          </thead>
          <tbody>{subscriptions.map(renderSubscription)}</tbody>
        </table>
      </div>
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
