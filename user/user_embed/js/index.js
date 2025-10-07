document.addEventListener("keydown", (event) => {
  if (event.key !== "Enter") return;
  const pageId = getInputValue("_pageid");
  if (pageId === "login") return Login();
  if (pageId === "signup") return Signup();
});

async function Login() {
  const username = getInputValue("_username");
  const password = getInputValue("_password");
  
  //if username contains '@', treat it as email
  if (username.includes("@") && username.includes(".")) {
    if (!isValidEmail(username)) return showMsg("Invalid email");
  }
  if (!isValidPassword(password)) return showMsg("Invalid password");
  
  try {
    const response = await sendRequest("/login", { username, password });
    handleResponse(response, "Login");
  } catch (error) {
    showMsg(`Login: ${error.message}`);
  }
}

async function Signup() {
  const username = getInputValue("_username");
  const email = getInputValue("_email");
  const password = getInputValue("_password");
  const password2 = getInputValue("_password2");
  const orgname = getInputValue("_orgname");
  const count = parseInt(getInputValue("_count"));

  if (!validatePasswords(password, password2)) return;

  try {
    const response = await sendRequest("/signup", { username, email, password, password2, orgname, count });
    handleResponse(response, "Signup");
  } catch (error) {
    showMsg(`Signup: ${error.message}`);
  }
}

function validatePasswords(psw1, psw2) {
  if (psw1 !== psw2) return showMsg("Passwords do not match"), false;
  if (!isValidPassword(psw1)) return showMsg("Passwords must be at least 8 characters long"), false;
  if (containsSqlInjection(psw1) || containsSqlInjection(psw2)) return showMsg("Invalid Password"), false;
  return true;
}

function isValidEmail(username) {
  if (username.length < 6 || username.length > 255) return showMsg("Username must be between 6 and 255 characters"), false;
  if (containsSqlInjection(username)) return showMsg("Invalid Email"), false;
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(username);
}

function isValidPassword(password) {
  return password.length >= 8 && password.length <= 255;
}

function containsSqlInjection(text) {
  const patterns = [
    /['"]/g,
    /[;:]/g,
    /\b(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|TRUNCATE|RENAME|EXEC)\b/gi,
    /\b(WHERE|FROM|JOIN|ON|INTO|VALUES|SET)\b/gi,
    /\b(AND|OR|NOT)\b/gi,
  ];
  return patterns.some((pattern) => pattern.test(text));
}

function getInputValue(id) {
  const el = document.getElementById(id);
  return el ? el.value : "";
}

async function sendRequest(url, data) {
  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  const responseData = await response.json();
  if (!response.ok) throw new Error(responseData.message || "Server error");
  return responseData;
}

function handleResponse(response, action) {
  // if response == 200, redirect to the URL provided in response.redirect
  //alert(JSON.stringify(response));
  if (response && response.redirect) {
    // alert("Redirecting to: " + response.redirect);
        window.location.href = response.redirect;
  } else {
    showMsg(`${action} failed: ${response.message || "Unknown error"}`);
  }
}

function Cancel(val) {
  window.location.href = val;
}

function showMsg(msg) {
  const el = document.getElementById("_msg");
  if (!el) return;
  el.innerHTML = msg;
  el.style.borderColor = "red";
  el.style.width = "inherit";
  el.style.padding = "5px";
  el.style.textAlign = "center";
  setTimeout(() => {
    el.innerHTML = "";
    el.style.borderColor = "transparent";
    el.style.width = "0px";
    el.style.padding = "0px";
  }, 5000);
}
