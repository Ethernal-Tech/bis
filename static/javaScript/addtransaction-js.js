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
        hasSCL = false
        hasCFM = false
        data.forEach(function(currentValue){
            if (currentValue.Name == "Saction Check List") {
                document.getElementById("country").innerText = "(" + currentValue.Country + ")"
                document.getElementById("country").style.color = "rgb(66, 127, 109)"
                document.getElementById("SCL").style.color = "rgb(66, 127, 109)"
                document.getElementById("SCL").style.textDecorationLine = "none"    
                hasSCL = true
            }

            if (currentValue.Name == "Capital Flow Management") {
                document.getElementById("amount").innerText = "(" + currentValue.Amount + ")"
                document.getElementById("amount").style.color = "rgb(66, 127, 109)"
                document.getElementById("CFM").style.color = "rgb(66, 127, 109)"
                document.getElementById("CFM").style.textDecorationLine = "none"       
                hasCFM = true
            }
        })

        if (hasSCL == false) {
            document.getElementById("country").innerText = "(-)"
            document.getElementById("country").style.color = "rgb(183, 183, 183)"
            document.getElementById("SCL").style.color = "rgb(183, 183, 183)"  
            document.getElementById("SCL").style.textDecorationLine = "line-through"       
        }

        if (hasCFM == false) {
            document.getElementById("amount").innerText = "(-)"
            document.getElementById("amount").style.color = "rgb(183, 183, 183)"
            document.getElementById("CFM").style.color = "rgb(183, 183, 183)"  
            document.getElementById("CFM").style.textDecorationLine = "line-through"       
        }

    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}