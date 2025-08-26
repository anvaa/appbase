
function noteContent(content, uuid) {
    // Implement your content display logic here
    const newText = content.replace(/\n/g, '<br>');
    document.getElementById('_content-' + uuid).innerHTML = newText;
};

async function noteAddUpd(_uuid) {
    
    const _type = document.getElementById('_type').value;
    const _content = document.getElementById('_content').value;
    const _header = document.getElementById('_header').value;
    const _projid = document.getElementById('_projid').value;
    const _persid = document.getElementById('_persid').value;

    if (_uuid > 0) {
        if (!_content.trim() || !_header.trim()) {
            ShowMsg('error', 'Note is empty');
            return;
        }
    }

    try {
        const response = await fetch(`/note/addupd`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                uuid: parseInt(_uuid, 10),
                projid: parseInt(_projid, 10),
                persid: parseInt(_persid, 10),
                type: parseInt(_type, 10),
                content: _content.trim(),
                header: _header.trim()
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result = await response.json();
        if (result.message === 'success') {
            // Redirect to the projects list or show a success message
            //window.location.reload();
            treeData(result.redirect);
        } else {
            ShowMsg('error', result.message);
        }
    }
    catch (error) {
        console.error('Error:', error);
        ShowMsg('error', 'An error occurred while saving the Note. Please try again.');
    }
}

async function noteDel(_uuid) {
    if (!confirm("Are you sure you want to delete this note?")) {
        return;
    }

    _uuid = parseInt(_uuid);

    try {
        const response = await fetch(`/note/${_uuid}`, {
                                    method: 'DELETE',
                                    headers: {'Content-Type': 'application/json'}
                                });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        if (result.message === 'success') {
            // Redirect to the projects list or show a success message
            treeData(result.redirect);
        } else {
            ShowMsg('error', result.message);
        }
    } catch (error) {
        console.error('Error:', error);
        ShowMsg('error', 'An error occurred while deleting the Note. Please try again.');
    }
}

async function NoteSort(val) {
    fetch(`/note/sort/${val}`, {
        method: 'POST',
    })
}
