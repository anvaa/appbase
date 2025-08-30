document.addEventListener("keydown", handleKeyDown);

function handleKeyDown(event) {
  if (event.key === "Enter") {
    const pageId = document.getElementById("_pageid").value;

    if (pageId === "login") {
      Login();
    } else if (pageId === "signup") {
      Signup();
    }
  }
}

async function Login() {
  const email = getInputValue("_email");
  const password = getInputValue("_password");

  if (!validEmail(email)) return ShowMsg("Invalid email");
  if (!validPassword(password)) return ShowMsg("Invalid password");

  const data = { email, password };

  try {
    const response = await sendRequest("/login", data);
    handleResponse(response, "Login");
  } catch (error) {
    ShowMsg(`Login: ${error.message}`);
  }
}

async function Signup() {
  const email = getInputValue("_email");
  const password = getInputValue("_password");
  const password2 = getInputValue("_password2");
  const orgname = getInputValue("_orgname");
  const count = parseInt(getInputValue("_count"));


  if (!validatePasswords(password, password2)) return;

  const data = { email, password, password2, orgname, count };

  try {
    const response = await sendRequest("/signup", data);
    handleResponse(response, "Signup");
  } catch (error) {
    ShowMsg(`Signup: ${error.message}`);
  }
}

function validatePasswords(psw1, psw2) {
  if (psw1 !== psw2) return ShowMsg("Passwords do not match"), false;
  if (psw1.length < 8 || psw1.length > 50) {
    return ShowMsg("Passwords must be at least 8 characters long"), false;
  }
  if (noSqlInText(psw1)) {
    return ShowMsg("Invalid Password"), false;
  }
  if (noSqlInText(psw2)) {
    return ShowMsg("Invalid Password"), false;
  }
  return true;
}

function validEmail(email) {
  if (email.length < 5 || email.length > 50) {
    ShowMsg("Email must be between 5 and 50 characters");
    return false;
  }
  if (noSqlInText(email)) {
    ShowMsg("Invalid Email");
    return false;
  }
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(email);
}

function validPassword(password) {
  return password.length >= 8 && password.length <= 50;
  
}

function noSqlInText(text) {
  const noSqlPatterns = [
    /['"]/g, // Single or double quotes
    /[;:]/g, // Semicolon or colon
    /\b(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|TRUNCATE|RENAME|EXEC)\b/gi, // SQL keywords
    /\b(WHERE|FROM|JOIN|ON|INTO|VALUES|SET)\b/gi, // SQL clauses
    /\b(AND|OR|NOT)\b/gi, // Logical operators
  ];

  return noSqlPatterns.some((pattern) => pattern.test(text));
}

function getInputValue(id) {
  return document.getElementById(id).value;
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
  if (response.url) {
    window.location.href = response.url;
  } else {
    throw new Error(`${action} failed`);
  }
}

function Cancel(val) {
  window.location.href = val;
}

function ShowMsg(_msg) {
  const messageElement = document.getElementById("_msg");
  messageElement.innerHTML = _msg;
  messageElement.style.borderColor = "red";
  messageElement.style.width = "inherit";
  messageElement.style.padding = "5px";
  messageElement.style.textAlign = "center";

  setTimeout(() => {
    messageElement.innerHTML = "";
    messageElement.style.borderColor = "transparent";
    messageElement.style.width = "0px";
    messageElement.style.padding = "0px";
  }
  , 5000);
}