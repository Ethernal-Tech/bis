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
        populateStatuses()
        populateFilters()
    })
}

function populateStatuses() {
    var statuses = document.querySelectorAll('.compliance-check-status');
    statuses.forEach((status) => {
        switch (status.textContent.trim()) {
            case 'Failed':
                status.classList.add('state-error');
              break;
            case 'Successful':
                status.classList.add('state-success');
              break;
            case 'Pending':
                status.classList.add('state-pending');
              break;
            default:
              break;
          }
      });
}

/******* FILTER START  *********/
function openPopup(triggerElement, popupId) {
    var popup = document.getElementById(popupId);
    
    // Check if the popup is currently visible
    var isVisible = popup.style.display === 'block';

    // Hide all popups if not clicking on an already visible one
    document.querySelectorAll('.popup-container, .popup-container-amount').forEach(function(p) {
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
    var popups = document.querySelectorAll('.popup-container, .popup-container-amount');
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
var searchField;
var fromAmountField;
var toAmountField;
var startDate = null;
var endDate = null;  

function populateFilters() {
    table = document.getElementById('compliance-check-table');
    filterColumns = [0, 1, 2, 3, 4, 6];
    filters = [
        document.getElementById('originating-bank-list'),
        document.getElementById('originator-list'),
        document.getElementById('beneficiary-bank-list'),
        document.getElementById('beneficiary-list'),
        document.getElementById('currency-list'),
        document.getElementById('status-list'),
    ];
    
    searchField = document.getElementById('compliance-check-search-field');
    fromAmountField = document.getElementById('from-amount');
    toAmountField = document.getElementById('to-amount');

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
                    <label class="body-semibold text-strong" for="filter-${index}-${value}">${value}</label>
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

    // Attach event listener to search field
    searchField.addEventListener('input', filterTable);
    var searchElements = document.querySelectorAll('.compliance-check-search');
    searchElements.forEach(search => {
        var icons = search.querySelectorAll('svg path');
        var input = search.querySelector('input');

        input.addEventListener('focus', () => {
            icons.forEach(svg => svg.classList.remove('icon-stroke-placeholder'));
            icons.forEach(svg => svg.classList.add('icon-stroke-strong'));
        })
        input.addEventListener('blur', () => {
            icons.forEach(svg => svg.classList.remove('icon-stroke-strong'));
            icons.forEach(svg => svg.classList.add('icon-stroke-placeholder'));
        })
    })

    fromAmountField.addEventListener('input', filterTable);
    toAmountField.addEventListener('input', filterTable);

    addSearchOnFilters()
}

function filterTable() {
    const selectedValues = filters.map(filter => {
        return Array.from(filter.querySelectorAll('input[type="checkbox"]:checked')).map(checkbox => checkbox.value);
    });
    const searchValue = searchField.value.toLowerCase();
    const fromAmount = parseFloat(fromAmountField.value) || -Infinity;
    const toAmount = parseFloat(toAmountField.value) || Infinity;

    Array.from(table.tBodies[0].rows).forEach(row => {
        const matchesFilters = filterColumns.every((colIndex, filterIndex) => {
            const cell = row.cells[colIndex];
            return selectedValues[filterIndex].length === 0 || selectedValues[filterIndex].includes(cell.textContent);
        });

        const matchesSearch = filterColumns.slice(0, 4).some(colIndex => {
            const cell = row.cells[colIndex];
            return cell.textContent.toLowerCase().includes(searchValue);
        });

        const amountCell = row.cells[5];
        const amount = parseFloat(amountCell.textContent) || 0;
        const matchesAmount = amount >= fromAmount && amount <= toAmount;

        const dateCell = row.cells[7];
        const rowDate = new Date(dateCell.textContent);
        rowDate.setHours(0, 0, 0, 0);
        const matchesDate = (!startDate || !endDate) || (rowDate >= startDate && rowDate <= endDate);

        row.style.display = matchesFilters && matchesSearch && matchesAmount && matchesDate ? "" : "none";
    });
}
/******* FILTER END  *********/

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

flatpickr("#calendar", {
    mode: "range",
    showMonths: 2, 
    onChange: function(selectedDates) {
        if (selectedDates.length === 2) {
            startDate = selectedDates[0];
            endDate = selectedDates[1];
            filterTable();
            updateDatesRange();
        }
    }
});

function updateDatesRange() {
    if (startDate && endDate) {
        var formattedStartDate = formatDate(startDate);
        var formattedEndDate = formatDate(endDate);
        document.getElementById('calendar').innerText = formattedStartDate + ' - ' + formattedEndDate;
    }
}

function formatDate(date) {
    const options = { day: 'numeric', month: 'short' };
    return date.toLocaleDateString('en-GB', options);
}

function addSearchOnFilters() {
    const originatingBankcSearchInput = document.getElementById('originating-bank-search');
    const originatingBankList = document.getElementById('originating-bank-list');
    
    originatingBankcSearchInput.addEventListener('input', function() {
        const searchTerm = originatingBankcSearchInput.value.toLowerCase();
        
        Array.from(originatingBankList.children).forEach(li => {
            const text = li.textContent.toLowerCase();
            if (text.includes(searchTerm)) {
                li.style.display = '';
            } else {
                li.style.display = 'none';
            }
        });
    });

    const originatorSearchInput = document.getElementById('originator-search');
    const originatorList = document.getElementById('originator-list');
    
    originatorSearchInput.addEventListener('input', function() {
        const searchTerm = originatorSearchInput.value.toLowerCase();
        
        Array.from(originatorList.children).forEach(li => {
            const text = li.textContent.toLowerCase();
            if (text.includes(searchTerm)) {
                li.style.display = '';
            } else {
                li.style.display = 'none';
            }
        });
    });

    const beneficiaryBankSearchInput = document.getElementById('beneficiary-bank-search');
    const beneficiaryBankList = document.getElementById('beneficiary-bank-list');
    
    beneficiaryBankSearchInput.addEventListener('input', function() {
        const searchTerm = beneficiaryBankSearchInput.value.toLowerCase();
        
        Array.from(beneficiaryBankList.children).forEach(li => {
            const text = li.textContent.toLowerCase();
            if (text.includes(searchTerm)) {
                li.style.display = '';
            } else {
                li.style.display = 'none';
            }
        });
    });

    const beneficiarySearchInput = document.getElementById('beneficiary-search');
    const beneficiaryList = document.getElementById('beneficiary-list');
    
    beneficiarySearchInput.addEventListener('input', function() {
        const searchTerm = beneficiarySearchInput.value.toLowerCase();
        
        Array.from(beneficiaryList.children).forEach(li => {
            const text = li.textContent.toLowerCase();
            if (text.includes(searchTerm)) {
                li.style.display = '';
            } else {
                li.style.display = 'none';
            }
        });
    });

    const currencySearchInput = document.getElementById('currency-search');
    const currencyList = document.getElementById('currency-list');
    
    currencySearchInput.addEventListener('input', function() {
        const searchTerm = currencySearchInput.value.toLowerCase();
        
        Array.from(currencyList.children).forEach(li => {
            const text = li.textContent.toLowerCase();
            if (text.includes(searchTerm)) {
                li.style.display = '';
            } else {
                li.style.display = 'none';
            }
        });
    });
}