
document.addEventListener("keydown", function (event) {
    
    if (event.key == "Enter") {
        
    }
});

function ShowText() {
    var _txt = document.getElementById("_text").value;
    if (_txt == "") {
        ShowMsg("Please enter text.");
        return;
    }
    document.getElementById("_show").innerHTML = _txt;

}