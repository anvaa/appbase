
var is_edit = "";
var sel_Txt = "";
var mnu_ID = 0;
var sub_UUID = 0;


function resetPage() {
    window.location.reload();
    return;
};

document.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
        
        if (is_edit == "") {
            return;
        }
        
        if (is_edit == "titles") {
            TitlesUpd(mnu_ID);
            return;
        } else {
            subAddUpd(mnu_ID);
            return;
        }
    }
});

function isEdit(val,mnuid) {
    is_edit = val;
    mnu_ID = mnuid;
    
    mnu_ID = document.getElementById("_idmnu"+mnu_ID).value;
    sel_Txt = document.getElementById("_txtmnu"+mnu_ID).value;
    //alert("isEdit: "+is_edit+" mnu_ID: "+mnu_ID+" sel_Txt: "+sel_Txt);
}

function printTxt() {
    const is_printtxt = document.getElementById("_printtxt").checked;
    if (is_printtxt) {
        document.getElementById("_printtxt").value = "0";
        document.getElementById("_printtxt").checked = false;
    } else {
        document.getElementById("_printtxt").value = "1";
        document.getElementById("_printtxt").checked = true;
    }
    printConf();
    return;
}

function startPageFocus() {
    // get radio button value
    var radioValue = document.querySelector('input[name="foc_select"]:checked').value;

    var data = {
        "start_page_focus": radioValue,
    };
    // alert(JSON.stringify(data));
    try {
        const response = fetch("/tools/appconf", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });

        const responseData = response.json();
        if (!response.ok) {
            alert(responseData.Error);
            throw new Error(responseData.Error);
        }
        
    } catch (error) {
        // ShowMsg("Start Page Focus failed: " + error.message);
    }

    return;
}

function printConf() {

    var _printtxt = document.getElementById("_printtxt").value;
    var _height = document.getElementById("_print_height").value;
    var _width = document.getElementById("_print_width").value;
    var _margin = document.getElementById("_print_margin").value;
    var _fontSize = document.getElementById("_print_fontsize").value;

    var data = {
        print_txt : parseInt(_printtxt),
        font_size : parseInt(_fontSize),
        height : parseInt(_height),
        width : parseInt(_width),
        margin : parseInt(_margin),
    };
    // alert(JSON.stringify(data));
    try {
        const response = fetch("/tools/printconf", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });

        const responseData = response.json();
        if (!response.ok) {
            alert(responseData.Error);
            throw new Error(responseData.Error);
        }
        
        document.getElementById("_printtxt").checked = responseData.data.print_txt;
        document.getElementById("_print_height").style.height = responseData.height;
        document.getElementById("_print_width").style.width = responseData.width;
        document.getElementById("_print_margin").value = responseData.margin;
        document.getElementById("_print_fontsize").value = responseData.font_size;
    
    } catch (error) {
        // ShowMsg("Print Config failed: " + error.message);
    }
}

// MENUS //
function subSel(sub_uuid, mnuid) {
    mnu_ID = mnuid;
    sub_UUID = sub_uuid;
    is_edit = "subitem";

    if (mnu_ID === 100) {
        is_edit = "status";
    } 
    if (mnu_ID === 200) {
        is_edit = "titles";
    }

    sel_Txt = document.getElementById("_selsub"+sub_UUID).value;
    document.getElementById("_txtmnu"+mnu_ID).value = sel_Txt;
    document.getElementById("_idsub"+mnu_ID).value = sub_UUID;
    mnu_ID = document.getElementById("_idmnu"+mnuid).value;
}

async function TitlesUpd() {
    var _txt = document.getElementById("_txtmnu200").value;
    var _uuid = parseInt(document.getElementById("_idsub200").value);

    if (_uuid == 0 || _txt == "") {
        return; 
    }
    
    var data = {
        "mnu_id": 200,
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
        alert("Upd Titles "+txt+" failed: " + error.message);
        window.location.reload();
    }
}

async function subAddUpd(mnu_id) {
    
    var txt = document.getElementById("_txtmnu"+mnu_ID).value;

    if (txt == "") {
        alert("Status: Nothing to add or update");
        return; 
    }

    var mnu_id = parseInt(document.getElementById("_idmnu"+mnu_ID).value);
    var sub_uuid = parseInt(document.getElementById("_idsub"+mnu_ID).value);
    
    if (sub_uuid == "") {
        id = 0;
    }

    var data = {
        "txt": txt,
        "mnu_id": mnu_id,
        "sub_uuid": sub_uuid,
    };
    // alert(JSON.stringify(data));
    try {
        const response = await fetch("/"+is_edit+"/addupd", {
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
        alert("Add/Update "+is_edit+" failed: " + error.message);
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
        const response = await fetch("/"+is_edit+"/delete", {
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
        alert("Delete "+is_edit+" failed: " + error.message);
        window.location.reload();
    }
}