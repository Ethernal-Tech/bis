var view

window.addEventListener('load', function() {
    user = document.getElementById("home-user")
    logoutWindow = document.getElementById("home-logout-window")
    document.addEventListener('click', function(event){
    if (!logoutWindow.contains(event.target) && !user.contains(event.target)) {
        logoutWindow.classList.remove("visible");
    }
    })

    user.addEventListener('click', function(e){
        logoutWindow.classList.toggle("visible");
    })

    view = document.getElementById('home-view')

    fetch("/complianceCheckIndex", {
        method: 'GET',
    })
    .then(response => response.text())
    .then(partialHTML => {
        view.innerHTML = ""
        view.innerHTML = partialHTML
        loadScript('compliance_checks')
    })

    document.getElementById('home-compliance-check-view').addEventListener('click', function(){
        fetch("/complianceCheckIndex", {
            method: 'GET',
        })
        .then(response => response.text())
        .then(partialHTML => {
            view.innerHTML = ""
            view.innerHTML = partialHTML
            loadScript('compliance_checks')
        })
    })
    
    var sideBarAnalytics = document.getElementById('home-analytics-view')
    if (sideBarAnalytics) {
        sideBarAnalytics.addEventListener('click', function(){
            fetch("/analytics", {
                method: 'POST',
            })
            .then(response => response.text())
            .then(partialHTML => {
                view.innerHTML = ""
                view.innerHTML = partialHTML
                loadScript('analytics')
            })
        })
    }

    document.getElementById('home-rules-management-view').addEventListener('click', function(){
        fetch("/policies", {
            method: 'POST',
        })
        .then(response => response.text())
        .then(partialHTML => {
            view.innerHTML = ""
            view.innerHTML = partialHTML
            loadScript('policies')
        })
    })
})

function loadScript(partial) {
    const existingScript = document.getElementById('dynamicScript');
    if (existingScript) {
        existingScript.remove();
    }

    const script = document.createElement('script');
    script.id = 'dynamicScript';
    script.src = `/app/web/static/javascripts/${partial}.js`;
    document.body.appendChild(script);
}
