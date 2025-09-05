
async function orgAddUpd(uuid) {
    const org = {
        UUID: uuid,
        Name: document.getElementById("_name").value,
        Note: document.getElementById("_note").value
    };

    const response = await fetch(`/v/org/addupd`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(org)
    });

    if (response.ok) {
        window.location.href = "/v/orgs";
    } else {
        ShowMsg("error","Error updating organization");
    }
}

async function deleteOrg(uuid) {

    if (!confirm("Are you sure you want to delete this organization? This action cannot be undone.")) {
        return;
    }

    const response = await fetch(`/v/org/${uuid}`, {
        method: "DELETE"
    });

    if (response.ok) {
        ShowMsg("info","Organization deleted successfully");
        window.location.href = "/v/orgs";
    } else {
        ShowMsg("error","Error deleting organization");
    }
}
