let mnu_TYPE = "";
let sel_TXT = "";
let mnu_ID = "";
let sub_UUID = 0;

function resetPage() {
    window.location.reload();
}

document.addEventListener("keydown", (event) => {
    if (event.key !== "Enter" || !mnu_TYPE) return;
    if (mnu_TYPE === "titles") {
        TitlesUpd();
    } else {
        subAddUpd(mnu_ID);
    }
});

function isEdit(mnu_type, mnu_id) {
    mnu_TYPE = mnu_type;
    mnu_ID = mnu_id;
    const txtElem = document.getElementById(`_txt_${mnu_ID}`);
    sel_TXT = txtElem ? txtElem.value : "";
}

// MENUS //
function subSel(sub_uuid, mnu_id, mnu_type) {
    mnu_ID = mnu_id;
    sub_UUID = parseInt(sub_uuid, 10);
    mnu_TYPE = mnu_id === 500 ? "titles" : mnu_type;

    const optElem = document.getElementById(`_opt_${sub_UUID}`);
    const txtElem = document.getElementById(`_txt_${mnu_ID}`);
    const subElem = document.getElementById(`_sub_${mnu_ID}`);

    sel_TXT = optElem ? optElem.value : "";
    if (txtElem) txtElem.value = sel_TXT;
    if (subElem) subElem.value = sub_UUID;
}

async function TitlesUpd() {
    const txtElem = document.getElementById("_txt_500");
    const subElem = document.getElementById("_sub_500");
    const _txt = txtElem ? txtElem.value : "";
    const _uuid = subElem ? parseInt(subElem.value, 10) : 0;

    if (!_uuid || !_txt) return;

    const data = {
        mnu_id: 500,
        sub_uuid: _uuid,
        mnu_title: _txt,
    };

    try {
        const response = await fetch("/title/upd", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        const responseData = await response.json();
        if (!response.ok) throw new Error(responseData.Error);
        window.location.reload();
    } catch (error) {
        ShowMsg("Upd Titles failed: " + error.message, "error");
        window.location.reload();
    }
}

async function subAddUpd(mnu_id) {
    const txtElem = document.getElementById(`_txt_${mnu_id}`);
    const subElem = document.getElementById(`_sub_${mnu_id}`);
    let txt = txtElem ? txtElem.value.trim() : "";

    if (!txt) {
        alert("Nothing to add or update" + mnu_id);
        return;
    }

    const sub_uuid = subElem ? parseInt(subElem.value, 10) : 0;
    const _mnuid = mnu_id.toString().substring(3);

    const data = {
        mnu_id: parseInt(_mnuid, 10),
        sub_uuid,
        val: txt,
    };

    const url = `/${mnu_TYPE}/addupd`;
    try {
        const response = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        const responseData = await response.json();
        if (!response.ok) throw new Error(responseData.Error);
        window.location.reload();
    } catch (error) {
        alert(`Add/Update ${mnu_TYPE} failed: ${error.message}`);
    }
}

async function lstDel(mnu_id) {
    const txtElem = document.getElementById(`_txt_${mnu_ID}`);
    const subElem = document.getElementById(`_sub_${mnu_ID}`);
    const txt = txtElem ? txtElem.value : "";
    const sub_uuid = subElem ? parseInt(subElem.value, 10) : 0;

    if (!txt) {
        alert("Nothing to delete");
        return;
    }

    if (!confirm(`Delete '${txt}'?`)) return;

    const data = {
        sub_uuid: sub_uuid,
    };

    try {
        const response = await fetch(`/${mnu_TYPE}/delete`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        const responseData = await response.json();
        if (response.status !== 200) throw new Error(responseData.error);
        window.location.reload();
    } catch (error) {
        alert(`Delete ${mnu_TYPE} failed: ${error.message}`);
        window.location.reload();
    }
}
