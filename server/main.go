package main

import (
    "net/http"
    "log"
    "database/sql"
    "fmt"
    "io/ioutil"
    "os"
    "strconv"

    "ktrn.com/service"
    "ktrn.com/dbhandler"

    _ "github.com/lib/pq"
)

type Config struct {
  host     string
  port     int
  user     string
  password string
  dbname   string
}

func main() {
    var err error

    log.SetPrefix("[LOG] ")
    log.SetFlags(3)

    db_pass_file, err := ioutil.ReadFile(os.Getenv("DB_PASSWORD_FILE"))
    if err != nil {
        log.Println("ERROR! .db-secret not loaded")
        log.Fatal(err)
    }

    port_value, _ := strconv.Atoi(os.Getenv("DB_PORT"))

    Cfg := Config{
        host:     os.Getenv("DB_HOST"),
        port:     port_value,
        user:     os.Getenv("DB_USER"),
        password: string(db_pass_file),
        dbname:   os.Getenv("DB_NAME"),
    }

    log.Printf("Server started successfully")

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
                            "password=%s dbname=%s sslmode=disable",
    Cfg.host, Cfg.port, Cfg.user, Cfg.password, Cfg.dbname)
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
    http.HandleFunc("/api/checkLogin", service.PostCheckLogin)
    http.HandleFunc("/api/sendCode", service.PostSendCode)
    http.HandleFunc("/api/validateCode", service.PostValidateCode)
    http.HandleFunc("/api/getFile", service.GetGetFile)
    http.HandleFunc("/api/addPass", service.GetAddPass)
    http.HandleFunc("/api/multiGeneratePass", service.PostMultiGeneratePass)
    http.HandleFunc("/api/generatePass", service.PostGeneratePass)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
