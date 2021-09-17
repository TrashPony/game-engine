package auth

import (
	"encoding/json"
	"github.com/TrashPony/game_engine/src/dbConnect"
	"log"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type message struct {
	Login    string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//t, _ := template.ParseFiles("static/registration/index.html")
		//t.Execute(w, nil)
	}
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var msg message
		err := decoder.Decode(&msg)
		if err != nil {
			panic(err)
		}

		if msg.Login == "" || msg.Email == "" || msg.Password == "" || msg.Confirm == "" {
			resp := response{Success: false, Error: "form is empty"}
			json.NewEncoder(w).Encode(resp)
		} else {
			if msg.Password == msg.Confirm {

				checkLogin := checkAvailableLogin(msg.Login)
				checkEmail := true //checkAvailableEmail(msg.Email)

				if checkLogin && checkEmail {
					SuccessRegistration(msg.Login, msg.Email, msg.Password, 0)
					resp := response{Success: true, Error: ""}
					json.NewEncoder(w).Encode(resp)
				} else {
					if !checkLogin {
						resp := response{Success: false, Error: "login busy"} // error "этот логин уже занят"
						json.NewEncoder(w).Encode(resp)
					} else {
						resp := response{Success: false, Error: "email busy"}
						json.NewEncoder(w).Encode(resp)
					}
				}
			} else {
				resp := response{Success: false, Error: "password error"}
				json.NewEncoder(w).Encode(resp)
			}
		}
	}
}

func checkAvailableLogin(login string) (checkLogin bool) {
	user := GetUsersByName(login)

	if user.Name != "" {
		checkLogin = false
	} else {
		checkLogin = true
	}

	return
}

func checkAvailableEmail(email string) bool {

	user := GetUsersByMail(email)
	// если пользователь с таким мылом существуте то значит проверка не пройдена
	if user.Mail != "" {
		return false
	} else {
		return true
	}
}

func SuccessRegistration(login, email, password string, vkID int) {
	hashPassword := ""

	if password != "" {
		var err error
		hashPassword, err = HashPassword(login, password)
		if err != nil {
			println(err)
		}
	}

	tx, err := dbConnect.GetDBConnect().Begin()
	defer tx.Rollback()

	userID := 0

	err = tx.QueryRow(`INSERT INTO users (name, password, mail, vk_user_id) 
		VALUES ($1, $2, $3, $4) RETURNING id`, login, hashPassword, email, vkID).Scan(&userID)
	if err != nil {
		log.Fatal("registration new user " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("registration" + err.Error())
	}
}
