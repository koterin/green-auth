'use strict';

const previousLogin = document.getElementById('previousLogin');
const codeField = document.getElementById('codeField');
const sendBtn = document.getElementById('sendBtn');
const infoField = document.getElementById('infoField');
const timerTextField = document.getElementById('timerTextField');
const timerField = document.getElementById('timerField');
const timerLimit = 30;
const backBtn = document.getElementById('backBtn');

const sendCodeUrlProd = 'https://green.auth.ktrn.com/api/sendCode';
const validateCodeUrlProd = 'https://green.auth.ktrn.com/api/validateCode';

window.onload = () => {
    if (localStorage.getItem('login') == null) {
        location.href = "index.html";
    }

    sendCode();
}

document.addEventListener("DOMContentLoaded", startTimer);
previousLogin.innerHTML = localStorage.getItem('login');

backBtn.addEventListener('click', pageBack);
sendBtn.addEventListener('click', sendCode);

codeField.addEventListener('input', function (event) {
    infoField.innerHTML = "";
    
    if (codeField.checkValidity()) {
        validateCode();
        codeField.disabled = true;
    }
});

function errorSendCodeMsg() {
    infoField.innerHTML = "something went wrong while sending code. c`mon again";
}

function errorWrongCodeMsg() {
    infoField.innerHTML = "ouch. that's a wrong code";
}

function errorAttemptsLimitMsg() {
    infoField.innerHTML = "you gotta get a new code";
}

function errorSendCodeLimitMsg() {
    infoField.innerHTML = "you can't get more than 5 codes in 5 minutes.\nlet's wait";
}

function startTimer() {
    timerTextField.innerHTML = "You can ask for another code after: &nbsp";
    var seconds = countdownTimer(timerLimit);
    var interval = setInterval(() => {
        seconds = countdownTimer(seconds);
        if (seconds == 0) {
            clearInterval(interval);
            timerTextField.innerHTML = "Code is outdated";
            timerField.innerHTML = "";
            sendBtn.disabled = false;
            codeField.disabled = false;
        }
    }, 1000);
}

function countdownTimer(seconds) {
    timerField.innerHTML = seconds;
    return (seconds = seconds - 1);
}

function pageBack() {
    location.href = "index.html";
}

async function sendCode() {
    var requestBody = {
        login : localStorage.getItem('login')
    }
    var jsonBody = JSON.stringify(requestBody);

    startTimer();
    infoField.innerHTML = "";
    codeField.value = "";

    fetch(sendCodeUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: "same-origin",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '400' || response.status == '500') {
            console.log("Error while sending code");
            errorSendCodeMsg();
            return;
        } else if (response.status == '429') {
            console.log("SendCode limit error");
            errorSendCodeLimitMsg();
            sendBtn.disabled = true;
            return;
        } else if (response.status == '200') {
            console.log("ok");
            sendBtn.disabled = true;
            return;
        } else {
	        return Promise.reject(response);
        }
    });
}

async function validateCode() {
    var requestBody = {
        login : localStorage.getItem('login'),
        code : codeField.value
    }
    var jsonBody = JSON.stringify(requestBody);
    console.log(jsonBody);

    fetch(validateCodeUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        codeField.disabled = false;

        if (response.status == '500') {
            console.log("Error while sending code");
            errorSendCodeMsg();
            return;
        } else if (response.status == '400') {
            console.log("Wrong code");
            errorWrongCodeMsg();
            return;
        } else if (response.status == '429') {
            console.log("Too many attempts");
            errorAttemptsLimitMsg();
            return;
        } else if (response.status == '200') {
            console.log("Access granted");
            location.href = "home.html";
            return;
        } else {
	        return Promise.reject(response);
        }
    });
}
