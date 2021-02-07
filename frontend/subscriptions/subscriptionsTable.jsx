import React, { useEffect, useState } from "react";
import formatAmountTwoDecimals from "../util/formatNumbers";

const calendarSvg = (
  <svg
    className="w-6 h-6"
    fill="currentColor"
    viewBox="0 0 20 20"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fillRule="evenodd"
      d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
      clipRule="evenodd"
    ></path>
  </svg>
);
const binSvg = (
  <svg
    className="w-6 h-6"
    fill="currentColor"
    viewBox="0 0 20 20"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fillRule="evenodd"
      d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
      clipRule="evenodd"
    ></path>
  </svg>
);

function SubscriptionsTable(props) {
  useEffect(() => {
    const url = "/api/subscriptions";
    fetch(url)
      .then((response) => {
        return response.json();
      })
      .then((payload) => {
        payload.sort((a, b) => {
          const dayA = new Date(a.dateDue).getDate();
          const dayB = new Date(b.dateDue).getDate();
          if (dayA < dayB) {
            return -1;
          } else if (dayA > dayB) {
            return 1;
          }
          return 0;
        });

        props.setSubscriptions(payload);
      });
  }, []);

  if (props.subscriptions.length === 0) {
    return <p>You don't have any subscriptions</p>;
  }
  return (
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
        <tbody>
          {props.subscriptions.map((subscription) => {
            return (
              <Subscription subscription={subscription} key={subscription.id} />
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

function Subscription(props) {
  function handleDeleteSubscription(subscriptionID) {
    const url = "/api/subscriptions/" + subscriptionID;
    const options = {
      method: "DELETE",
    };
    fetch(url, options).then((response) => {
      if (response.status !== 200) {
        console.log("There was an error deleting the subscription", response);
        return;
      }
      newSubscriptions = props.subscriptions.filter((subscription) => {
        return subscription.id !== subscriptionID;
      });
      props.setSubscriptions(newSubscriptions);
    });
  }

  function handleSendReminder(subscriptionID) {
    const url = "/api/reminders";
    const options = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: subscriptionID,
      }),
    };
    fetch(url, options).then((response) => {
      if (response.status !== 200) {
        console.log("There was an error sending a reminder", response);
        return;
      }
      console.log("it worked!");
    });
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

  function _formatDateAsDay(dateString) {
    const date = new Date(dateString);
    let d = date.getDate();
    return d + _getOrdinal(d);
  }

  return (
    <tr>
      <th scope="row">{props.subscription.name}</th>
      <td>{formatAmountTwoDecimals(props.subscription.amount)}</td>
      <td>{_formatDateAsDay(props.subscription.dateDue)}</td>
      <td>Monthly</td>
      <td>
        <button
          type="button"
          className="icon-button"
          onClick={() => handleSendReminder(props.subscription.id)}
          id="reminder-button"
        >
          {calendarSvg}
        </button>
        <button
          type="button"
          className="icon-button"
          onClick={() => handleDeleteSubscription(props.subscription.id)}
        >
          {binSvg}
        </button>
      </td>
    </tr>
  );
}

export default SubscriptionsTable;
