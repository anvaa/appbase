document.addEventListener("keydown", function (event) {
    if (event.ctrlKey) {
        switch (event.key) {
            case "1":
                Start();
                break;
            case "q":
                Logout();
                break;
        }
    }
});

function navigateTo(url) {
    window.location.href = url;
}

function Start() {
    navigateTo("/app");
}

function Tools() {
    navigateTo("/tools/titles");
}

function Users() {
    navigateTo("/v/users");
}

function MyAccount() {
    navigateTo("/v/myaccount");
}

function AppInfo() {
    navigateTo("/info");
}

function Logout() {
    navigateTo("/logout");
}

// Tools menus
function toolsTitles() {
    navigateTo("/tools/titles");
}

function toolsStatus() {
    navigateTo("/tools/status");
}

function toolsTypes() {
    navigateTo("/tools/types");
}

async function fetchPage(url) {
    try {
        // Dynamically load related JS file based on URL
        const urlParts = url.split('/');
        const jsfile = urlParts[1];
        if (jsfile) {
            const script = document.createElement('script');
            script.src = `/js/${jsfile}.js`;
            script.async = true;
            document.head.appendChild(script);
        }

        const response = await fetch(url);
        if (!response.ok) throw new Error('Network response was not ok');
        const html = await response.text();

        // Replace body content with new page
        document.body.innerHTML = "";

        const content = document.createElement("div");
        content.id = "content";
        content.innerHTML = html;
        document.body.appendChild(content);

        // Call global js() if defined
        if (typeof window.js === "function") {
            window.js();
        }
    } catch (error) {
        console.error('Error loading page:', error);
        ShowMsg("Error loading page: " + error.message, "error");
    }
}

function ShowMsg(type, val) {
    const msgEl = document.getElementById("_msg");
    if (!msgEl) return;

    document.body.style.pointerEvents = "none";

    msgEl.innerHTML = val;
    Object.assign(msgEl.style, {
        margin: "2px",
        fontSize: "large",
        fontWeight: "bold",
        fontFamily: "monospace",
        width: "fit-content",
        padding: "2px",
        borderRadius: "5px",
        border: `3px solid ${getBorderColor(type)}`,
        color: getBorderColor(type),
        display: "block",
        opacity: "1",
        transition: ""
    });

    setTimeout(() => {
        msgEl.style.transition = "opacity 0.5s";
        msgEl.style.opacity = "0";
        setTimeout(() => {
            msgEl.style.display = "none";
            msgEl.style.opacity = "1";
            document.body.style.pointerEvents = "auto";
        }, 500);
    }, 3000);
}

function getBorderColor(type) {
    switch (type) {
        case "error": return "red";
        case "success": return "green";
        case "info": return "blue";
        default: return "black";
    }
}
