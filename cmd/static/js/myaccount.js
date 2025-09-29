async function NewPsw(_uuid) {
    const _psw1 = document.getElementById("_psw1").value;
    const _psw2 = document.getElementById("_psw2").value;

    if (!validatePasswords(_psw1, _psw2)) return;

    const userData = { uuid: _uuid, psw1: _psw1, psw2: _psw2 };

    try {
        const response = await fetch("/user/psw", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(userData),
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || "Failed to change password");
        }

        ShowMsg("Password changed successfully", "success");
    } catch (error) {
        ShowMsg(`Change password failed: ${error.message}`, "error");
    }
}

function validatePasswords(password, password2) {
    if (!password || !password2) {
        ShowMsg("Passwords cannot be empty", "info");
        return false;
    }
    if (password !== password2) {
        ShowMsg("Passwords do not match", "info");
        return false;
    }
    if (password.length < 8) {
        ShowMsg("Passwords must be at least 8 characters long", "info");
        return false;
    }
    return true;
}
