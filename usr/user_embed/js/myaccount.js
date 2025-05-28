
async function setPsw() {

    var _id = document.getElementById("_id").value;
    var _psw1 = document.getElementById("_psw1").value;
    var _psw2 = document.getElementById("_psw2").value;

    if (!validatePasswords(_psw1, _psw2)) {
        return; // Message is set inside the validatePasswords function
    }

    var userData = {
        id: _id,
        psw1: _psw1,
        psw2: _psw2,
    };

    try {
        const response = await fetch("/user/psw", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(userData),
        });
    
        if (!response.ok) {
            ShowMsg(error.message);
        }
    
            ShowMsg("Password changed successfully");
        } catch (error) {
            ShowMsg("Change password failed: " + error.message);
        
    }
}

function validatePasswords(password, password2) {

    if (password === "" || password2 === "") {
        ShowMsg("Passwords cannot be empty");
    return false;
    }

    if (password !== password2) {
        ShowMsg("Passwords do not match");
    return false;
    }

    if (password.length < 8) {
        ShowMsg("Passwords must be at least 8 characters long");
    return false;
    }

    return true;
}