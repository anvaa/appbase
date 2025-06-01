
var mnu_TYPE = "";
var sel_TXT = "";
var mnu_ID = "";
var sub_UUID = 0;


function resetPage() {
    window.location.reload();
    return;
};

document.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
        
        if (mnu_TYPE == "") {
            return;
        }
        
        if (mnu_TYPE == "titles") {
            TitlesUpd(mnu_ID);
            return;
        } else {
            subAddUpd(mnu_ID);
            return;
        }
    }
});

function isEdit(mnu_type, mnu_id) {
    mnu_TYPE = mnu_type;
    mnu_ID = mnu_id;
    
    sel_TXT = document.getElementById("_txt_"+mnu_ID).value;

}


// MENUS //
function subSel(sub_uuid, mnu_id, mnu_type) {
    mnu_ID = mnu_id;
    sub_UUID = parseInt(sub_uuid);
    mnu_TYPE = mnu_type;

    if (mnu_ID === 500) {
        mnu_TYPE = "titles";
    }

    // alert("subSel: "+mnu_ID+" "+sub_UUID+" "+mnu_TYPE);
    
    sel_TXT = document.getElementById("_opt_"+sub_UUID).value;
    document.getElementById("_txt_"+mnu_ID).value = sel_TXT;
    document.getElementById("_sub_"+mnu_ID).value = sub_UUID;

}

async function TitlesUpd() {
    var _txt = document.getElementById("_txt_500").value;
    var _uuid = parseInt(document.getElementById("_sub_500").value);

    if (_uuid == 0 || _txt == "") {
        return; 
    }
    
    var data = {
        "mnu_id": 500,
        "sub_uuid": _uuid,
        "mnu_title": _txt,
    };
    // alert(JSON.stringify(data));
    try {
        const response = await fetch("/title/upd", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });

        const responseData = await response.json();
        if (!response.ok) {
            throw new Error(responseData.Error);
        }

        window.location.reload();
    } catch (error) {
        ShowMsg("Upd Titles failed: " + error.message, "error");
        window.location.reload();
    }
}

async function subAddUpd(mnu_id) {
    
    var txt = document.getElementById("_txt_"+mnu_id).value;

    if (txt == "") {
        alert("Nothing to add or update"+mnu_id);
        return; 
    }

    txt = txt.trim();

    var sub_uuid = parseInt(document.getElementById("_sub_"+mnu_id).value);
    
    if (sub_uuid == "") {
        id = 0;
    }

    // remove first three chars
    _mnuid = mnu_id.toString().substring(3);

    var data = {
        mnu_id: parseInt(_mnuid),
        sub_uuid: parseInt(sub_uuid),
        val: txt,
    };

    var url = "/"+mnu_TYPE+"/addupd";
    // alert(JSON.stringify(data)+" "+url);
    try {
        const response = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });

        const responseData = await response.json();
        if (!response.ok) {
            throw new Error(responseData.Error);
        }
        
        window.location.reload();
    } catch (error) {
        alert("Add/Update "+mnu_TYPE+" failed: " + error.message);
    }
}

async function lstDel(mnu_id) {

    var txt = document.getElementById("_txtmnu"+mnu_ID).value;
    var sub_uuid = parseInt(document.getElementById("_idsub"+mnu_ID).value);

    if (txt == "") {
        alert("Nothing to delete");
        return; 
    }

    if (!confirm("Delete '"+txt+"'?")) {
        return;
    }

    var data = {
        "sub_uuid": sub_uuid,
        "url": window.location.pathname,
    };

    try {
        const response = await fetch("/"+mnu_TYPE+"/delete", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });

        const responseData = await response.json(); 
       if (response.status != 200) {
            throw new Error(responseData.error);
        }

        window.location.reload();
    } catch (error) {
        alert("Delete "+mnu_TYPE+" failed: " + error.message);
        window.location.reload();
    }
}