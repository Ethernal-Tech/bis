window.onload = function () {
    document.getElementById("cancel-button").addEventListener("click", downgradeView)
    document.getElementById("next-button").addEventListener("click", upgradeView)
    document.getElementById("amount").addEventListener("input", addSeparators)

    setElements()

    addSeparators()
}

let currentView = 1

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
}

function downgradeView() {
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
        amountTI.innerText= amount.value + " " + currency.value
        transactionTypeTI.innerText = transactionType.value
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
        view3CL.remove("not-display")
        view3CL.add("display")

        getPolicies()
    } else if (currentView == 3) {
        // submit compliance check
        var insertedData = {
            "senderLei": document.getElementById("sender-lei").value,
            "senderName": document.getElementById("sender-name").value,
            "beneficiaryLei": document.getElementById("beneficiary-lei").value,
            "beneficiaryName": document.getElementById("beneficiary-name").value,
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
    //bankId = document.getElementById("select-bank").value
    //transactionTypeid = document.getElementById("select-type").value

    data = {
        BankId: "549300BUPYUQGB5BFX94",
        TransactionTypeId: "1"
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
                console.log('no policies to be applied')
            }
            else {
                console.log(data)
            }
        })
        .catch(error => {
            console.error('Fetch error:', error);
        });
}