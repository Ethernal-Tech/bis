window.onload = function() {
    document.getElementById("to-checklist-button").addEventListener("click", colapse)
}

function colapse() {
    document.getElementById("add-transaction-data").style.display = "none"
    document.getElementById("confirm-buttons").style.display = "flex"
    document.getElementById("to-checklist-buttons").style.display = "none"
    document.getElementById("policies").style.display = "flex"

    getPolicies()
}
 
function getPolicies() {
    bankId = document.getElementById("select-bank").value
    transactionTypeid = document.getElementById("select-type").value

    data = {
        BankId : bankId,
        TransactionTypeId : transactionTypeid
    }

    fetch("/api/getpolicies", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json', 
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        policiesDiv = document.getElementById("policies")
        data.forEach(function(currentValue){
            newDiv = document.createElement("div")
            newDiv.classList.add("policies-row")
            newDiv.innerHTML = `
                <div class="policy-applied">` + currentValue.Name + `</div>
                <div class="policy-applied">` + currentValue.Parameter + `</div>`
            policiesDiv.appendChild(newDiv)
        })
    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}