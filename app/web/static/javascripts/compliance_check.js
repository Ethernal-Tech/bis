var data = {
    value: "",
    from: "",
    to: "",
    statusId: ""
}

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
var calendarBtn = document.getElementById('datesRange');
var calendarEl = document.getElementById('calendar');

var calendar = new FullCalendar.Calendar(calendarEl, {
    initialView: 'dayGridMonth', 
    height: '100%', 
    aspectRatio: 1.5,
    contentHeight: 'auto',
    selectable: true, 
    select: function(info) {
        var startDate = info.startStr;
        var start = new Date(startDate);
        var endDate = info.endStr;
        var end = new Date(endDate);
        end.setDate(end.getDate() - 1);
        endDate = end.toISOString().split('T')[0];

        var formattedStartDate = formatDate(startDate);
        var formattedEndDate = formatDate(endDate);

        document.getElementById('datesRange').innerText = `${formattedStartDate} - ${formattedEndDate}`;
        calendarWindow.style.display = 'none';
    }
});

function formatDate(dateStr) {
    const options = { day: '2-digit', month: 'short' };
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', options);
  }

calendarBtn.addEventListener('click', function() {
    calendarWindow.style.display = 'flex';
    calendar.render(); 
});

document.addEventListener('click', function(event) {
    if (!calendarWindow.contains(event.target) && !calendarBtn.contains(event.target)) {
        calendarWindow.style.display = 'none';
    }
});