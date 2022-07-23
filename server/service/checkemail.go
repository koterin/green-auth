package service

import (
        "net/http"
        "log"
        "encoding/json"

        "green-auth/dbhandler"
)

func CheckEmail(email string) (int, string) {
    res, err := dbhandler.FindUser(email)
    if err != nil {
       log.Println("08: Error while executing .FindUser()")

        return 500, INTERNAL_ERROR_MSG
    }

    if res == 0 {
        log.Println("04: Email not found in allowed: ", email)

        return 400, "Email not found in allowed"
    } else {
        return 200, "Email accepted"
    }
}

func PostCheckEmail(w http.ResponseWriter, req *http.Request) {
    var respdata ResponseData

    log.Println("POST /api/checkEmail")
    AddBasicHeaders(w);

    if (req.Method== OPTIONS) {
        CorsHandler(w, req)

        return
    }

    err := ReadJson(w, req, &respdata)
    if err != nil {
        log.Print("05: Error unmarshalling JSON")

        return
    }

    status, msg := CheckEmail(respdata.Email)

    w.Header().Set("Origin", "http://127.0.0.1:8080")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    answer := Answer {
                        Status: status,
                        Response: msg,
                     }

    json.NewEncoder(w).Encode(answer)
}
