package service

import (
        "net/http"
        "log"
        "encoding/json"
        "crypto/rand"
        "io"
        "time"

        "ktrn.com/dbhandler"
)

var PassClient = &http.Client{Timeout: 10 * time.Second}

var table = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func generate_code(length int) string {
	b := make([]byte, length)
	_, err := io.ReadAtLeast(rand.Reader, b, length)

    if err != nil {
		panic(err)
	}

    for i := 0; i < length; i++ {
		b[i] = table[int(b[i]) % len(table)]
	}

    return string(b)
}

func sendTelegramMsg(chat_id string, code string) {
    code = "`" + code + "`"
    requestURL := TELEGRAM_URL + "bot" + TELEGRAM_BOT_KEY + "/sendMessage" +
                  "?chat_id=" + chat_id + "&text=" + code + "&parse_mode=markdown"
    resp, err := PassClient.Get(requestURL)

    if err != nil {
        log.Println("21: Error while sending request to Telegram")

        return
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)

    if err != nil {
        log.Println("22: Error while reading TG response body")
        return
    }

    sb := string(body)
    sb = sb[6:10]

    if sb != "true" {
		log.Println("23: Error while sending code via Telegram")
	}
}

func SendCode(Login string) (int, string) {
    user_id, err := dbhandler.GetUserId(Login)

    if err != nil {
        log.Println("24: Error while executing .GetUserId()")

        return 500, INTERNAL_ERROR_MSG
    }

    attempts, err := dbhandler.CheckSendCodeAttempts(user_id)

    if err != nil {
        return 500, INTERNAL_ERROR_MSG
    }

    if attempts > 5 {
        log.Println("25: User made more than 5 SendCode attempts in 5 minutes: ", Email)

        return 429, "Attempts limit in period exceeded"
    }

    code := generate_code(CODE_LENGTH)

    err = dbhandler.AddCode(code, user_id);

    if err != nil {
        log.Println("26: Error while executing .AddCode() for user ", user_id)

        return 500, INTERNAL_ERROR_MSG
    }

    chat_id, err := dbhandler.GetChatId(user_id)

    if err != nil {
        log.Println("27: Error while executing .GetChatId()")

        return 500, INTERNAL_ERROR_MSG
    }

    sendTelegramMsg(chat_id, code)

    return 200, "Code sent"
}

func PostSendCode(w http.ResponseWriter, req *http.Request) {
    var respdata ResponseData

    log.Println("POST /api/sendCode")

    if (req.Method== OPTIONS) {
        CorsHandler(w, req)

        return
    }

    err := ReadJson(w, req, &respdata)

    if err != nil {
        log.Print("28: Error unmarshalling JSON")

        return
    }

    checkStatus, checkMsg := CheckEmail(respdata.Email)

    AddBasicHeaders(w)

    if checkStatus == 200 {
        status, msg := SendCode(respdata.Email)
        w.WriteHeader(status)

        answer := Answer {
                            Status: status,
                            Response: msg,
                         }
        json.NewEncoder(w).Encode(answer)
    } else {
        w.WriteHeader(checkStatus)
        answer := Answer {
                            Status: checkStatus,
                            Response: checkMsg,
                         }
        json.NewEncoder(w).Encode(answer)
    }
}
