package service

import (
        "net/http"
        "log"
        "encoding/json"
        "os/exec"
        "io/ioutil"
)

func GitAdd() int {
    cmd := exec.Command("git", "-C", "./green-auth", "add", ".")
    err := cmd.Run()

    if err != nil {
        log.Println("15: Error while adding changes to passwords file")

        return 1
    }

    return GitCommit()
}

func GitCommit() int {
    cmd := exec.Command("git", "-C", "./green-auth", "commit", "-m", "\"new pass added\"")
    err := cmd.Run()

    if err != nil {
        log.Println("16: Error while committing file")

        return 1
    }

    return GitPush()
}

func GitPush() int {
    cmd := exec.Command("git", "-C", "./green-auth", "push")
    err := cmd.Run()

    if err != nil {
        log.Println("17: Error while pushing changes")

        return 1
    }

    return 0
}

func GitPullRepo() (int, string) {
    cmd := exec.Command("git", "-C", "./green-auth", "pull")
    err := cmd.Run()

    if err != nil {
        log.Println("18: Error while pulling passwords file")

        return 500, INTERNAL_ERROR_MSG
    }

    b, err := ioutil.ReadFile(Filename)

    if err != nil {
        log.Println("19: Error while opening passwords file")

        return 500, INTERNAL_ERROR_MSG
    }

    strBody := string(b)

    return 200, strBody
}

func GetGetFile(w http.ResponseWriter, req *http.Request) {
    log.Println("POST /api/getFile")
    AddBasicHeaders(w);

    if (req.Method == OPTIONS) {
        CorsHandler(w, req)

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
