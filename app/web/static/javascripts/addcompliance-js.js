window.onload = function () {
    document.getElementById("cancel-button").addEventListener("click", downgradeView)
    document.getElementById("next-button").addEventListener("click", upgradeView)
    document.getElementById("amount").addEventListener("input", addSeparators)

    setElements()

    addSeparators()
}

let currentView = 1
let buttonSelectable = true

let senderLei
let senderName
let beneficiaryLei
let beneficiaryName
let paymentType
let currency
let amount
let transactionType
let beneficiaryBank

let senderLeiTI
let senderNameTI
let beneficiaryLeiTI
let beneficiaryNameTI
let paymentTypeTI
let currencyTI
let amountTI
let transactionTypeTI
let beneficiaryBankTI

let view1CL = document.getElementById("view-1").classList
let view2CL = document.getElementById("view-2").classList
let view2indCL = document.querySelector(".view-indicator > div:nth-child(2)").classList
let view3CL = document.getElementById("view-3").classList
let view3indCL = document.querySelector(".view-indicator > div:nth-child(3)").classList

let loader

function setElements() {
    senderLei = document.getElementById("sender-lei")
    senderName = document.getElementById("sender-name")
    beneficiaryLei = document.getElementById("beneficiary-lei")
    beneficiaryName = document.getElementById("beneficiary-name")
    paymentType = document.getElementById("payment-type")
    currency = document.getElementById("currency")
    amount = document.getElementById("amount")
    transactionType = document.getElementById("transaction-type")
    beneficiaryBank = document.getElementById("beneficiary-bank")

    senderLeiTI = document.getElementById("ti-sender-lei")
    senderNameTI = document.getElementById("ti-sender-name")
    beneficiaryLeiTI = document.getElementById("ti-beneficiary-lei")
    beneficiaryNameTI = document.getElementById("ti-beneficiary-name")
    paymentTypeTI = document.getElementById("ti-payment-type")
    currencyTI = document.getElementById("ti-currency")
    amountTI = document.getElementById("ti-amount")
    transactionTypeTI = document.getElementById("ti-transaction-type")
    beneficiaryBankTI = document.getElementById("ti-beneficiary-bank")

    loader = document.getElementById("loader")
}

function downgradeView() {
    if (!buttonSelectable) {
        return
    }

    let view1CL = document.getElementById("view-1").classList
    let view2CL = document.getElementById("view-2").classList
    let view2indCL = document.querySelector(".view-indicator > div:nth-child(2)").classList
    let view3CL = document.getElementById("view-3").classList
    let view3indCL = document.querySelector(".view-indicator > div:nth-child(3)").classList

    if (currentView == 1) {
        window.location.href = "/home"
    } else if (currentView == 2) {
        currentView--
        view2CL.remove("display")
        view2CL.add("not-display")
        view2indCL.remove("add-color")
        view2indCL.add("remove-color")
        view1CL.remove("not-display")
        view1CL.add("display")
    } else if (currentView == 3) {
        currentView--
        view3CL.remove("display")
        view3CL.add("not-display")
        view3indCL.remove("add-color")
        view3indCL.add("remove-color")
        view2CL.remove("not-display")
        view2CL.add("display")
    }
}

function upgradeView() {
    if (!buttonSelectable) {
        return
    }

    let view1CL = document.getElementById("view-1").classList
    let view2CL = document.getElementById("view-2").classList
    let view2indCL = document.querySelector(".view-indicator > div:nth-child(2)").classList
    let view3CL = document.getElementById("view-3").classList
    let view3indCL = document.querySelector(".view-indicator > div:nth-child(3)").classList

    if (currentView == 1) {
        currentView++

        if (senderLei.value == "") {
            senderLeiTI.innerText = "-"
        } else {
            senderLeiTI.innerText = senderLei.value
        }
        senderNameTI.innerText = senderName.value
        if (beneficiaryLei.value == "") {
            beneficiaryLeiTI.innerText = "-"
        } else {
            beneficiaryLeiTI.innerText = beneficiaryLei.value
        }
        beneficiaryNameTI.innerText = beneficiaryName.value
        paymentTypeTI.innerText = paymentType.value
        amountTI.innerText = amount.value + " " + currency.value
        transactionTypeTI.innerText = transactionType.options[transactionType.selectedIndex].text
        beneficiaryBankTI.innerText = beneficiaryBank.value

        view1CL.remove("display")
        view1CL.add("not-display")
        view2indCL.remove("remove-color")
        view2indCL.add("add-color")
        view2CL.remove("not-display")
        view2CL.add("display")
    } else if (currentView == 2) {
        currentView++
        view2CL.remove("display")
        view2CL.add("not-display")
        view3indCL.remove("remove-color")
        view3indCL.add("add-color")

        getPolicies()
    } else if (currentView == 3) {
        // submit compliance check
        var insertedData = {
            "senderLei": senderLei.value,
            "senderName": senderName.value,
            "beneficiaryLei": beneficiaryLei.value,
            "beneficiaryName": beneficiaryName.value,
            "paymentType": paymentType.value,
            "transactionType": transactionType.value,
            "currency": currency.options[currency.selectedIndex].text,
            "amount": amount.value,
            // TODO: Change to dropdown select
            "beneficiaryBank": beneficiaryBank.value
        };

        fetch("/addtransaction", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(insertedData)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                window.location.replace("/home");
            })
    }
}

function addSeparators(event) {
    let inputValue = parseFloat(event.target.value.replace(/,/g, ''));

    if (!isNaN(inputValue)) {
        event.target.value = inputValue.toLocaleString();
    } else {
        event.target.value = '';
    }
}

function getPolicies() {
    data = {
        BeneficiaryBankGlobalIdentifier: beneficiaryBank.value,
        TransactionTypeId: transactionType.value
    }

    showLoader()
    fetch("api/getbeneficiarybankpolicies", {
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
                console.log('no policies to be applied')
            }
            else {
                hideLoader()
                console.log(data)
                let policiesDiv = document.getElementById('policies')
                policiesDiv.innerHTML = ''
                data.policies.forEach(policy => {
                    const policyDiv = document.createElement('div');
                    policyDiv.classList.add('policy')

                    let param

                    if (policy.code == "CFM") {
                        let inputValue = parseFloat(policy.params.replace(/,/g, ''));

                        if (!isNaN(inputValue)) {
                            param = inputValue.toLocaleString();
                        } else {
                            param = '';
                        }
                    } else {
                        param = policy.params
                    }

                    policyDiv.innerHTML = `
                    <svg class="policy-icon" xmlns="http://www.w3.org/2000/svg" width="32px" height="32px" fill="rgb(247, 233, 156)" class="bi bi-file-earmark-text-fill" viewBox="0 0 16 16">
                        <path d="M9.293 0H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V4.707A1 1 0 0 0 13.707 4L10 .293A1 1 0 0 0 9.293 0M9.5 3.5v-2l3 3h-2a1 1 0 0 1-1-1M4.5 9a.5.5 0 0 1 0-1h7a.5.5 0 0 1 0 1zM4 10.5a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5m.5 2.5a.5.5 0 0 1 0-1h4a.5.5 0 0 1 0 1z"/>
                    </svg>
                    <div class="policy-content">
                        <div class="policy-code-name">
                            ${policy.code} - ${policy.name}
                        </div>
                        <div class="policy-params">
                            ${param}
                        </div>
                    </div>
                    `;

                    policiesDiv.appendChild(policyDiv);
                });
                view3CL.remove("not-display")
                view3CL.add("display")
            }
        })
        .catch(error => {
            console.error('Fetch error:', error);
        });
}

function showLoader() {
    buttonSelectable = false
    loader.style.display = 'block';
}

function hideLoader() {
    buttonSelectable = true
    loader.style.display = 'none';
}