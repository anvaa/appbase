let mnu_TYPE = "";
let sel_TXT = "";
let mnu_ID = "";
let sub_UUID = 0;

function resetPage() {
    window.location.reload();
}

document.addEventListener("keydown", (event) => {
    if (event.key === "Enter" && mnu_TYPE) {
        mnu_TYPE === "titles" ? TitlesUpd() : subAddUpd(mnu_ID);
    }
});

function isEdit(mnu_type, mnu_id) {
    mnu_TYPE = mnu_type;
    mnu_ID = mnu_id;
    const txtElem = document.getElementById(`_txt_${mnu_ID}`);
    sel_TXT = txtElem?.value || "";
}

// MENUS //
function subSel(sub_uuid, mnu_id, mnu_type) {
    mnu_ID = mnu_id;
    sub_UUID = Number(sub_uuid) || 0;
    mnu_TYPE = mnu_id === 500 ? "titles" : mnu_type;

    const optElem = document.getElementById(`_opt_${sub_UUID}`);
    const txtElem = document.getElementById(`_txt_${mnu_ID}`);
    const subElem = document.getElementById(`_sub_${mnu_ID}`);

    sel_TXT = optElem?.value || "";
    if (txtElem) txtElem.value = sel_TXT;
    if (subElem) subElem.value = sub_UUID;
    document.getElementById("_tooltest").innerHTML = `Selected: ${sel_TXT} _opt_${sub_UUID}`;
}

async function TitlesUpd() {
    const txtElem = document.getElementById("_txt_500");
    const subElem = document.getElementById("_sub_500");
    const _txt = txtElem?.value || "";
    const _uuid = Number(subElem?.value) || 0;

    if (!_uuid || !_txt) return;

    const data = { mnu_id: 500, sub_uuid: _uuid, mnu_title: _txt };

    try {
        const response = await fetch("/title/upd", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        const responseData = await response.json();
        if (!response.ok) throw new Error(responseData.Error);
        resetPage();
    } catch (error) {
        ShowMsg("Upd Titles failed: " + error.message, "error");
        resetPage();
    }
}

async function subAddUpd(mnu_id) {
    const txtElem = document.getElementById(`_txt_${mnu_id}`);
    const subElem = document.getElementById(`_sub_${mnu_id}`);
    const idElem = document.getElementById(`_id_${mnu_id}`);
    const _txt = txtElem?.value.trim() || "";

    if (!_txt) {
        alert("Nothing to add or update");
        return;
    }

    const _sub_uuid = Number(subElem?.value) || 0;
    const _mnuid = Number(idElem?.value) || 0;

    const data = { mnu_id: _mnuid, sub_uuid: _sub_uuid, val: _txt };
    const url = `/${mnu_TYPE}/addupd`;

    try {
        const response = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        const responseData = await response.json();
        if (!response.ok) throw new Error(responseData.Error);
        resetPage();
    } catch (error) {
        alert(`Add/Update ${mnu_TYPE} failed: ${error.message}`);
    }
}

async function lstDel(mnu_id) {
    const txtElem = document.getElementById(`_txt_${mnu_ID}`);
    const subElem = document.getElementById(`_sub_${mnu_ID}`);
    const txt = txtElem?.value || "";
    const sub_uuid = Number(subElem?.value) || 0;

    if (!txt) {
        alert("Nothing to delete");
        return;
    }

    if (!confirm(`Delete '${txt}'?`)) return;

    const data = { sub_uuid };

    try {
        const response = await fetch(`/${mnu_TYPE}/delete`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        const responseData = await response.json();
        if (response.status !== 200) throw new Error(responseData.error);
        resetPage();
    } catch (error) {
        alert(`Delete ${mnu_TYPE} failed: ${error.message}`);
        resetPage();
    }
}
