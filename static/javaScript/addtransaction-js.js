window.onload = function() {
    document.getElementById("select-bank").addEventListener("change", getPolicies)
    document.getElementById("select-type").addEventListener("change", getPolicies)

    getPolicies()
}

function getPolicies() {
    bankId = document.getElementById("select-bank").value
    transactionTypeid = document.getElementById("select-type").value

    data = {
        BankId : bankId,
        TransactionTypeId : transactionTypeid
    }

    fetch("/getpolicies", {
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
        hasSCL = false
        hasCFM = false
        data.forEach(function(currentValue){
            if (currentValue.Name == "Saction Check List") {
                document.getElementById("country").innerText = "(" + currentValue.Country + ")"
                document.getElementById("SCL-sign").style.backgroundColor = "#706f6f"
                hasSCL = true
            }

            if (currentValue.Name == "Capital Flow Management") {
                document.getElementById("amount").innerText = "(" + currentValue.Amount + ")"
                document.getElementById("CFM-sign").style.backgroundColor = "#706f6f"
                hasCFM = true
            }
        })

        if (hasSCL == false) {
            document.getElementById("country").innerText = "(-)"
            document.getElementById("SCL-sign").style.backgroundColor = "#ffffff"            
        }

        if (hasCFM == false) {
            document.getElementById("amount").innerText = "(-)"
            document.getElementById("CFM-sign").style.backgroundColor = "#ffffff"
        }

    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}