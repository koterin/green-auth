package dbhandler

import (
        "database/sql"
        "log"
)

var DB *sql.DB

func FindUser(email string) (uint, error) {
    var uuid string

    row := DB.QueryRow(`SELECT id FROM users WHERE email=$1;`, email)

    switch err := row.Scan(&uuid); err {
        case sql.ErrNoRows:
            return 0, nil
        case nil:
            return 1, nil
        default:
            return 0, err
    }
}

func FindCode(user_id string) (string, string, int, error) {
    var code string
    var codeId string
    var attempts int

    row := DB.QueryRow(`SELECT code, id, attempts FROM codes
                        WHERE user_id=$1
                        ORDER BY created_at DESC
                        LIMIT 1;`, user_id)

    switch err := row.Scan(&code, &codeId, &attempts); err {
        case nil:
            return code, codeId, attempts, nil
        default:
            return "", "", 0, err
    }
}

func GetUserId(email string) (string, error) {
    var uuid string

    row := DB.QueryRow(`SELECT id FROM users WHERE email=$1;`, email)

    switch err := row.Scan(&uuid); err {
        case nil:
            return uuid, nil
        default:
            return "", err
    }
}

func GetChatId(user_id string) (string, error) {
    var chat_id string

    row := DB.QueryRow(`SELECT chat_id FROM users WHERE id=$1;`, user_id)

    switch err := row.Scan(&chat_id); err {
        case nil:
            return chat_id, nil
        default:
            return "", err
    }
}

func AddCode(code string, user_id string) error {
    _, err := DB.Exec(`INSERT INTO codes (id, code, attempts, user_id, created_at, updated_at)
                        VALUES (DEFAULT, $1, $2, $3, DEFAULT, DEFAULT);`, code, 0, user_id)
    return err
}

func CheckSendCodeAttempts(user_id string) (int, error) {
    var id string
    var attempts int

    attempts = 1
    rows, err := DB.Query(`SELECT id FROM codes WHERE user_id = $1 AND
                           created_at > (NOW() - INTERVAL '5 minutes');`, user_id)
    if err != nil {
        log.Println("1: Error: CheckSendCodeAttempts Query")
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&id)
        if err != nil {
            log.Println("2: Error: CheckSendCodeAttempts row iterations")
            return attempts, err
        }
        attempts++
    }

    err = rows.Err()
    if err != nil {
        log.Println("3: Error: CheckSendCodeAttempts row end")
        return attempts, err
    }

    return attempts, nil
}

func IncreaseAttempts(codeId string) error {
    _, err := DB.Exec(`UPDATE codes SET attempts = attempts + 1, updated_at = NOW()
                       WHERE id = $1;`, codeId)
    return err
}

func InsertSession(userId string, sessionId string) error {
    _, err := DB.Exec(`INSERT INTO sessions (id, user_id, created_at)
                       VALUES ($1, $2, now());`, sessionId, userId)
    return err
}

func CheckSession(sessionId string) error {
    var userId string

    row := DB.QueryRow(`SELECT user_id FROM sessions WHERE id=$1;`, sessionId)

    switch err := row.Scan(&userId); err {
        case nil:
            return nil
        default:
            return err
    }
}

func CheckAdminSession(sessionId string) (int, error) {
    var userId string
    var role string

    row := DB.QueryRow(`SELECT user_id FROM sessions WHERE id=$1;`, sessionId)

    err := row.Scan(&userId)
    if err != nil {
        return 0, err
    }

    row = DB.QueryRow(`SELECT role FROM users WHERE id=$1;`, userId)

    err = row.Scan(&role)
    if err != nil {
        return 0, err
    }

    if ((role != "admin") && (role != "service")) {
        return 0, nil
    } else {
        return 1, nil
    }
}

func CheckDevSession(sessionId string) (int) {
    var userId string
    var role string

    row := DB.QueryRow(`SELECT user_id FROM sessions WHERE id=$1;`, sessionId)

    err := row.Scan(&userId)
    if err != nil {
        return 0
    }

    row = DB.QueryRow(`SELECT role FROM users WHERE id=$1;`, userId)

    err = row.Scan(&role)
    if err != nil {
        return 0
    }

    if ((role != "dev") && (role != "admin")) {
        return 0
    } else {
        return 1
    }
}

/*
func ifNull(session string) sql.NullString {
    if session == "NULL" {
        return sql.NullString{};
    }

    return sql.NullString {
         String: session,
         Valid: true,
    }
}
*/
