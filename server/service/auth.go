package service

import (
        "log"
        "net/http"
)

func GetAuth(w http.ResponseWriter, req *http.Request) {
    var status int

    log.Println("GET /api/auth")
    AddBasicHeaders(w);

    if !CheckApiKey(req) {
        status = 400
	return
    }

    access := GetCookie(req)
    if !access {
       status = 401
    } else {
       status = 201
    }

    w.WriteHeader(status)
}
