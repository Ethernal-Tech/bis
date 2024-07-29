var view

window.addEventListener('load', function() {
    user = document.getElementById("home-user")
    logoutWindow = document.getElementById("home-logout-window")
    document.addEventListener('click', function(event){
        if (!logoutWindow.contains(event.target) && !user.contains(event.target)) {
            logoutWindow.style.display = "none"
        }
    })

    user.addEventListener('click', function(){
        logoutWindow.style.display = "flex"
    })

    view = document.getElementById('home-view')

    fetch("/compliancechecks", {
        method: 'POST',
    })
    .then(response => response.text())
    .then(partialHTML => {
        view.innerHTML = ""
        view.innerHTML = partialHTML

        //add compliancechecks-js.js after compliancechecks.html is loaded
        const script = document.createElement('script');
        script.src = '/app/web/static/javascripts/compliancechecks-js.js';
        script.defer = true;
        script.onload = function() {
            console.log('compliancechecks.js loaded successfully');
        };
        document.head.appendChild(script);
        //*****************
    })

    document.getElementById('home-compliance-check-view').addEventListener('click', function(){
        fetch("/compliancechecks", {
            method: 'POST',
        })
        .then(response => response.text())
        .then(partialHTML => {
            view.innerHTML = ""
            view.innerHTML = partialHTML

            //add compliancechecks-js.js after compliancechecks.html is loaded
            const script = document.createElement('script');
            script.src = '/app/web/static/javascripts/compliancechecks-js.js';
            script.defer = true;
            script.onload = function() {
                console.log('compliancechecks.js loaded successfully');
            };
            document.head.appendChild(script);
            //***************** 
        })
    })
    
    document.getElementById('home-analytics-view').addEventListener('click', function(){
        fetch("/analytics", {
            method: 'POST',
        })
        .then(response => response.text())
        .then(partialHTML => {
            view.innerHTML = ""
            view.innerHTML = partialHTML
        })
    })
    
    document.getElementById('home-rules-management-view').addEventListener('click', function(){
        fetch("/policies", {
            method: 'POST',
        })
        .then(response => response.text())
        .then(partialHTML => {
            view.innerHTML = ""
            view.innerHTML = partialHTML
        })
    })
})

