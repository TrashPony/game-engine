package auth

import (
	"encoding/json"
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/dbConnect"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var cookieStore = sessions.NewCookieStore([]byte("dick, mountain, sky ray")) // мало понимаю в шифрование сессии внутри указан приватный ключь шифрования

func init() {
	if _const.Config.GetParams("cookiesSecure") == "true" {
		cookieStore.Options.SameSite = http.SameSiteNoneMode
		cookieStore.Options.Secure = true
	}
}

const cookieName = "infinity-key" // имя куки в браузере юзера

const (
	login = 1
	id    = 2
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//t, _ := template.ParseFiles("static/login/index.html")
		//t.Execute(w, nil)
	}
	if r.Method == "POST" { // получаем данные с фронтенда
		decoder := json.NewDecoder(r.Body)
		var msg message
		err := decoder.Decode(&msg)
		if err != nil {
			panic(err)
		}

		// отправляет эти данные на проверку если прошло то возвращает пользователя и пропуск
		user := GetUsersByName(msg.Login)

		passed := false
		if user.Password != "" {
			passed = CheckPasswordHash(msg.Login, msg.Password, user.Password)
		}

		if passed {
			//отправляет пользователя на получение токена подключения
			GetCookie(w, r, user)
		} else {
			json.NewEncoder(w).Encode(response{Success: false, Error: "not allow"})
			println("Соеденение не разрешено: не авторизован")
		}
	}
}

func GetCookie(w http.ResponseWriter, r *http.Request, user User) {
	// берет сеанс из браузера пользователя
	ses, _ := cookieStore.Get(r, cookieName)
	// если есть куки подписаные не правильным ключем то вылетает ошибка

	ses.Values[login] = user.Name // ложит данные в сессию
	ses.Values[id] = user.Id      // ложит данные в сессию

	//возвращает ответ с сохранение сессии в браузере
	err := cookieStore.Save(r, w, ses)

	//http.Redirect(w, r, "http://642e0559eb9c.sn.mynetname.net:8080/lobby/", 302)

	json.NewEncoder(w).Encode(response{Success: true, Error: ""})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func CheckCookie(w http.ResponseWriter, r *http.Request) (string, int) {
	// берет сеанс из браузера пользователя
	ses, err := cookieStore.Get(r, cookieName)
	// если есть куки подписаные не правильным ключем то вылетает ошибка
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "", 0
	}

	// смотрит пустая сессия или нет, если нет то присваивает переменной логин логин
	// ищет значение в сессии и присваивает переменной [Login] - ключь .(string) - тип данных ok - удалось ли получить

	login, ok := ses.Values[login].(string)
	id, ok := ses.Values[id].(int)

	if !ok { // если пустая то говорит что ты анонимус
		return "", 0
	}
	return login, id
}

func GetUsersByName(name string) User {
	rows, err := dbConnect.GetDBConnect().Query("Select id, name, mail, password FROM users WHERE name=$1 and vk_user_id=$2", name, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var user User

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Mail, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
	}
	return user
}

func GetUsersByVkID(vkID int) User {
	rows, err := dbConnect.GetDBConnect().Query("Select id, name, mail, password FROM users WHERE vk_user_id=$1", vkID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var user User

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Mail, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
	}
	return user
}

func GetUsersByMail(mail string) User {
	rows, err := dbConnect.GetDBConnect().Query("Select id, name, mail, password FROM users WHERE mail=$1", mail)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var user User

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Mail, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
	}
	return user
}

type User struct {
	Id       int
	Name     string
	Mail     string
	Password string
	VkID     int
}
