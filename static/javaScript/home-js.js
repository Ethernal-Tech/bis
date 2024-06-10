window.onload = function() {
    fetch('/home/Transactions') 
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
};

window.addEventListener("load", () => {
    amounts = document.getElementsByName("Amount")
    amounts.forEach(element => {
        element.innerHTML = parseFloat(element.textContent.replace(/,/g, '')).toLocaleString()
    });

    const searchFiled = document.getElementById('searchField');
    searchFiled.addEventListener('input', function (e) {
        fetch('/home/Transactions?searchValue=' + e.target.value) 
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