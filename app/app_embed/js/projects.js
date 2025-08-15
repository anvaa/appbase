
async function projUpd(_uuid) {

    const _name = document.getElementById('_name').value;
    const _staid = document.getElementById('_status').value;
    const _typid = document.getElementById('_type').value;
    
    const data = {
        uuid: parseInt(_uuid),
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