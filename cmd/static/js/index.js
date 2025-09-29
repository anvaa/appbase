document.addEventListener("keydown", (event) => {
  if (event.key !== "Enter") return;
  const pageId = getInputValue("_pageid");
  if (pageId === "login") return Login();
  if (pageId === "signup") return Signup();
});

async function Login() {
  const username = getInputValue("_username");
  const password = getInputValue("_password");

  if (!isValidUsername(username)) return showMsg("Invalid username");
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
  const password = getInputValue("_password");
  const password2 = getInputValue("_password2");
  const orgname = getInputValue("_orgname");
  const count = parseInt(getInputValue("_count"));

  if (!validatePasswords(password, password2)) return;

  try {
    const response = await sendRequest("/signup", { username, password, password2, orgname, count });
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

function isValidUsername(username) {
  if (username.length < 6 || username.length > 50) return showMsg("Username must be between 6 and 50 characters"), false;
  if (containsSqlInjection(username)) return showMsg("Invalid Username"), false;
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(username);
}

function isValidPassword(password) {
  return password.length >= 8 && password.length <= 50;
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
  // if response == 200, redirect to the URL provided in response.url
  //alert(JSON.stringify(response));
  if (response && response.url) {
    // alert("Redirecting to: " + response.url);
    window.location.href = response.url;
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
