
const messageElement = document.getElementById("_message");

function mngUsers() {
    window.location.href = "/v/users";
}

function newUsers() {
    window.location.href = "/v/newusers";
}

function addUser() {
    window.location.href = "/signup/121209";
}

function mngOrgs() {
    window.location.href = "/v/orgs";
}

function mngDatabase() {
    window.location.href = "/v/database";
}

function backToApp() {
    window.location.href = "/app";
}

async function Logout() {
    window.location.href = "/logout";
}


function ShowMsg(type, val) {
    const msgEl = document.getElementById("_msg");
    if (!msgEl) return;

    document.body.style.pointerEvents = "none";

    msgEl.innerHTML = val;
    Object.assign(msgEl.style, {
        margin: "2px",
        fontSize: "large",
        fontWeight: "bold",
        fontFamily: "monospace",
        width: "fit-content",
        padding: "2px",
        borderRadius: "5px",
        border: `3px solid ${getBorderColor(type)}`,
        color: getBorderColor(type),
        display: "block",
        opacity: "1",
        transition: ""
    });

    setTimeout(() => {
        msgEl.style.transition = "opacity 0.5s";
        msgEl.style.opacity = "0";
        setTimeout(() => {
            msgEl.style.display = "none";
            msgEl.style.opacity = "1";
            document.body.style.pointerEvents = "auto";
        }, 500);
    }, 3000);
}

function getBorderColor(type) {
    switch (type) {
        case "error": return "red";
        case "success": return "green";
        case "info": return "blue";
        default: return "black";
    }
}