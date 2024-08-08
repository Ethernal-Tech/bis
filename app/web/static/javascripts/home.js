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

    document.getElementsByClassName('home-hamburger')[0].addEventListener('click', function() {
        let homeBody = document.getElementsByClassName('home-body')[0];
        let homeSideBar = document.getElementsByClassName('home-side-bar')[0];
        let homeView = document.getElementById('home-view');
        if (homeBody.classList.contains('hide-home-side-bar')) {
            homeBody.classList.remove('hide-home-side-bar');
            homeBody.classList.add('show-home-side-bar');
            homeView.classList.remove('expend-home-view');
            homeView.classList.add('shrink-home-view');
            setTimeout(() => {
                homeSideBar.classList.remove('hidden');
                homeSideBar.classList.add('visible');
            }, 100);
        }
        else {
            homeBody.classList.remove('show-home-side-bar');
            homeBody.classList.add('hide-home-side-bar');
            homeView.classList.remove('shrink-home-view');
            homeView.classList.add('expend-home-view');
            setTimeout(() => {
                homeSideBar.classList.remove('visible');
                homeSideBar.classList.add('hidden');
            }, 300);
        }
    });
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
