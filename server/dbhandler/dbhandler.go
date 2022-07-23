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
        log.Println("13: Error: CheckSendCodeAttempts Query")
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&id)
        if err != nil {
            log.Println("14: Error: CheckSendCodeAttempts row iterations")
            return attempts, err
        }
        attempts++
    }

    err = rows.Err()
    if err != nil {
        log.Println("15: Error: CheckSendCodeAttempts row end")
        return attempts, err
    }

    return attempts, nil
}

func IncreaseAttempts(codeId string) error {
    _, err := DB.Exec(`UPDATE codes SET attempts = attempts + 1, updated_at = NOW()
                       WHERE id = $1;`, codeId)
    return err
}

func InsertSession(userId string, sessionId string, host string) error {
    _, err := DB.Exec(`INSERT INTO sessions (id, user_id, created_at, origin)
                       VALUES ($1, $2, now(), $3);`, sessionId, userId, host)
    return err
}

func CheckSession(sessionId string, host string) bool {
    var origin string

    row := DB.QueryRow(`SELECT origin FROM sessions WHERE id=$1;`, sessionId)

    switch err := row.Scan(&origin); err {
        case nil:
           if origin != host {
               log.Println("dbhandler.CheckSession: origin != host")
	       log.Println("origin in DB: ", origin)
	       log.Println("host: ", host)
	       return false
	   } else {
	        return true
	   }
        default:
            return false
    }
}

func CheckRole(sessionId string, role string) bool {
    var userId string
    var userRole string

    row := DB.QueryRow(`SELECT user_id FROM sessions WHERE id=$1;`, sessionId)

    err := row.Scan(&userId)
    if err != nil {
        return false
    }

    row = DB.QueryRow(`SELECT role FROM users WHERE id=$1;`, userId)

    err = row.Scan(&userRole)
    if err != nil {
        return false
    }

    if userRole != role {
        return false
    }
    return true
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
