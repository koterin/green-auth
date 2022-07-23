package service

import (
        "net/http"
        "log"
        "encoding/json"
        "time"
        "net/url"

        "github.com/google/uuid"

        "green-auth/dbhandler"
)

func generateCookie(email string, host string) (*http.Cookie, int, string) {
    userId, err := dbhandler.GetUserId(email)
    if err != nil {
        log.Println("10: Error while executing .GetUserId()")

        return nil, 500, INTERNAL_ERROR_MSG
    }

    sessionId := uuid.New()
    err = dbhandler.InsertSession(userId, sessionId.String(), host)

    if err != nil {
        log.Println("23: Error while executing .InsertSession()")

        return nil, 500, INTERNAL_ERROR_MSG
    }

    cookie := &http.Cookie{
        Name:   "sessionId",
        Value:  sessionId.String(),
        Expires: (time.Now().Add(time.Duration(168) * time.Hour)),
        Path: "/",
        HttpOnly: true,
        Secure: true,
    }

    return cookie, 200, "Code is valid"
}

func ValidateCode(email string, code string) (int, string) {
    userId, err := dbhandler.GetUserId(email)

    if err != nil {
        log.Println("10: Error while executing .GetUserId()")

        return 500, INTERNAL_ERROR_MSG
    }

    trueCode, codeId, attempts, err := dbhandler.FindCode(userId)

    if err != nil {
        log.Println("20: Error while executing .FindCode")

        return 500, INTERNAL_ERROR_MSG
    }

    if attempts >= 5 {
        return 429, "Too many attempts"
    }

    err = dbhandler.IncreaseAttempts(codeId)

    if err != nil {
        log.Println("21: Error while executing .IncreaseAttempts")

        return 500, INTERNAL_ERROR_MSG
    }

    if (code != trueCode) {
        return 400, "Wrong Code"
    }

    return 200, "Code is valid"
}

func PostValidateCode(w http.ResponseWriter, req *http.Request) {
    var respdata ResponseData

    log.Println("POST /api/validateCode")
    AddBasicHeaders(w);

    err := ReadJson(w, req, &respdata)
    if err != nil {
        log.Print("07: Error unmarshalling JSON")
        w.WriteHeader(400)

        return
    }

    checkStatus, checkMsg := CheckEmail(respdata.Email)

    host := req.Host

    redirect := req.Header.Get("X-Redirect-To")
    if redirect != "null" {
        redirectUrl, err := url.Parse(redirect)
    if err != nil {
        redirect = "null"
    } else {
        host = redirectUrl.Host
        }
    }

    if checkStatus == 200 {
        codeStatus, codeMsg := ValidateCode(respdata.Email, respdata.Code)

        if codeStatus == 200 {
            cookie, status, msg := generateCookie(respdata.Email, host);

            if status == 200 {
                if redirect != "null" {
                    w.Header().Set("X-Green-Token", cookie.Value)
                    w.WriteHeader(200)

                return
            } else {
                http.SetCookie(w, cookie)
                }
            }

            w.WriteHeader(status)

            answer := Answer {
                                Status: status,
                                Response: msg,
                             }
            json.NewEncoder(w).Encode(answer)
        } else {
            w.WriteHeader(codeStatus)

            answer := Answer {
                            Status: codeStatus,
                            Response: codeMsg,
                         }
            json.NewEncoder(w).Encode(answer)
        }
    } else {
        w.WriteHeader(checkStatus)

        answer := Answer {
                            Status: checkStatus,
                            Response: checkMsg,
                         }
        json.NewEncoder(w).Encode(answer)
    }
}
