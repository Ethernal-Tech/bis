window.addEventListener("load", () => {
    amounts = document.getElementsByName("Amount")
    amounts.forEach(element => {
        console.log(element.textContent)
        element.innerHTML = parseFloat(element.textContent.replace(/,/g, '')).toLocaleString()
    });
})