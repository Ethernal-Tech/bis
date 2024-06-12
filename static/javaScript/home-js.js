window.onload = function() {
    var searchData = {
        "value": null,
        "from": null,
        "to": null,
        // "amount": null
    };

    fetch('/transactions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(searchData)
        }) 
        .then(response => response.text())
        .then(partialHTML => {
            var transactionsDiv = document.getElementById('transactionsDiv');
            if (!transactionsDiv) {
                console.error("Element with id 'transactionsDiv' not found.");
                return;
            }

            transactionsDiv.innerHTML = "";
            transactionsDiv.innerHTML = partialHTML;

            amounts = document.getElementsByName("Amount")
            amounts.forEach(element => {
                element.innerHTML = parseFloat(element.textContent.replace(/,/g, '')).toLocaleString()
            });

        })
        .catch(error => console.error('Error fetching API:', error));

        //****modal****
        var modal = document.getElementById("advancedSearchModal");
        var btn = document.getElementById("advancedSearchBtn");
        var span = document.getElementsByClassName("close")[0];
        
        btn.onclick = function() {
            modal.style.display = "block";
        }
        
        span.onclick = function() {
            modal.style.display = "none";
        }
        
        window.onclick = function(event) {
            if (event.target == modal) {
                modal.style.display = "none";
            }
        }
        
        document.getElementById("searchBtn").onclick = function() {
            var searchValue = document.getElementById("searchValue").value;
            var dateFrom = document.getElementById("dateFrom").value;
            var dateTo = document.getElementById("dateTo").value;
            // var amount = document.getElementById("amount").value;
            
            var searchData = {
                "value": searchValue,
                "from": dateFrom.replace('T', ' '),
                "to": dateTo.replace('T', ' '),
                // "amount": amount
            };

            console.log("Search parameters: ", { searchData });
            
            fetch("/transactions", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(searchData)
            })
            // .then(response => {
            //     if (!response.ok) {
            //         throw new Error('Network response was not ok');
            //     }})
            .then(response => response.text())
            .then(partialHTML => {
                var transactionsDiv = document.getElementById('transactionsDiv');
                if (!transactionsDiv) {
                    console.error("Element with id 'transactionsDiv' not found.");
                    return;
                }

                transactionsDiv.innerHTML = "";
                transactionsDiv.innerHTML = partialHTML;

                amounts = document.getElementsByName("Amount")
                amounts.forEach(element => {
                    element.innerHTML = parseFloat(element.textContent.replace(/,/g, '')).toLocaleString()
                });
                
                // Close the modal after search
                modal.style.display = "none";
            })
        }
        //****modal****
};

window.addEventListener("load", () => {
    amounts = document.getElementsByName("Amount")
    amounts.forEach(element => {
        element.innerHTML = parseFloat(element.textContent.replace(/,/g, '')).toLocaleString()
    });

    const searchFiled = document.getElementById('searchField');
    searchFiled.addEventListener('input', function (e) {
        var searchData = {
            "value": e.target.value,
            "from": null,
            "to": null,
            // "amount": null
        };
        fetch('/transactions', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(searchData)
            }) 
            .then(response => response.text())
            .then(partialHTML => {
                var transactionsDiv = document.getElementById('transactionsDiv');
                if (!transactionsDiv) {
                    console.error("Element with id 'transactionsDiv' not found.");
                    return;
                }
    
                transactionsDiv.innerHTML = "";
                transactionsDiv.innerHTML = partialHTML;

                amounts = document.getElementsByName("Amount")
                amounts.forEach(element => {
                    element.innerHTML = parseFloat(element.textContent.replace(/,/g, '')).toLocaleString()
                });

            })
            .catch(error => console.error('Error fetching API:', error));
    });
})