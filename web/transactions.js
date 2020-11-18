loadTransactions();

function loadTransactions() {
    _getTransactions(_showTransactions);
}


function _getTransactions(callback) {
    let xhttp = new XMLHttpRequest();
    let path = '/api/listoftransactions';
    xhttp.onreadystatechange = function() {
        if (xhttp.readyState === 4 && xhttp.status === 200) {
            console.log(xhttp.responseText)
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
    let transactionsHTML = "";
    if (transactions.length > 0) {
        transactionsHTML = _formatTransactionsTable(transactions);
    } else {
        transactionsHTML = "You don't have any subscriptions";
    }
    document.getElementById("transactions").innerHTML = transactionsHTML;
}

function _formatTransactionsTable(transactions) {
    let tableHTML = `<table id=\"table-subscriptions\" style=\"width:100%\">
                                <tr>
                                    <td>Transaction</td>
                                    <td>Date</td>
                                    <td>Amount</td>
                                </tr>`;

    transactions.forEach(function(transaction) {
        tableHTML += _formatTransaction(transaction);
    });

    tableHTML += "</table>";
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