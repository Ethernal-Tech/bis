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
        populateFilters()
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
        populateFilters()
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

/******* FILTER START  *********/
function openPopup(triggerElement, popupId) {
    var popup = document.getElementById(popupId);
    
    // Check if the popup is currently visible
    var isVisible = popup.style.display === 'block';

    // Hide all popups if not clicking on an already visible one
    document.querySelectorAll('.popup-container').forEach(function(p) {
        p.style.display = 'none';
    });

    // If the popup is not visible, show it
    if (!isVisible) {
        var rect = triggerElement.getBoundingClientRect();
        popup.style.top = (rect.bottom + window.scrollY) + 'px';
        popup.style.left = (rect.left + window.scrollX) + 'px';
        popup.style.display = 'block';
    }
}

// Hide the popup when clicking outside of it
document.addEventListener('click', function(event) {
    var popups = document.querySelectorAll('.popup-container');
    var triggers = document.querySelectorAll('.compliance-check-advanced-filter-item');

    // Check if click is outside any popup or trigger
    var clickOutside = true;
    popups.forEach(function(popup) {
        if (popup.contains(event.target)) {
            clickOutside = false;
        }
    });
    triggers.forEach(function(trigger) {
        if (trigger.contains(event.target)) {
            clickOutside = false;
        }
    });

    if (clickOutside) {
        popups.forEach(function(popup) {
            popup.style.display = 'none';
        });
    }
});

var table;
var filterColumns;
var filters;

function populateFilters() {
    table = document.getElementById('compliance-check-table');
    filterColumns = [0, 1, 2, 3, 4];
    filters = [
        document.getElementById('originating-bank-list'),
        document.getElementById('originator-list'),
        document.getElementById('beneficiary-bank-list'),
        document.getElementById('beneficiary-list'),
        document.getElementById('currency-list')
    ];
    const uniqueValues = filterColumns.map(() => []);

    Array.from(table.tBodies[0].rows).forEach(row => {
        filterColumns.forEach((colIndex, filterIndex) => {
            const cell = row.cells[colIndex];
            if (!uniqueValues[filterIndex].includes(cell.textContent)) {
                uniqueValues[filterIndex].push(cell.textContent);
            }
        });
    });

    // Clear existing filter options
    filters.forEach(filter => {
        filter.innerHTML = '';
    });

    uniqueValues.forEach((values, index) => {
        values.sort();
        values.forEach(value => {
            filters[index].insertAdjacentHTML('beforeend', 
                `<div class="popup-filter-item">
                    <input type="checkbox" id="filter-${index}-${value}" value="${value}">
                    <label for="filter-${index}-${value}">${value}</label>
                </div>`
            );
        });
    });

    // Attach event listeners to all checkboxes
    filters.forEach((filter, index) => {
        filter.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
            checkbox.addEventListener('change', filterTable);
        });
    });
}

function filterTable() {
    const selectedValues = filters.map(filter => {
        return Array.from(filter.querySelectorAll('input[type="checkbox"]:checked')).map(checkbox => checkbox.value);
    });

    Array.from(table.tBodies[0].rows).forEach(row => {
        const shouldShow = filterColumns.every((colIndex, filterIndex) => {
            const cell = row.cells[colIndex];
            return selectedValues[filterIndex].length === 0 || selectedValues[filterIndex].includes(cell.textContent);
        });
        row.style.display = shouldShow ? "" : "none";
    });
}

/******* FILTER END  *********/


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