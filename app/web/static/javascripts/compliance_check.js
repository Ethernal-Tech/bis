var data = {
    Value: "",
	OriginatingBank:[],
	Originator: [],
	BeneficiaryBank: [],
	Beneficiary: [],
	Currency: [],
	AmountFrom: "",
	AmountTo: "",
	StatusId: "",
	From: "",
	To: ""
}
GetComplianceChecks()

var searchField = document.getElementById('compliance-check-search-field');
var timeout = null;

searchField.addEventListener('input', function() {
    if (timeout) {
        clearTimeout(timeout);
    }
    timeout = setTimeout(function() {
        data.Value = searchField.value;
        GetComplianceChecks();
        timeout = null;
    }, 500);
});

function GetComplianceChecks() {
    fetch("/compliancecheck", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.text())
    .then(partialHTML => {
        var partial = document.getElementById('compliance-check-partial')
        partial.innerHTML = ""
        partial.innerHTML = partialHTML
    })
}

function showAdvancedFilter(){
    var divToCheck = document.getElementById('compliance-check-advanced-filter');
    if (divToCheck) {
        var isVisible = window.getComputedStyle(divToCheck).display !== 'none';

        if (isVisible) {
            divToCheck.style.display = 'none';
        } else {
            divToCheck.style.display = 'flex';
        }
    }
}

/*****Calendar*****/
var calendarWindow = document.getElementById('calendar-window');
var calendarBtn = document.getElementById('calendarBtn');
var calendarEl = document.getElementById('calendar');

var calendar = new FullCalendar.Calendar(calendarEl, {
    initialView: 'dayGridMonth', 
    height: '100%', 
    aspectRatio: 1.5,
    contentHeight: 'auto' 
});

calendarBtn.addEventListener('click', function() {
    calendarWindow.style.display = 'flex';
    calendar.render(); 
});

document.addEventListener('click', function(event) {
    if (!calendarWindow.contains(event.target) && !calendarBtn.contains(event.target)) {
        calendarWindow.style.display = 'none';
    }
});