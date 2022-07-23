'use strict';

const textContent = document.getElementById('textContent');
const getFileBtn = document.getElementById('getFileBtn');
const addPassBtn = document.getElementById('addPassBtn');
const infoField = document.getElementById('infoField');
const infoField2 = document.getElementById('infoField2');
const genPassBtn = document.getElementById('genPassBtn');
const multiAddPassBtn = document.getElementById('multiAddPassBtn');
const multiGenPassBtn = document.getElementById('multiGenPassBtn');

const loader = document.getElementById('loader');
const loginInputMsg = document.getElementById('loginInputMsg');
const loginInput = document.getElementById('loginInput');
const loginValidateMsg = document.getElementById('loginValidateMsg');

const passModal = document.getElementById('passModal');
const closeModal = document.getElementById('closeModal');
const pass = document.getElementById('pass');
const passes = document.getElementById('passes');
const modalMsg = document.getElementById('modalMsg');
const clipPassBtn = document.getElementById('clipPassBtn');
const clipPassesBtn = document.getElementById('clipPassesBtn');

const hostUrl = 'https://password.berizaryad.ru';
const getFileUrlProd = hostUrl + '/api/getFile';
const addPassUrlProd = hostUrl + '/api/addPass';
const genPassUrlProd = hostUrl + '/api/generatePass';
const multiGenPassUrlProd = hostUrl + '/api/multiGeneratePass';

getFileBtn.addEventListener('click', preGetFile);
addPassBtn.addEventListener('click', preAddPass);
multiAddPassBtn.addEventListener('click', preMultiAddPass);
genPassBtn.addEventListener('click', generatePass);
multiGenPassBtn.addEventListener('click', multiGeneratePass);

clipPassesBtn.addEventListener('click', clipPasses);
clipPassBtn.addEventListener('click', clipPass);
backBtn.addEventListener('click', pageBack);

window.onload = () => {
    textContent.style.display = "none";
    loginInput.style.display = "none";
    loginInputMsg.style.display = "none";
    loginValidateMsg.style.display = "none";
    genPassBtn.style.display = "none";
    multiGenPassBtn.style.display = "none";
}

window.onclick = function(event) {
  if (event.target == passModal) {
    passModal.style.display = "none";
  }
}

closeModal.onclick = function() {
  passModal.style.display = "none";
}

loginInput.addEventListener('keyup', function (event) {
    loginValidateMsg.innerHTML = "";
});

function pageBack() {
    location.href = "index.html";
}

function errorReloadPageMsg(id) {
    id.innerHTML = "ой-ёй, что-то не то. надо обновить страницу";
}

function errorNoAdminAccessMsg(id) {
    id.innerHTML = "не трогай, это для админов";
}

function preGetFile() {
    loader.innerHTML = "Загрузка...";
    getFile();
}

function preAddPass() {
    loader.innerHTML = "Загрузка...";
    addPass();
}

function preMultiAddPass() {
    loader.innerHTML = "Загрузка...";
    multiAddPass();
}

function clipPasses() {
    passes.select();
    document.execCommand("copy");
    alert("Данные скопированы");
}

function clipPass() {
  var TempText = document.createElement("input");
  TempText.value = pass.innerHTML;
  document.body.appendChild(TempText);
  TempText.select();
  
  document.execCommand("copy");
  document.body.removeChild(TempText);
  alert("Данные скопированы");
}

async function getFile() {
    fetch(getFileUrlProd, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'X-Green-Origin': 'password.berizaryad.ru'
	},
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '400' || response.status == '500') {
            console.log("Error while getting file");
    	    loader.innerHTML = "ошибка";
            return Promise.reject(response);
	} else if (response.status == '200') {
    	    loader.innerHTML = "";
            return (response.json());
        } else {
    	    loader.innerHTML = "ошибка";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        textContent.innerHTML = json.response;
        textContent.style.display = "block";
        loginInputMsg.style.display = "none";
        loginValidateMsg.style.display = "none";
        loginInput.style.display = "none";
        multiGenPassBtn.style.display = "none";
        genPassBtn.style.display = "none";
    });
}

async function addPass() {
   loader.style.display = "block";

   fetch(addPassUrlProd, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'X-Green-Origin': 'password.berizaryad.ru'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '500') {
    	    loader.innerHTML = "";
            errorReloadPageMsg(infoField);
            return Promise.reject(response);
        } else if (response.status == '403') {
    	    loader.innerHTML = "";
            errorNoAdminAccessMsg(infoField);
            return Promise.reject(response);
        } else if (response.status == '200') {
    	    loader.innerHTML = "";
            return (response.json());
        } else {
    	    loader.innerHTML = "ошибка";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        multiGenPassBtn.style.display = "none";
        textContent.style.display = "none";
        loginInputMsg.innerHTML = "Введите новый логин";
        loginInputMsg.style.display = "block";
        loginInput.value = json.login;
        loginInput.style.display = "block";
        genPassBtn.style.display = "block";
    });
}

async function generatePass() {
    var requestBody = {
        login : loginInput.value
    }
    var jsonBody = JSON.stringify(requestBody);

    loginValidateMsg.style.display = "block";
    loginValidateMsg.innerHTML = "Загрузка...";

    fetch(genPassUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json',
            'X-Green-Origin': 'password.berizaryad.ru'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '204') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "пустые логины предлагаем? нехорошо это";
            return Promise.reject(response);
        } else if (response.status == '400') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "в логине не должно быть пробелов";
            return Promise.reject(response); 
        } else if (response.status == '409') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "такой логин уже есть. давай другой";
            return Promise.reject(response); 
        } else if (response.status == '200') {
            loginValidateMsg.innerHTML = "";
            return (response.json());
        } else {
            loginValidateMsg.style.display = "block";
	    loginValidateMsg.innerHTML = "ой-ёй, что-то не то. надо обновить страницу";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        pass.innerHTML = json.password;
        modalMsg.innerHTML = "Сохраните этот пароль! Он будет показан только один раз";
	passModal.style.display = "block";
	pass.style.display = "block";
	passes.style.display = "none";
	clipPassBtn.style.display = "inline-block";
	clipPassesBtn.style.display = "none";
    });
}

async function multiGeneratePass() {
    var requestBody = {
        quantity : loginInput.value
    }
    var jsonBody = JSON.stringify(requestBody);
        
    loginValidateMsg.style.display = "block";
    loginValidateMsg.innerHTML = "Загрузка...";

    fetch(multiGenPassUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json',
            'X-Green-Origin': 'password.berizaryad.ru'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '400') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "неверный формат количества";
            return Promise.reject(response); 
        } else if (response.status == '400') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "за раз можно сгенерировать от 2 до 100 записей";
            return Promise.reject(response); 
        } else if (response.status == '200') {
            loginValidateMsg.innerHTML = "";
            return (response.json());
        } else {
            loginValidateMsg.style.display = "block";
	    loginValidateMsg.innerHTML = "ой-ёй, что-то не то. надо обновить страницу";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        passes.innerHTML = json.response;
        modalMsg.innerHTML = "Сохраните этот список! Он будет показан только один раз";
        passModal.style.display = "block";
	passes.style.display = "block";
	pass.style.display = "none";
	clipPassBtn.style.display = "none";
	clipPassesBtn.style.display = "inline-block";
    });
}

async function multiAddPass() {
   loader.style.display = "block";

   fetch(addPassUrlProd, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'X-Green-Origin': 'password.berizaryad.ru'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '500') {
    	    loader.innerHTML = "";
            errorReloadPageMsg(infoField2);
            return Promise.reject(response);
        } else if (response.status == '403') {
    	    loader.innerHTML = "";
            errorNoAdminAccessMsg(infoField2);
            return Promise.reject(response);
        } else if (response.status == '200') {
    	    loader.innerHTML = "";
            return (response.json());
        } else {
    	    loader.innerHTML = "ошибка";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        textContent.style.display = "none";
        loginInputMsg.innerHTML = "Введите количество новых записей";
        loginInputMsg.style.display = "block";
        loginInput.value = "2";
        loginInput.style.display = "block";
        genPassBtn.style.display = "none";
        multiGenPassBtn.style.display = "block";
    });
}
