package service

import (
        "net/http"
        "log"
        "encoding/json"
	"os/exec"
	"io/ioutil"
)

func GitAdd() int {
    cmd := exec.Command("git", "-C", "./servicepswd", "add", ".")
    err := cmd.Run()

    if err != nil {
        log.Println("32: Error while adding changes to passwords file")

        return 1
    }

    return GitCommit()
}

func GitCommit() int {
    cmd := exec.Command("git", "-C", "./servicepswd", "commit", "-m", "\"new pass added\"")
    err := cmd.Run()

    if err != nil {
        log.Println("33: Error while committing file")

        return 1
    }

    return GitPush()
}

func GitPush() int {
    cmd := exec.Command("git", "-C", "./servicepswd", "push")
    err := cmd.Run()

    if err != nil {
        log.Println("34: Error while pushing changes")

        return 1
    }

    return 0
}

func GitPullRepo() (int, string) {
    cmd := exec.Command("git", "-C", "./servicepswd", "pull")
    err := cmd.Run()

    if err != nil {
        log.Println("28: Error while pulling passwords file")

        return 500, INTERNAL_ERROR_MSG
    }

    b, err := ioutil.ReadFile(Filename)

    if err != nil {
        log.Println("27: Error while opening passwords file")

        return 500, INTERNAL_ERROR_MSG
    }

    strBody := string(b)

    return 200, strBody
}

func GetGetFile(w http.ResponseWriter, req *http.Request) {
    log.Println("POST /api/getFile")
    AddBasicHeaders(w);

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

    status, msg := GitPullRepo()
    w.WriteHeader(status)

    answer := Answer {
                         Status: status,
                         Response: msg,
                     }
    json.NewEncoder(w).Encode(answer)
}
