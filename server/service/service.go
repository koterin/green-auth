package service

import (
        "net/http"
        "encoding/json"
        "io/ioutil"
        "os"
        "path/filepath"
	"strings"

        "ktrn.com/dbhandler"
)

var TELEGRAM_URL = os.Getenv("TG_API_URL")
var TELEGRAM_BOT_KEY = os.Getenv("TG_BOT_KEY")
var CODE_LENGTH = 7
var UnathorizedPages [9]string = [9]string{"/main.css", "/index.js", "/index.html",
                                  "/", "/authenticate.html", "/auth.js", "/theme.js",
                                  "/images/tg_qr.png", "/favicon.ico"}
var Loginpage = "https://password.berizaryad.ru/"
var Homepage = "https://password.berizaryad.ru/home.html"
var PASSWORD_LENGTH = 8
var Filename = "servicepswd/.service"
var INTERNAL_ERROR_MSG = "Internal Error"
var OPTIONS = "OPTIONS"

type ResponseData struct {
    Email    string  `json:"email"`
    Code     string  `json:"code"`
    File     string  `json:"file"`
    Login    string  `json:"login"`
    Quantity string  `json:"quantity"`
}

type Answer struct {
    Status      int     `json:"status"`
    Response    string  `json:"response"`
    Login       string  `json:"login"`
    Password    string  `json:"password"`
}

func CorsHandler(w http.ResponseWriter, req *http.Request) {
    AddBasicHeaders(w)
    w.WriteHeader(http.StatusOK)

    answer := Answer {
                        Status: 200,
                        Response: "CORS, I see you mthfckr!1!",
                    }
    json.NewEncoder(w).Encode(answer)

    return
}

func ReadJson(w http.ResponseWriter, req *http.Request, respdata *ResponseData) (error) {
    resp, err := ioutil.ReadAll(req.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)

        return err
    }

    err = json.Unmarshal([]byte(resp), respdata)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)

        return err
    }

    return err
}

func AddBasicHeaders(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "http://test.password.ru")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    w.Header().Set("Content-Type", "application/json")
}

func CheckAdminRole(sessionId string) (int) {
    res, err := dbhandler.CheckAdminSession(sessionId);
    if ((err != nil) || (res == 0)) {
        return 403
    }

    return 200
}

type AuthFileServer struct {
    Path string
}

func AuthFileServerHandler(path string) http.Handler {
    return AuthFileServer{ path }
}

func (afs AuthFileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // Autoredirect to Homepage if user is logged in
    if req.URL.Path == "/" {
        sessionCookie, err := req.Cookie("sessionId")
        if err == nil {
	    dbErr := dbhandler.CheckSession(sessionCookie.Value)
            if (dbErr == nil) {
                http.Redirect(w, req, Homepage, http.StatusSeeOther)

		return
            }
        }
    }

    // Free access to log in pages for everyone
    for _, uri := range UnathorizedPages {
        if req.URL.Path == uri {
            path := filepath.Join(afs.Path, req.URL.Path)
            http.ServeFile(w, req, path)

            return
        }
    }

    // If user is not logged in, he is being redirected to the log in page
    sessionCookie, err := req.Cookie("sessionId")
    if err != nil {
        http.Redirect(w, req, Loginpage, http.StatusSeeOther)

	return
    }

    err = dbhandler.CheckSession(sessionCookie.Value)
    if err != nil {
        http.Redirect(w, req, Loginpage, http.StatusSeeOther)

	return
    }

    // If user is asking for swagger, he must have dev or admin role
    if (strings.Contains(req.URL.Path, "/swagger")) {
        sessionCookie, err := req.Cookie("sessionId")
        if err == nil {
	    dbErr := dbhandler.CheckDevSession(sessionCookie.Value)
	    if (dbErr != 1) {
                http.Redirect(w, req, Homepage, http.StatusSeeOther)

		return
            }
        }
    }

    path := filepath.Join(afs.Path, req.URL.Path)
    http.ServeFile(w, req, path)
}
