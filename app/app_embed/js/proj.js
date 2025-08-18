
async function projAddUpd(_uuid) {

    _uuid = parseInt(_uuid);
    var _projid = 0;
    var _name = "";
    var _staid = 0;
    var _typid = 0;

    if (_uuid > 0) {
        _projid = parseInt(document.getElementById('_projid').value);
        _name = document.getElementById('_name').value;
        _staid = document.getElementById('_status').value;
        _typid = document.getElementById('_type').value;
    } else {
        _name = "New Project";
        _staid = 0;
        _typid = 0;
    }

    const data = {
        projid: _projid,
        uuid: _uuid,
        name: String(_name).trim(),
        staid: parseInt(_staid),
        typid: parseInt(_typid),
    };

    const options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    };
    //alert(JSON.stringify(data));
    try {
        const response = await fetch(`/proj/addupd`, options);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        if (result.message === 'success') {
            window.location.href = result.redirect;
        } else {
            ShowMsg('error',result.message);
        }
    } catch (error) {
        console.error('Error:', error);
        ShowMsg('error','An error occurred while updating the project. Please try again.');
    }
}

async function projDel(_uuid) {

    // Confirm deletion
    if (!confirm("Are you sure you want to delete this project?")) {
        return;
    }

    _uuid = parseInt(_uuid);

    try {
        const response = await fetch(`/proj/${_uuid}`, {
                                    method: 'DELETE',
                                    headers: {'Content-Type': 'application/json'}
                                });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        if (result.message === 'success') {
            // Redirect to the projects list or show a success message
            window.location.href = result.redirect;
        } else {
            ShowMsg('error', result.message);
        }
    } catch (error) {
        console.error('Error:', error);
        ShowMsg('error', 'An error occurred while deleting the project. Please try again.');
    }
}
