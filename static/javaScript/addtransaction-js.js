window.onload = function () {
    document.getElementById("to-checklist-button").addEventListener("click", colapse)
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

            if (data.length == 0) {
                document.getElementById("policy-not-applied").style.display = "flex"
            }
            else {
                policiesDiv = document.getElementById("policies")
                data.forEach(function (currentValue) {
                    newDiv = document.createElement("div")
                    newDiv.classList.add("policies-row")
                    newDiv.innerHTML = `
                    <div class="policy-applied">` + currentValue.Name + `</div>
                    <div class="policy-applied">` + currentValue.Parameter + `</div>`
                    policiesDiv.appendChild(newDiv)
                })
            }
        })
        .catch(error => {
            console.error('Fetch error:', error);
        });
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