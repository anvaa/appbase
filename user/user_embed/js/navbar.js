
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

function backToApp() {
    window.location.href = "/app";
}

async function Logout() {
    window.location.href = "/logout";
}


function ShowMsg(val) {
    document.getElementById("_msg").innerHTML = val;
    document.getElementById("_msg").style.display = "block";
    document.getElementById("_msg").style.color = "black";
    document.getElementById("_msg").style.fontSize = "16px";
    document.getElementById("_msg").style.width = "fit-content";

    // resetFields();

    setTimeout(function() {
        document.getElementById("_msg").style.display = "none";
        window.location.reload();
    }, 3000);
    return;
}