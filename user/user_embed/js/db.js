
async function saveDbConfig() {
    const _sqlitePath = document.getElementById("_dbpath").value;
    const _dbType = document.getElementById("_dbtype").value;
    const _dbHost = document.getElementById("_dbhost").value;
    const _dbName = document.getElementById("_dbname").value;
    const _dbPort = document.getElementById("_dbport").value;
    const _dbUser = document.getElementById("_dbuser").value;
    const _dbPass = document.getElementById("_dbpass").value;

    var body = {
        type: _dbType,
        path: _sqlitePath,
        host: _dbHost,
        name: _dbName,
        port: _dbPort,
        user: _dbUser,
        password: _dbPass,
    }
    
    
    const response = await fetch("/v/dbconf", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(body)
    });
    const result = await response.json();
    if (result.success) {
        ShowMsg('info', result.message);
    } else {
        ShowMsg('error', result.message);
    }
}
