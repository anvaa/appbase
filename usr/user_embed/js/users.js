
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

async function Role(_uuid) {

    var _role = document.getElementById("_role").value;
    
    var userData = {
        uuid: _uuid,
        role: _role,
    };

    try {
        const response = await fetch("/user/role", {
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

async function Orgname(_uuid) {

    var _orgname = document.getElementById("_orgname").value;
    
    var userData = {
        uuid: _uuid,
        orgname: _orgname,
    };

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
        uuid: parseInt(_uuid),
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