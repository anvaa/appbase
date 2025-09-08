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

function Start() {
    navigateTo("/app");
    // fetchPage("/app");
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
    //navigateTo("/info");
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

async function fetchPopup(url) {
    try {
        const response = await fetch(url);
        if (!response.ok) throw new Error('Network response was not ok');
        const html = await response.text();

        // Create popup container
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

        // Add content to popup
        popup.innerHTML = html;

        // Add close button
        const closeButton = document.createElement("button");
        closeButton.innerText = "X";
        Object.assign(closeButton.style, {
            position: "absolute",
            top: "1px",
            right: "1px"
        });
        closeButton.onclick = () => document.body.removeChild(popup);
        popup.appendChild(closeButton);

        // Append popup to body
        document.body.appendChild(popup);

        // Call global js() if defined
        if (typeof window.js === "function") {
            window.js();
        }
    } catch (error) {
        console.error('Error loading popup:', error);
        ShowMsg("Error loading popup: " + error.message, "error");
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
