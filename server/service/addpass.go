package service

import (
        "net/http"
        "log"
        "encoding/json"
        "strings"
        "os"
	"bufio"
	"strconv"
)

func suggestPassword() (int, string) {
    max, err := FindMaxLoginIndex()
    if err != nil {
        return 500, ""
    }

    login := ConstructNewLogin(max)

    return 200, login
}

func ConstructNewLogin(max int) string {
    login := "bz" + strconv.Itoa(max + 1)
    return login
}

func FindMaxLoginIndex() (int, error) {
    var record string
    var split []string
    var max int

    max = 0

    file, err := os.Open(Filename)
    if err != nil {
        log.Println("27: Error while opening passwords file")
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        record = scanner.Text()
        if record[0:2] == "bz" {
            split = strings.Split(record,":")
            s := strings.Split(split[0], "bz")
            num, _ := strconv.Atoi(s[1])

            if num > max {
                max = num
            }
        }
    }

    if err := scanner.Err(); err != nil {
        log.Println("28: Error while parsing passwords file")
        return 0, err
    }

    return max, nil
}

func GetAddPass(w http.ResponseWriter, r *http.Request) {
    log.Println("POST /api/addPass")
    AddBasicHeaders(w);

    if (r.Method== OPTIONS) {
        CorsHandler(w, r)

        return
    }

    if !GetCookie(r) {
        w.WriteHeader(401)
        return
    }

    adminStatus := CheckRole(r, "admin")
    if adminStatus != 200 {
        adminStatus = CheckRole(r, "service")
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

    gitStatus, gitMsg := GitPullRepo()

    if gitStatus != 200 {
        w.WriteHeader(gitStatus)

        answer := Answer {
                              Status: gitStatus,
                              Response: gitMsg,
                         }
        json.NewEncoder(w).Encode(answer)

        return
    }

    status, login := suggestPassword()

    w.WriteHeader(status)

    answer := Answer {
                            Status: status,
                            Login: login,
                     }
    json.NewEncoder(w).Encode(answer)
}
