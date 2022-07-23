package service

import (
        "net/http"
        "net/url"
        "encoding/json"
        "io/ioutil"
        "os"
        "path/filepath"
        "net/http/httputil"
	"log"

        "ktrn.com/dbhandler"
)

var TELEGRAM_URL = os.Getenv("TG_API_URL")
var TELEGRAM_BOT_KEY = os.Getenv("TG_BOT_KEY")
var API_KEY = os.Getenv("API_KEY")
var CODE_LENGTH = 7
var PASSWORD_LENGTH = 8

var UnathorizedPages [10]string = [10]string{"/main.css", "/index.js", "/index.html",
                                  "/", "/authenticate.html", "/auth.js", "/theme.js",
                                  "/images/tg_qr.png", "/favicon.ico", "/wrongredirect.html"}
var HOST_URL = os.Getenv("HOST_URL")
var Loginpage = HOST_URL
var Homepage = HOST_URL + "/home.html"
var WrongRedirectPage = HOST_URL + "/wrongredirect.html"

var Filename = "servicepswd/.service"
var INTERNAL_ERROR_MSG = "Internal Error"
var OPTIONS = "OPTIONS"

type ResponseData struct {
    Email    string  `json:"email"`
    Code     string  `json:"code"`
    File     string  `json:"file"`
    Login    string  `json:"login"`
    Quantity string  `json:"quantity"`
    Redirect string  `json:"redirect"`
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
    w.Header().Set("Access-Control-Allow-Origin", HOST_URL)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    w.Header().Set("Content-Type", "application/json")
}

func requestDebugger(req *http.Request) {
    log.Println("____DEBUGGING____")
    reqDump, _ := httputil.DumpRequest(req, true)
    log.Println(string(reqDump))
}

func CheckRole(req *http.Request, role string) (int) {
    sessionCookie, err := req.Cookie("sessionId")
    if err != nil {
        return 403
    }

    if !dbhandler.CheckRole(sessionCookie.Value, role) {
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

func GetCookie(req *http.Request) bool {
    sessionCookie, err := req.Cookie("sessionId")
    if err != nil {
        return false
    }

    origin := req.Header.Get("X-Green-Origin");
    if ((origin == "null") || (origin == "")) {
	aUrl, _ := url.Parse(HOST_URL)
	origin = aUrl.Host
    }

    if !dbhandler.CheckSession(sessionCookie.Value, origin) {
        return false
    }

    return true
}

func checkUrl(redirect string) uint {
    _, err := url.Parse(redirect)
    if err != nil {
        return 401
    }

    _, err = url.ParseRequestURI(redirect)
    if err != nil {
        return 401
    }

    return 1
}

func checkRedirect(w http.ResponseWriter, req *http.Request) uint {
    redirect := req.URL.Query().Get("redirect")
    if redirect != "" {
	return checkUrl(redirect)
    }

    return 0
}

func CheckApiKey(req *http.Request) bool {
    key := req.Header.Get("Api-Key")
    if key == API_KEY {
        return true
    }

    log.Println("Wrong API Key: ", key)

    return false
}

func (afs AuthFileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // Check if redirect query is present and correct
    redir := checkRedirect(w, req)
    if redir == 401 {
        http.Redirect(w, req, WrongRedirectPage, http.StatusSeeOther)
        return
    }

    // Autoredirect to Homepage if user is logged in
    // But not if he is being redirected from other resource
    access := GetCookie(req)
    if ((req.URL.Path == "/") && (redir == 0)) {
        if access {
            http.Redirect(w, req, Homepage, 302)
            return
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

    if !access {
        http.Redirect(w, req, Loginpage, http.StatusSeeOther)
        return
    }

    path := filepath.Join(afs.Path, req.URL.Path)
    http.ServeFile(w, req, path)
}
