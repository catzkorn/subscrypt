loadTransactions();

function loadTransactions() {
    _getTransactions(_showTransactions);
}


function _getTransactions(callback) {
    let xhttp = new XMLHttpRequest();
    let path = '/api/listoftransactions';
    xhttp.onreadystatechange = function() {
        if (xhttp.readyState === 4 && xhttp.status === 200) {
            let transactions = _convertToTransactions(xhttp.responseText);
            callback(transactions);
        }
    };
    xhttp.open("GET", path, true);
    xhttp.send();
}

function _convertToTransactions(res) {
    let resTransactions = JSON.parse(res);
    let transactions = [];
    if (resTransactions === null) {
        return transactions;
    } else {
        resTransactions.transactions.forEach(function(transaction) {
            let transactionObj = new Transaction(transaction.name, transaction.date, transaction.amount);
            transactions.push(transactionObj);
        });
        return transactions;
    }
}

function _showTransactions(transactions) {
    let transactionsHTML = `<h4>Subscriptions</h4>`;
    if (transactions.length > 0) {
        transactionsHTML += _formatTransactionsTable(transactions);
    } else {
        transactionsHTML += "<p>You don't have any Transactions</p>";
    }
    document.getElementById("transactions").innerHTML = transactionsHTML;
}

function _formatTransactionsTable(transactions) {
    let tableHTML = `<table class="table" id="table-transactions">
                        <thead>
                            <tr>
                                <th scope="col">Transaction</th>
                                <th scope="col">Date</th>
                                <th scope="col">Amount</th>
                            </tr>
                        </thead>
                        <tbody>`;

    transactions.forEach(function(transaction) {
        tableHTML += _formatTransaction(transaction);
    });

    tableHTML += "</tbody></table>";
    return tableHTML;
}

function _formatTransaction(transaction) {
    return `<tr>
            <td>${transaction.name}</td>
            <td>${transaction.date}</td>
            <td>${_formatAmountTwoDecimals(transaction.amount)}</td>
            </tr>`;
}

function _formatAmountTwoDecimals(amount) {
    return parseFloat(amount).toFixed(2);
}