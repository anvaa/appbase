
document.addEventListener("keydown", function (event) {  
    
    if (event.ctrlKey && event.key === "1") {
        Start();
    }

    if (event.ctrlKey && event.key === "q") {
        Logout();
    }
});

function Start() {
    window.location.href = "/app/start";
}

function Tools() {
    window.location.href = "/tools/titles";
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

// Tools menus
function toolsTitles() {
    window.location.href = "/tools/titles";
}

function toolsStatus() {
    window.location.href = "/tools/status";
}

function toolsTypes() {
    window.location.href = "/tools/types";
}


async function fetchPage(url) {
    // Get url
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        window.location.href = url;
    } catch (error) {
        console.error("Error fetching page:", error);
        ShowMsg("Failed to load page", "error");
    }
}



function ShowMsg(val, type) {
    const msgEl = document.getElementById("_msg");
    if (!msgEl) return;

    // loxck page for input
    document.body.style.pointerEvents = "none";

    msgEl.innerHTML = val;
    msgEl.style.margin = "2px";
    msgEl.style.fontSize = "large";
    msgEl.style.fontWeight = "bold";
    msgEl.style.fontFamily = "monospace";
    msgEl.style.width = "fit-content";
    msgEl.style.padding = "2px";
    msgEl.style.borderRadius = "5px";

    let borderColor;
    switch (type) {
        case "error":
            borderColor = "red";
            break;
        case "success":
            borderColor = "green";
            break;
        case "info":
            borderColor = "blue";
            break;
        default:
            borderColor = "black";
    }
    msgEl.style.border = `3px solid ${borderColor}`;
    msgEl.style.color = `${borderColor}`;

    setTimeout(() => {
        // fade out the message
        msgEl.style.transition = "opacity 0.5s";
        msgEl.style.opacity = "0";
        // after the fade out, hide the message
        setTimeout(() => {
            msgEl.style.display = "none";
            msgEl.style.opacity = "1"; // reset opacity for next message
            document.body.style.pointerEvents = "auto"; // unlock page for input
        }, 500);
    }, 3000);

}
