window.addEventListener("load", () => {
    amounts = document.getElementsByName("Amount")
    amounts.forEach(element => {
        element.innerHTML = parseFloat(element.textContent.replace(/,/g, '')).toLocaleString()
    });

    const searchFiled = document.getElementById('searchField');
    console.log(document.getElementById('transactionsDiv'))
    searchFiled.addEventListener('input', function (e) {
        fetch('/home/searchTransaction?searchValue=' + e.target.value) 
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