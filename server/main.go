package main

import (
    "net/http"
    "log"
    "database/sql"
    "fmt"
    "io/ioutil"

    _ "github.com/lib/pq"

    "green-auth/service"
    "green-auth/dbhandler"
)


func main() {
    log.SetPrefix("[LOG] ")
    log.SetFlags(3)

    service.ValidateEnv()

    db_pass, _ := ioutil.ReadFile(service.Args.DB_PASSWORD_FILE)

    log.Printf("Server started successfully")

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
                            "password=%s dbname=%s sslmode=disable",
                            service.Args.DB_HOST, service.Args.DB_PORT, service.Args.DB_USER,
                            string(db_pass), service.Args.DB_NAME)

    var err error
    dbhandler.DB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Println("00: DB credentials are invalid")
        panic(err)
    }

    err = dbhandler.DB.Ping()
    if err != nil {
        log.Println("01: DB connection is not established")
        panic(err)
    }

    log.Println("Successfully connected to DB")

    http.Handle("/", service.AuthFileServerHandler("./static/"))
    http.HandleFunc("/api/checkEmail", service.PostCheckEmail)
    http.HandleFunc("/api/sendCode", service.PostSendCode)
    http.HandleFunc("/api/validateCode", service.PostValidateCode)
    http.HandleFunc("/api/auth", service.GetAuth)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
