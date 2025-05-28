var is_edit = false;
var sel_field = "";

function resetPage() {
    window.location.reload();
    return;
};

document.addEventListener("keydown", function (event) {
    
    if (event.key == "Enter") {
        if (sel_field = "uuid") {
            var uuid = document.getElementById("_uuid").value;
            if (uuid != "") {
                // remove all spaces
                uuid = uuid.replace(/\s/g, '');
                // check its numbers
                var regex = /^[0-9]+$/;
                if (!regex.test(uuid)) {
                    Msg("ID must be a number");
                    return;
                }

                QuickSearch('uuid',uuid);
            }
        }
        if (sel_field = "menu8") {
            var menu8 = document.getElementById("_menu8").value;
            if (menu8 != "") {
                QuickSearch('menu8',menu8);
            }
        }

        if (sel_field = "scan") {
            var scan = document.getElementById("_scan").value;
            if (scan != "") {
                // replace spaces with 
                scan = scan.replace(/\s/g, '@')
                QuickSearch('scan',scan);
            }
        }
    }
});

function isEdit(val) {
    is_edit = true;
    sel_field = val;
    emtyField(sel_field);
}

function setFocus(val) {
    document.getElementById(val).focus();
}

function emtyField() {
    if (sel_field == "menu8") {
        document.getElementById("_uuid").value = "";
        document.getElementById("_scan").value = "";
    }
    if (sel_field == "uuid") {
        document.getElementById("_menu8").value = "";
        document.getElementById("_scan").value = "";
    }
    if (sel_field == "scan") {
        document.getElementById("_uuid").value = "";
        document.getElementById("_menu8").value = "";
    }
}

async function QuickSearch(act,val) {

    var url = `/search/${act}/${val}`;
    var response = await fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    });
    // if status = 302 redirect to url
    if (response.status == 302) {
        var data = await response.json();
        if (data.redirect) {
            window.location.href = data.redirect;
        }
    } else if (response.status == 404) {
        var data = await response.json();
        if (data.msg) {
            Msg(data.msg);
        }
    } else {
        console.error('Error:', response.statusText);
    }

}

function Msg(val) {
    document.getElementById("_msg").innerHTML = val;
    document.getElementById("_msg").style.display = "block";
    document.getElementById("_msg").style.color = "red";
    document.getElementById("_msg").style.fontSize = "16px";
    document.getElementById("_msg").style.width = "fit-content";

    resetFields();

    setTimeout(function() {
        document.getElementById("_msg").style.display = "none";
    }, 3000);
    return;
}

function resetFields() {
    document.getElementById("_uuid").value = "";
    document.getElementById("_menu8").value = "";
    document.getElementById("_scan").value = "";
}