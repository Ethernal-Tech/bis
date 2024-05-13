window.onload = function () {
    document.getElementById("to-checklist-button").addEventListener("click", colapse)
    if (document.getElementById("amount")) {
        amount = document.getElementById("amount").value
        document.getElementById("amount").value = formatAmountValue(amount);
        addSeparators()
    }
}

function colapse() {
    document.getElementById("add-transaction-data").style.display = "none"
    document.getElementById("confirm-buttons").style.display = "flex"
    document.getElementById("to-checklist-buttons").style.display = "none"
    document.getElementById("policies").style.display = "flex"
    document.getElementById("edit-policy-form").style.height = "auto"

    getPolicies()
}

function getPolicies() {
    bankCountry = document.getElementById("bankCountry").value
    policyId = document.getElementById("policyId").value

    data = {
        BankCountry: bankCountry,
        PolicyId: policyId
    }

    fetch("/api/getpolicy", {
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

                var centeredTextDiv = document.createElement("div");
                centeredTextDiv.textContent = "Original policy:";
                centeredTextDiv.style.textAlign = "center";
                centeredTextDiv.style.color = "red";
                centeredTextDiv.style.fontFamily = "Segoe UI";

                newDiv = document.createElement("div")

                newDiv.classList.add("policies-row")
                if (data.Code == "CFM") {
                    newDiv.innerHTML = `
                        <div class="policy-applied" style="color:red">` + data.Name + `</div>
                        <div class="policy-applied" style="color:red">` + parseFloat(data.Parameter.replace(/,/g, '')).toLocaleString() + ` </div>`
                } else {
                    newDiv.innerHTML = `
                        <div class="policy-applied"  style="color:red">` + data.Name + `</div>
                        <div class="policy-applied"  style="color:red">` + data.Parameter + `</div>`
                }

                policiesDiv.appendChild(centeredTextDiv);
                policiesDiv.appendChild(newDiv)

                centeredTextDiv = document.createElement("div");
                centeredTextDiv.textContent = "Updated policy:";
                centeredTextDiv.style.textAlign = "center";
                centeredTextDiv.style.color = "rgb(66, 127, 109)";
                centeredTextDiv.style.fontFamily = "Segoe UI";

                newDiv = document.createElement("div")

                newDiv.classList.add("policies-row")
                if (data.Code == "CFM") {
                    newDiv.innerHTML = `
                        <div class="policy-applied">` + document.getElementById("policyName").value + `</div>
                        <div class="policy-applied">` + document.getElementById("amount").value + ` </div>`
                } else {
                    newDiv.innerHTML = `
                        <div class="policy-applied" style="text-align: center; font-family: Segoe UI">` +
                        "By confirming you are updateing sanctions list to the latest version pubished on " +
                        `<a href="https://www.opensanctions.org/datasets/un_sc_sanctions/">open sactions</a>` +
                        `</div>`
                }

                policiesDiv.appendChild(centeredTextDiv);
                policiesDiv.appendChild(newDiv)
            }
        })
        .catch(error => {
            console.error('Fetch error:', error);
        });
}

function addSeparators() {
    const numberInput = document.getElementById('amount');

    numberInput.addEventListener('input', function (event) {
        editAmountField(event)
    });

    numberInput.addEventListener('beforeinput', function (event) {
        editAmountField(event)
    });

    numberInput.addEventListener('change', function (event) {
        editAmountField(event)
    });

    numberInput.addEventListener('cut', function (event) {
        editAmountField(event)
    });

    numberInput.addEventListener('copy', function (event) {
        editAmountField(event)
    });

    numberInput.addEventListener('paste', function (event) {
        editAmountField(event)
    });
}

function editAmountField(event) {
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
}

function formatAmountValue(amount) {
    return parseFloat(amount.replace(/,/g, '')).toLocaleString();
}

function extractFilenameFromPath(filePath) {
    // Split the file path by backslashes (\) to separate directories and filename
    const parts = filePath.split("\\");
    // Get the last part which should be the filename
    var filename = parts[parts.length - 1];
    filename = filename.replace(/\.csv$/, "");
    return filename;
}