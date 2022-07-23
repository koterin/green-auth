package service

import (
        "io/ioutil"
        "net/url"

        "github.com/alexflint/go-arg"
)

var Args struct {
    DB_PASSWORD_FILE  string  `arg:"env"`
    DB_PORT           int     `arg:"env"`
    DB_HOST           string  `arg:"env"`
    DB_USER           string  `arg:"env"`
    DB_NAME           string  `arg:"env"`
    TG_API_URL        string  `arg:"env"`
    TG_BOT_KEY        string  `arg:"env"`
    API_KEY           string  `arg:"env"`
    HOST_URL          string  `arg:"env"`
    CODE_LENGTH       int     `arg:"env"`
}

var UnathorizedPages [10]string = [10]string{"/main.css",
                                            "/index.js",
                                            "/index.html",
                                            "/",
                                            "/authenticate.html",
                                            "/auth.js",
                                            "/theme.js",
                                            "/images/tg_qr.png",
                                            "/favicon.ico",
                                            "/wrongredirect.html"}

var Loginpage = Args.HOST_URL
var Homepage = Args.HOST_URL + "/home.html"
var WrongRedirectPage = Args.HOST_URL + "/wrongredirect.html"
var INTERNAL_ERROR_MSG = "Internal Error"

func ValidateEnv() {
    p := arg.MustParse(&Args)

    _, err := ioutil.ReadFile(Args.DB_PASSWORD_FILE)
    if err != nil {
        p.Fail("Path to .db-secret is empty or invalid")
    }

    if Args.DB_HOST == "" {
        p.Fail("DB_HOST env is not set")
    }

    if Args.DB_USER == "" {
        p.Fail("DB_USER env is not set")
    }

    if Args.DB_NAME == "" {
        p.Fail("DB_NAME env is not set")
    }

    if Args.TG_API_URL == "" {
        p.Fail("TG_API_URL env is not set")
    }

    if Args.TG_BOT_KEY == "" {
        p.Fail("TG_BOT_KEY env is not set")
    }

    if Args.API_KEY == "" {
        p.Fail("API_KEY env is not set")
    }

    if Args.HOST_URL == "" {
        p.Fail("HOST_URL env is not set")
    }

    _, err = url.Parse(Args.HOST_URL)
    if err != nil {
        p.Fail("HOST_URL env is invalid. Requried format: https://test.domain.com")
    }

    _, err = url.ParseRequestURI(Args.HOST_URL)
    if err != nil {
        p.Fail("HOST_URL env is invalid. Requried format: https://test.domain.com")
    }
}
