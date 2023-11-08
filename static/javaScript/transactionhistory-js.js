window.addEventListener("load", () => {
    elements = document.getElementsByName("Amount")
    elements.forEach(element => {
        element.innerHTML = parseFloat(element.textContent.replace(/,/g, '')).toLocaleString()
    });
})