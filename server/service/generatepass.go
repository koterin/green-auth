package service

import (
       "net/http"
       "log"
       "fmt"
       "time"
       "encoding/json"
       "os"
       "bufio"
       "strings"
       "math/rand"
       "crypto/md5"
)

func CheckUniqueLogin(login string) (int, string) {
    var record string
    var split []string

    if login == "" {
       return 204, "Empty login"
    }

    for _, ch := range login {
       if ch == 32 {
          return 400, "No spaces allowed"
       }
    }

    file, err := os.Open(Filename)
    if err != nil {
       log.Println("27: Error while opening passwords file")
       return 500, INTERNAL_ERROR_MSG
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
       record = scanner.Text()
       split = strings.Split(record, ":")

       if split[0] == login {
          return 409, "Login already exists"
       }
    }

    if err := scanner.Err(); err != nil {
       log.Println("28: Error while parsing passwords file")
       return 500, INTERNAL_ERROR_MSG
    }

    return 200, "Allowed"
}

func WritePassToFile(record string) (int, string) {
    file, err := os.OpenFile(Filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
    if err != nil {
       log.Println("27: Error while opening passwords file")

	    return 500, INTERNAL_ERROR_MSG
    }
    defer file.Close()

    _, err = file.WriteString(record)

    if err != nil {
       log.Println("29: Error while writing into passwords file")
       log.Println(err)
       file.Close()

       return 500, INTERNAL_ERROR_MSG
    }

    file.Close()

    return 200, ""
}

func GitComplete() (int, string) {
    res := GitAdd()
    if res != 0 {
        log.Println("30: Error while pushing changes")
        // TODO: GitRestore();

        return 500, INTERNAL_ERROR_MSG
    }

    return 200, "Success"
}

func CreateNewRecord(login string, pass string) string {
    record := login + ":" + "$apr1$" + md5pass(pass)
    return record
}

func GeneratePass(length int) string {
    const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

    b := make([]byte, length)
    for i := range b {
       b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
    }

    return string(b)
}

func md5pass(pass string) string {
    data := []byte(pass)
    result := fmt.Sprintf("%x", md5.Sum(data))

    return result
}

func PostGeneratePass(w http.ResponseWriter, req *http.Request) {
    var respdata ResponseData

    rand.Seed(time.Now().UnixNano())

    log.Println("POST /api/generatePass")
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

    checkStatus, checkMsg := CheckUniqueLogin(respdata.Login)
    if checkStatus != 200 {
        w.WriteHeader(checkStatus)

        answer := Answer {
                             Status: checkStatus,
                             Response: checkMsg,
                         }
        json.NewEncoder(w).Encode(answer)

        return
    }

    pass := GeneratePass(PASSWORD_LENGTH)
    record := CreateNewRecord(respdata.Login, pass)
    wStatus, wMsg := WritePassToFile(record + "\n")
    if wStatus != 200 {
        w.WriteHeader(wStatus)

        answer := Answer {
                             Status: wStatus,
                             Response: wMsg,
                         }
        json.NewEncoder(w).Encode(answer)

        return
    }

    status, msg := GitComplete()

    w.WriteHeader(status)

    answer := Answer {
                        Status: status,
                        Response: msg,
                        Password: pass,
                     }
    json.NewEncoder(w).Encode(answer)
}
