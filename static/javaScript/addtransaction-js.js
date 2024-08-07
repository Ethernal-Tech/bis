window.onload = function () {
    document.getElementById("to-checklist-button").addEventListener("click", colapse)
    document.getElementById("search-senders").addEventListener("click", openSendersWindow)
    document.getElementById("senders-window-back-button").addEventListener("click", closeSendersWindow)
    addSeparators()
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
        BankId: bankId,
        TransactionTypeId: transactionTypeid
    }

    fetch("/api/getbeneficiarybankpolicies", {
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
        if (data.length == 0) {
            document.getElementById("policy-not-applied").style.display = "flex"
        }
        else {
            policiesDiv = document.getElementById("policies")
            data.forEach(function (currentValue) {
            newDiv = document.createElement("div")
            newDiv.classList.add("policies-row")
            
            if (currentValue.Code == "CFM") {
                newDiv.innerHTML = `
                    <div class="policy-applied">` + currentValue.Name + `</div>
                    <div class="policy-applied">` + parseFloat(currentValue.Parameter.replace(/,/g, '')).toLocaleString() + ` (Loan ID = ` + document.getElementById("loanId").value + `) </div>`                        
            } 
            else {
                newDiv.innerHTML = `
                    <div class="policy-applied">` + currentValue.Name + `</div>
                        <div class="policy-applied">` + currentValue.Parameter + `</div>`
            }
            policiesDiv.appendChild(newDiv)
            })
        }
    })
    .catch(error => {
            console.error('Fetch error:', error);
    });
}

function openSendersWindow() {
    // the API is too slow to instantly set the display property of elements
    // document.getElementById("senders-window").style.display="block"
    // document.getElementById("add-transaction").style.display="none"
    
    var value = document.getElementById("input-sender").value

    var body = document.getElementById("senders-window-body")
    body.innerHTML = ""

    const apiUrl = 'https://api.gleif.org/api/v1/autocompletions?field=fulltext&q=' + value;
    var senders = [];
    fetch(apiUrl)
    .then(response => {
        if (!response.ok) {
        throw new Error('Network response was not ok');
        }

        return response.json();
    })
    .then(data => {
        data.data.forEach(el => {
            var newRow = document.createElement("div")
            newRow.textContent = el.attributes.value
            newRow.classList.add("senders-window-row")
            newRow.addEventListener('click', function () {
                document.getElementById("input-sender").value = el.attributes.value; 
                
                document.getElementById("senders-window").style.display="none"
                document.getElementById("add-transaction").style.display="block"
              });
            body.appendChild(newRow)
        })

        document.getElementById("senders-window").style.display="block"
        document.getElementById("add-transaction").style.display="none"
    })
    .catch(error => {
        // Handle any errors that occurred during the fetch
        console.error('There was a problem with the fetch operation:', error);
    });
}

function closeSendersWindow() {
    document.getElementById("senders-window").style.display="none"
    document.getElementById("add-transaction").style.display="block"
}

function addSeparators() {
    const numberInput = document.getElementById('numberInput');

    numberInput.addEventListener('input', function (event) {
        // Remove existing commas and parse the input value as a number
        const inputValue = parseFloat(event.target.value.replace(/,/g, ''));

        // Check if the input value is a valid number
        if (!isNaN(inputValue)) {
            // Format the number with thousand separators and update the input value
            event.target.value = inputValue.toLocaleString();
        } else {
            // If the input is not a valid number, clear the input field
            event.target.value = '';
        }
    });
}