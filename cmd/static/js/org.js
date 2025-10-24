async function orgAddUpd(uuid) {
    const org = {
        UUID: uuid,
        Name: document.getElementById("_name").value,
        Note: document.getElementById("_note").value
    };

    try {
        const response = await fetch("/v/org/addupd", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(org)
        });

        if (response.ok) {
            window.location.href = "/v/orgs";
        } else {
            ShowMsg("error", "Error updating organization");
        }
    } catch {
        ShowMsg("error", "Network error updating organization");
    }
}

async function deleteOrg(uuid) {
    if (!confirm("Are you sure you want to delete this organization? This action cannot be undone.")) return;

    try {
        const response = await fetch(`/v/org/${uuid}`, { method: "DELETE" });

        if (response.ok) {
            ShowMsg("info", "Organization deleted successfully");
            window.location.href = "/v/orgs";
        } else {
            ShowMsg("error", "Error deleting organization");
        }
    } catch {
        ShowMsg("error", "Network error deleting organization");
    }
}

async function modifyMember(url, orgUUID, userUUID, successMsg, errorMsg) {
    try {
        const response = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ org_id: Number(orgUUID), user_id: Number(userUUID) })
        });

        if (response.ok) {
            ShowMsg("info", successMsg);
            window.location.reload();
        } else {
            ShowMsg("error", errorMsg);
        }
    } catch {
        ShowMsg("error", `Network error ${errorMsg.toLowerCase()}`);
    }
}

function AddMember(orgUUID, userUUID) {
    return modifyMember("/v/org/members/add", orgUUID, userUUID, "Member added successfully", "Error adding member");
}

function RemoveMember(orgUUID, userUUID) {
    return modifyMember("/v/org/members/rem", orgUUID, userUUID, "Member removed successfully", "Error removing member");
}
