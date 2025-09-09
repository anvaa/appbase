document.addEventListener("keydown", handleKeyDown);

function handleKeyDown(event) {
    if (!event.ctrlKey) return;
    switch (event.key) {
        case "1":
            Start();
            break;
        case "q":
            Logout();
            break;
    }
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
    fetchPopup("/info");
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

function navigateTo(url) {
    window.location.href = url;
}

async function fetchPage(url) {
    try {
        loadScriptForUrl(url);
        const html = await fetchHtml(url);
        replaceBodyWithContent(html);
        callGlobalJs();
    } catch (error) {
        handleError("Error loading page: " + error.message);
    }
}

function loadScriptForUrl(url) {
    const jsfile = url.split('/')[1];
    if (jsfile) {
        const script = document.createElement('script');
        script.src = `/js/${jsfile}.js`;
        script.async = true;
        document.head.appendChild(script);
    }
}

async function fetchHtml(url) {
    const response = await fetch(url);
    if (!response.ok) throw new Error('Network response was not ok');
    return await response.text();
}

function replaceBodyWithContent(html) {
    document.body.innerHTML = "";
    const content = document.createElement("div");
    content.id = "content";
    content.innerHTML = html;
    document.body.appendChild(content);
}

function callGlobalJs() {
    if (typeof window.js === "function") {
        window.js();
    }
}

async function fetchPopup(url) {
    try {
        const html = await fetchHtml(url);
        showPopup(html);
        callGlobalJs();
    } catch (error) {
        handleError("Error loading popup: " + error.message);
    }
}

function showPopup(html) {
    const popup = document.createElement("div");
    popup.id = "popup";
    Object.assign(popup.style, {
        position: "fixed",
        top: "50%",
        left: "50%",
        transform: "translate(-50%, -50%)",
        backgroundColor: "white",
        border: "2px solid black",
        padding: "20px",
        zIndex: "1000",
        boxShadow: "0 4px 8px rgba(119, 117, 117, 1)"
    });

    popup.innerHTML = html;
    const closeButton = createCloseButton(() => document.body.removeChild(popup));
    popup.appendChild(closeButton);
    document.body.appendChild(popup);
}

function createCloseButton(onClick) {
    const btn = document.createElement("button");
    btn.innerText = "X";
    Object.assign(btn.style, {
        position: "absolute",
        top: "1px",
        right: "1px"
    });
    btn.onclick = onClick;
    return btn;
}

function handleError(msg) {
    console.error(msg);
    ShowMsg("error", msg);
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
