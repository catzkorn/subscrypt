import React, { useEffect, useState } from "react";
import formatAmountTwoDecimals from "../util/formatNumbers";

function Transactions() {
  const [transactions, setTransactions] = useState([]);

  useEffect(() => {
    const url = "/api/transactions";
    fetch(url)
      .then((response) => {
        return response.json();
      })
      .then((payload) => {
        console.log(payload);
        setTransactions(payload);
      });
  }, []);

  if (transactions.length === 0) {
    return <p>You don't have any transactions</p>;
  }

  return (
    <table className="table" id="table-transactions">
      <thead>
        <tr>
          <th scope="col">Transaction</th>
          <th scope="col">Date</th>
          <th scope="col">Amount</th>
        </tr>
      </thead>
      <tbody>
        {transactions.map((transaction) => {
          return (
            <Transaction
              transaction={transaction}
              key={transaction.name + transaction.date + transaction.amount}
            />
          );
        })}
      </tbody>
    </table>
  );
}

function Transaction(props) {
  return (
    <tr>
      <td>{props.transaction.name}</td>
      <td>{props.transaction.date}</td>
      <td>£{formatAmountTwoDecimals(props.transaction.amount)}</td>
    </tr>
  );
}

export default Transactions;
