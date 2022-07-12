'use strict';

const loginField = document.getElementById('loginField');
const proceedBtn = document.getElementById('proceedBtn');
const loginValidateMsg = document.getElementById('loginValidateMsg');
const qrCode = document.getElementById('qr');

const checkLoginUrlProd = 'https://green.auth.ktrn.com/api/checkLogin';  // Change to your URL
const sendCodeUrlProd = 'https://green.auth.ktrn.com/api/sendCode';  // Change to your URL
const TelegramBotUrl = 'https://t.me/berizaryad_green_auth_bot';  // Change to your Bot API URL

// Autocompletion - last used login is being filled in 
window.onload = () => {
    loginField.value = localStorage.getItem('login'); 

    if (loginField.checkValidity()) {
        proceedBtn.disabled = false;
    } else {
        proceedBtn.disabled = true;
    }
}

proceedBtn.addEventListener('click', checkLogin);
qrCode.addEventListener('click', openTG);

// Disable button if login is not valid
loginField.addEventListener('keyup', function (event) {
    if (loginField.checkValidity()) {
        proceedBtn.disabled = false;
    } else {
        proceedBtn.disabled = true;
    }
    clearErrorLoginMsg();
});

function errorLoginMsg() {
    loginValidateMsg.innerHTML = "incorrect login";
}

function errorSendCodeMsg() {
    loginValidateMsg.innerHTML = "something went wrong while sending code. c`mon again";
}

function errorSendCodeLimitMsg() {
    loginValidateMsg.innerHTML = "you can't get more than 5 codes in 5 minutes.\nlet's wait";
}

function clearErrorLoginMsg() {
    loginValidateMsg.innerHTML = "";
}

function errorReloadPageMsg() {
    loginValidateMsg.innerHTML = "yay, something's wrong. reload the page";
}

function openTG() {
    location.href = TelegramBotUrl;
}

// POST /api/checkLogin
async function checkLogin() {
    var userLogin = loginField.value;
    var requestBody = {
        login : userLogin
    }
    var jsonBody = JSON.stringify(requestBody);

    fetch(checkLoginUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
        withCredentials: true
    })
       .then(function (response) {
            if (response.status == '200') {
                console.log("Hey, I know you!");
                localStorage.setItem('login', userLogin);
                location.href = "authenticate.html";
                return;
            } else if (response.status == '400') {
                console.log("Go away, Stranger");
                proceedBtn.disabled = true;
                errorLoginMsg();
                return Promise.reject(response.status);
            } else {
                console.log("Internal Error");
                errorReloadPageMsg();
                return Promise.reject(response.status);
            };
        });
};
