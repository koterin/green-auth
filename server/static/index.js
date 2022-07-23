'use strict';

const emailField = document.getElementById('emailField');
const proceedBtn = document.getElementById('proceedBtn');
const emailValidateMsg = document.getElementById('emailValidateMsg');
const qrCode = document.getElementById('qr');

const hostUrl = 'https://password.berizaryad.ru';
const checkEmailUrlProd = hostUrl + '/api/checkEmail';
const sendCodeUrlProd = hostUrl + '/api/sendCode';
const TelegramBotUrl = 'https://t.me/berizaryad_password_manager_bot';

// Autocompletion - last used email is being filled in 
window.onload = () => {
    emailField.value = localStorage.getItem('email'); 

    if (emailField.checkValidity()) {
        proceedBtn.disabled = false;
    } else {
        proceedBtn.disabled = true;
    }

    const urlParams = new URLSearchParams(window.location.search);
    const redirectQuery = urlParams.get('redirect');
    if (redirectQuery != null) {
        sessionStorage.setItem('redirect', redirectQuery);
    }
}

proceedBtn.addEventListener('click', checkEmail);
qrCode.addEventListener('click', openTG);

// Disable button if email is not valid
emailField.addEventListener('keyup', function (event) {
    if (emailField.checkValidity()) {
        proceedBtn.disabled = false;
    } else {
        proceedBtn.disabled = true;
    }
    clearErrorEmailMsg();
});

function errorEmailMsg() {
    emailValidateMsg.innerHTML = "а ты зарегистрирован? если что - qr-код ниже";
}

function errorAccessMsg(resource) {
    emailValidateMsg.innerHTML = "у тебя нет доступа к " + resource;
}

function errorSendCodeMsg() {
    emailValidateMsg.innerHTML = "при отправке кода что-то пошло не так. давай еще";
}

function errorSendCodeLimitMsg() {
    emailValidateMsg.innerHTML = "получить код больше 5 раз в 5 минут нельзя.\nподождем";
}

function clearErrorEmailMsg() {
    emailValidateMsg.innerHTML = "";
}

function errorReloadPageMsg() {
    emailValidateMsg.innerHTML = "ой-ёй, что-то не то. надо обновить страницу";
}

function openTG() {
    location.href = TelegramBotUrl;
}

function deleteRedirectFromStorage() {
    sessionStorage.removeItem('redirect');
}

// POST /api/checkEmail
async function checkEmail() {
    var userEmail = emailField.value;
    var requestBody = {
        email : userEmail
    }
    var jsonBody = JSON.stringify(requestBody);

    fetch(checkEmailUrlProd, {
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
                localStorage.setItem('email', userEmail);
                location.href = "authenticate.html";
                return;
            } else if (response.status == '400') {
                console.log("Go away, Stranger");
                proceedBtn.disabled = true;
                errorEmailMsg();
                return Promise.reject(response.status);
            } else if (response.status == '401') {
                proceedBtn.disabled = true;
                errorAccessMsg();
                deleteRedirectFromStorage();
            } else {
                console.log("Internal Error");
                errorReloadPageMsg();
                return Promise.reject(response.status);
            };
        });
};
