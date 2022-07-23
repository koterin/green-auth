package service

import (
       "net/http"
       "log"
       "encoding/json"
       "strconv"
)

func multiGeneratePass(q string) (int, string) {
    var pass string
    var index int
    var login string
    var body string
    var record string
    var recordBody string

    qint, _ := strconv.Atoi(q)
    index, err := FindMaxLoginIndex()

    if err != nil {
        return 500, ""
    }

    for i := 0; i < qint; i++ {
        login = ConstructNewLogin(index)
        pass = GeneratePass(PASSWORD_LENGTH)
        body += login + ": " + pass + "\n"

        record = CreateNewRecord(login, pass)
        recordBody += record + "\n"
        index += 1
    }

    wStatus, _ := WritePassToFile(recordBody)
    if wStatus != 200 {
        return wStatus, ""
    }
    status, _ := GitComplete()

    return status, body
}

func checkQuantity(q string) (int, string) {
    qint, err := strconv.Atoi(q)

    if err != nil {
        return 400, "Bad Credentials"
    }

    if ((qint < 2) || (qint > 100)) {
        return 400, "Wrong quantity"
    }

    return 200, "Allowed"
}

func PostMultiGeneratePass(w http.ResponseWriter, req *http.Request) {
    var respdata ResponseData

    log.Println("POST /api/multiGeneratePass")
    AddBasicHeaders(w)

    if !GetCookie(req) {
        w.WriteHeader(401)
        return
    }

    adminStatus := CheckRole(req, "admin")
    if adminStatus != 200 {
        adminStatus = CheckRole(req, "service")
    }

    if adminStatus != 200 {
        w.WriteHeader(adminStatus)

        answer := Answer {
                            Status: adminStatus,
                            Response: "User not allowed",
                         }
        json.NewEncoder(w).Encode(answer)

        return
    }

    err := ReadJson(w, req, &respdata)
    if err != nil {
       log.Print("07: Error unmarshalling JSON")
       w.WriteHeader(400)

       return
    }

    checkStatus, checkMsg := checkQuantity(respdata.Quantity)
    if checkStatus != 200 {
        w.WriteHeader(checkStatus)

        answer := Answer {
                             Status: checkStatus,
                             Response: checkMsg,
                         }
        json.NewEncoder(w).Encode(answer)

        return
    }

    status, body := multiGeneratePass(respdata.Quantity)

    w.WriteHeader(status)

    answer := Answer {
                         Status: status,
                         Response: body,
                     }
    json.NewEncoder(w).Encode(answer)
}
