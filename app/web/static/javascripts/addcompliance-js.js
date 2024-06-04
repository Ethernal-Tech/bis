window.onload = function () {
    document.getElementById("cancel-button").addEventListener("click", downgradeView)
    document.getElementById("next-button").addEventListener("click", upgradeView)
}

let currentView = 1

function downgradeView() {
    let view1CL = document.getElementById("view-1").classList
    let view2CL = document.getElementById("view-2").classList
    let view2indCL = document.querySelector(".view-indicator > div:nth-child(2)").classList
    let view3CL = document.getElementById("view-3").classList
    let view3indCL = document.querySelector(".view-indicator > div:nth-child(3)").classList
    
    if (currentView == 1) {
        window.location.href = "/home"
    } else if (currentView == 2) {
        currentView--
        view2CL.remove("display")
        view2CL.add("not-display")
        view2indCL.remove("add-color")
        view2indCL.add("remove-color")
        view1CL.remove("not-display")
        view1CL.add("display")
    } else if (currentView == 3) {
        currentView--
        view3CL.remove("display")
        view3CL.add("not-display")
        view3indCL.remove("add-color")
        view3indCL.add("remove-color")
        view2CL.remove("not-display")
        view2CL.add("display")
    }
}

function upgradeView() {
    let view1CL = document.getElementById("view-1").classList
    let view2CL = document.getElementById("view-2").classList
    let view2indCL = document.querySelector(".view-indicator > div:nth-child(2)").classList
    let view3CL = document.getElementById("view-3").classList
    let view3indCL = document.querySelector(".view-indicator > div:nth-child(3)").classList

    if (currentView == 1) {
        currentView++
        view1CL.remove("display")
        view1CL.add("not-display")
        view2indCL.remove("remove-color")
        view2indCL.add("add-color")
        view2CL.remove("not-display")
        view2CL.add("display")
    } else if (currentView == 2) {
        currentView++
        view2CL.remove("display")
        view2CL.add("not-display")
        view3indCL.remove("remove-color")
        view3indCL.add("add-color")
        view3CL.remove("not-display")
        view3CL.add("display")
    } else if (currentView == 3) {
        // submit compliance check
    }
}

