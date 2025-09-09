async function saveDbConfig() {
    const fields = [
        { id: "_dbpath", key: "path" },
        { id: "_dbtype", key: "type" },
        { id: "_dbhost", key: "host" },
        { id: "_dbname", key: "name" },
        { id: "_dbport", key: "port" },
        { id: "_dbuser", key: "user" },
        { id: "_dbpass", key: "password" }
    ];

    const body = fields.reduce((acc, field) => {
        acc[field.key] = document.getElementById(field.id).value;
        return acc;
    }, {});

    try {
        const response = await fetch("/v/dbconf", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body)
        });
        const result = await response.json();
        if (result.message === "success") {
            ShowMsg('info', result.msg);
        } else {
            ShowMsg('error', result.error);
        }
    } catch (error) {
        ShowMsg('error', error.message || "Network error");
    }
}
