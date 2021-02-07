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

  function renderTransaction(transaction, index) {
    return (
      <tr key={transaction.name + transaction.date + transaction.amount}>
        <td>{transaction.name}</td>
        <td>{transaction.date}</td>
        <td>£{formatAmountTwoDecimals(transaction.amount)}</td>
      </tr>
    );
  }

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
      <tbody>{transactions.map(renderTransaction)}</tbody>
    </table>
  );
}

export default Transactions;
