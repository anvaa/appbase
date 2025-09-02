
async function Delete(_uuid) {

    var _email = document.getElementById("_email").value;

    const verify = confirm("Are you sure you want to delete user " + _email + 
                            "? \n\nYou can just remove auth to deny access.", "Delete user");
    if (!verify) {
        return;
    }

    var userData = {
        uuid: _uuid,
    };

    try {
        const response = await fetch("/user/delete", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(userData),
        });
    
        if (!response.ok) {
            ShowMsg(error.message);
            return;
        }
    
        window.location.href = "/v/users";
        } catch (error) {
        ShowMsg("Delete failed: " + error.message);

    }
}

async function setPsw(_uuid) {

    var _psw1 = document.getElementById("_psw1").value;
    var _psw2 = document.getElementById("_psw2").value;

    if (!validatePasswords(_psw1, _psw2)) {
        return; // Message is set inside the validatePasswords function
    }

    var userData = {
        uuid: _uuid,
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
            return;
        }
    
            ShowMsg("Password changed successfully");
        } catch (error) {
            ShowMsg("Change password failed: " + error.message);
        
    }
}

function validatePasswords(psw1, psw2) {

    if (psw1 === "" || psw2 === "") {
        ShowMsg("Passwords cannot be empty");
    return false;
    }

    if (psw1 !== psw2) {
        ShowMsg("Passwords do not match");
    return false;
    }

    if (psw1.length < 8) {
        ShowMsg("Passwords must be at least 8 characters long");
    return false;
    }

    return true;
}

async function AuthLevel(_uuid) {

    var _authlevel = document.getElementById("_authlevel").value;

    var userData = {
        uuid: _uuid,
        authlevel: parseInt(_authlevel),
    };
    // alert(JSON.stringify(userData));
    try {
        const response = await fetch("/user/authlevel", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(userData),
        });
    
        if (!response.ok) {
            ShowMsg(error.message);
            return;
        } 

        window.location.reload();
    } catch (error) {
        ShowMsg("Change role failed: " + error.message);
    }
}

function authSelect(_authlevel) {
    var currentAuthLevel = parseInt(document.getElementById("_authlevel").value, 10);
    if (_authlevel >= currentAuthLevel) {
        alert("You cannot assign an AuthLevel equal to or higher than your own.");
    }
}

async function Orgname(_uuid) {

    var _orgname = document.getElementById("_orgname").value;
    
    var userData = {
        uuid: _uuid,
        orgname: _orgname,
    };
    // alert(JSON.stringify(userData));
    try {
        const response = await fetch("/user/org", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(userData),
        });
    
        if (!response.ok) {
            ShowMsg(error.message);
            return;
        } 

        window.location.reload();
    } catch (error) {
        ShowMsg("Change Org failed: " + error.message);

    }
}

async function setAuth(_uuid, _auth) {
    
    // conv_auth to bool
    if (_auth == "true") {
        _auth = true;
    } else if (_auth == "false") {
        _auth = false;
    } else {
        ShowMsg("Invalid auth value");
        return;
    }

    var data = {
        uuid: _uuid,
        auth: _auth,
    };
    // alert(JSON.stringify(data));
    try {
        const response = await fetch("/user/auth", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });

        if (!response.ok) {
            ShowMsg(error.message);
            return;
        }

        window.location.reload();
    } catch (error) {
        ShowMsg("Change auth failed: " + error.message);

    }
}