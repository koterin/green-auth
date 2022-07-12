package service

import (
        "net/http"
        "log"
        "encoding/json"

        "ktrn.com/dbhandler"
)

func CheckLogin(login string) (int, string) {
    res, err := dbhandler.FindUser(login)
    if err != nil {
       log.Println("06: Error while executing .FindUser()")

        return 500, INTERNAL_ERROR_MSG
    }

    if res == 0 {
        log.Println("07: Login not found in allowed: ", login)

        return 400, "Login not found in allowed"
    } else {
        return 200, "Login accepted"
    }
}

func PostCheckLogin(w http.ResponseWriter, req *http.Request) {
    var respdata ResponseData

    log.Println("POST /api/checkLogin")
    AddBasicHeaders(w);

    if (req.Method== OPTIONS) {
        CorsHandler(w, req)

        return
    }

    err := ReadJson(w, req, &respdata)
    if err != nil {
        log.Print("08: Error unmarshalling JSON")

        return
    }

    status, msg := CheckLogin(respdata.Login)

    w.WriteHeader(status)

    answer := Answer {
                        Status: status,
                        Response: msg,
                     }

    json.NewEncoder(w).Encode(answer)
}
