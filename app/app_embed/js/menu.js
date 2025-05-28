
document.addEventListener("keydown", function (event) {
    if (event.key === "Escape") {
        resetPage();
    }

    // Ctrl + N
    if (event.ctrlKey && event.key === "n") {
        // NewItem();
    }  
    
    if (event.ctrlKey && event.key === "1") {
        Home();
    }
    
    if (event.ctrlKey && event.key === "2") {
        Search();
    }

    if (event.ctrlKey && event.key === "3") {
        Stats();
    }

    if (event.ctrlKey && event.key === "q") {
        Logout();
    }
});

function Start() {
    window.location.href = "/app";
}

function Tools() {
    window.location.href = "/app/tools";
}

function Users() {
    window.location.href = "/v/users";
}

function MyAccount() {
    window.location.href = "/v/myaccount";
}

function AppInfo() {
    window.location.href = "/info";
}

function Logout() {
    window.location.href = "/logout";
}



function ShowMsg(val,type) {
    document.getElementById("_msg").innerHTML = val;
    document.getElementById("_msg").style.display = "block";
    document.getElementById("_msg").style.color = "black";
    document.getElementById("_msg").style.fontSize = "16px";
    document.getElementById("_msg").style.width = "fit-content";

    if (type === "error") {
        document.getElementById("_msg").style.backgroundColor = "red";
    } else if (type === "success") {
        document.getElementById("_msg").style.backgroundColor = "green";
    } else if (type === "info") {
        document.getElementById("_msg").style.backgroundColor = "blue";
    }

    setTimeout(function() {
        document.getElementById("_msg").style.display = "none";
        window.location.reload();
    }, 3000);
    return;
}