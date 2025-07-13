
async function projUpd() {
    const _uuid = document.getElementById('_uuid').value;
    const _name = document.getElementById('_name').value;
    const _note = document.getElementById('_note').value;
    const _stasub_id = document.getElementById('_status').value;
    const _typsub_id = document.getElementById('_type').value;

    const _csrf = document.getElementById('_csrf').value;
    const url = `/app/proj/addupd`;

    const data = {
        uuid: _uuid,
        name: _name,
        note: _note,
        stasub_id: _stasub_id,
        typsub_id: _typsub_id,
        csrf: _csrf
    };
}