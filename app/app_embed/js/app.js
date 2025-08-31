
document.addEventListener("keydown", function (event) {
    
    if (event.key == "Enter") {
        
        var inputField = document.getElementById("_text");
        if (inputField === document.activeElement) {
            event.preventDefault(); // Prevent default action for Enter key
            ShowText(); // Call the ShowText function when Enter is pressed
        }
    }

    if (event.key == "Escape") {
        // reload page
        window.location.reload();    
    }

});

function ShowText() {
    var _txt = document.getElementById("_text").value;
    if (_txt == "") {
        ShowMsg("error", "Please enter text.");
        return;
    }
    if (_txt.length > 20) {
        ShowMsg("info", "Text is too long. Please limit to 100 characters.");
        return;
    }
    if (_txt.length < 3) {
        ShowMsg("info", "Text is too short. Please enter at least 3 characters.");
        return;
    }
    if (badword(_txt)) {
        ShowMsg("error", "Text contains a bad word. Please avoid using inappropriate language.");
        return;
    }
    document.getElementById("_show").innerHTML = "<b>You wrote:</b> "+_txt;

}

function ShowMsg(type = "info", msg) {
    const msgColors = {
        error: "red;font-weight: bold",
        success: "green",
        info: "blue"
    };
    const colorStyle = msgColors[type] || msgColors.info;
    const formattedMsg = `<span style='color: ${colorStyle}'>${msg}</span>`;

    const msgElem = document.getElementById("_msg");
    msgElem.style.fontFamily = "monospace";
    msgElem.style.fontSize = "18px";
    msgElem.innerHTML = formattedMsg;

    setTimeout(() => {
        msgElem.innerHTML = "";
    }, 4000);
}

function badword(txt) {
    // Example of bad words check
    badwordarray = ["badword", "badword2", "badword3", "badword4"];
    for (var i = 0; i < badwordarray.length; i++) {
        if (txt.includes(badwordarray[i])) {
            return true;
        }
    }
    return false;
}