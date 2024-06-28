window.addEventListener("load", () => {
    elements = document.getElementsByName("Amount")
    elements.forEach(element => {
        element.innerHTML = parseFloat(element.textContent.replace(/,/g, '')).toLocaleString()
    });
})

window.onload = function() {
    document.getElementById('add-policy-button').addEventListener('click', function(event) {
        event.preventDefault();
        fetch('/addpolicy', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
        }) 
        .then(response => response.text())
        .then(partialHTML => {
            var addPolicyDiv = document.getElementById('addPolicyDiv');
            if (!addPolicyDiv) {
                console.error("Element with id 'addPolicyDiv' not found.");
                return;
            }
    
            addPolicyDiv.innerHTML = "";
            addPolicyDiv.innerHTML = partialHTML;

            var modal = document.getElementById("add-policy-modal");
            var span = document.getElementsByClassName("close")[0];

            modal.style.display = "block";

            span.onclick = function() {
                modal.style.display = "none";
            }

            window.onclick = function(event) {
                if (event.target == modal) {
                    modal.style.display = "none";
                }
            }

        })
        .catch(error => console.error('Error fetching API:', error));
    });


}