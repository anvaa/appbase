async function Delete(uuid) {
    const username = document.getElementById("_username").value;
    const confirmMsg = `Are you sure you want to delete user ${username}?\n\nYou can just remove auth to deny access.`;
    if (!confirm(confirmMsg)) return;

    try {
        const response = await fetch("/user/delete", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ uuid }),
        });
        if (!response.ok) throw new Error("Delete failed");
        window.location.href = "/v/users";
    } catch (error) {
        ShowMsg(error.message);
    }
}

async function setPsw(uuid) {
    const psw1 = document.getElementById("_psw1").value;
    const psw2 = document.getElementById("_psw2").value;
    if (!validatePasswords(psw1, psw2)) return;

    try {
        const response = await fetch("/user/psw", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ uuid, psw1, psw2 }),
        });
        if (!response.ok) throw new Error("Change password failed");
        ShowMsg("Password changed successfully");
    } catch (error) {
        ShowMsg(error.message);
    }
}

function validatePasswords(psw1, psw2) {
    if (!psw1 || !psw2) {
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

async function AuthLevel(uuid) {
    const authlevel = parseInt(document.getElementById("_authlevel").value, 10);
    try {
        const response = await fetch("/user/authlevel", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ uuid, authlevel }),
        });
        if (!response.ok) throw new Error("Change role failed");
        window.location.reload();
    } catch (error) {
        ShowMsg(error.message);
    }
}

function authSelect(authlevel) {
    const currentAuthLevel = parseInt(document.getElementById("_authlevel").value, 10);
    if (authlevel >= currentAuthLevel) {
        alert("You cannot assign an AuthLevel equal to or higher than your own.");
    }
}

async function Orgname(uuid) {
    const orgname = document.getElementById("_orgname").value;
    try {
        const response = await fetch("/user/org", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ uuid, orgname }),
        });
        if (!response.ok) throw new Error("Change Org failed");
        window.location.reload();
    } catch (error) {
        ShowMsg(error.message);
    }
}

async function setAuth(uuid, auth) {
    const boolAuth = auth === "true" ? true : auth === "false" ? false : null;
    if (boolAuth === null) {
        ShowMsg("Invalid auth value");
        return;
    }
    try {
        const response = await fetch("/user/auth", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ uuid, auth: boolAuth }),
        });
        if (!response.ok) throw new Error("Change auth failed");
        window.location.reload();
    } catch (error) {
        ShowMsg(error.message);
    }
}
