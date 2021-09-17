package auth

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	_const "github.com/TrashPony/game_engine/src/const"
	"log"
	"net/http"
	"strconv"
)

var (
	redirectUri  = _const.Config.GetParams("redirectUri")
	clientID     = _const.Config.GetParams("clientID")
	clientSecret = _const.Config.GetParams("clientSecret")
	apiSecret    = _const.Config.GetParams("vkAppKey")
)

const (
	scope   = ""
	display = "popup"
)

type VkUser struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
	ScreenName  string `json:"screen_name"`
}

func getVaAuthUrl() string {
	return "https://oauth.vk.com/authorize?client_id=" + clientID + "&display=" + display + "&redirect_uri=" + redirectUri + "&scope=" + scope
}

func VkGetUrlToAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//https://oauth.vk.com/authorize?client_id=7718842&display=page&redirect_uri=http://localhost:8080/vk-oauth&scope=email
		jsonData, err := json.Marshal(map[string]string{"o_auth_url": getVaAuthUrl()})
		if err != nil {
			//todo
			log.Print(err)
		}

		_, err = w.Write(jsonData)
		if err != nil {
			//todo
			log.Print(err)
		}
	}
}

func VkAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		code := r.URL.Query().Get("code")
		errorAuth := r.URL.Query().Get("error")

		if errorAuth != "" {
			json.NewEncoder(w).Encode(response{Success: false, Error: errorAuth})
			return
		}

		vkUser, err := GetVkUser(code)
		if err != nil {
			json.NewEncoder(w).Encode(response{Success: false, Error: err.Error()})
			return
		}

		err = getScreenName(vkUser)
		if err != nil {
			json.NewEncoder(w).Encode(response{Success: false, Error: err.Error()})
			return
		}

		if vkUser.UserID > 0 {
			vkAuth(w, r, vkUser)
		} else {
			json.NewEncoder(w).Encode(response{Success: false, Error: "Не удалось"})
		}
	}
}

func VkAppLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		authKey := r.URL.Query().Get("auth_key")
		apiID := r.URL.Query().Get("api_id")
		viewerID := r.URL.Query().Get("viewer_id")
		accessToken := r.URL.Query().Get("access_token")

		//https://vk.com/dev/apps_init?f=3.%20auth_key - проверяем по полученым парметрам валидность авторизации
		//auth_key = md5(api_id + '_' + viewer_id + '_' + api_secret)
		checkAuthKey := md5.Sum([]byte(apiID + "_" + viewerID + "_" + apiSecret))

		if authKey == fmt.Sprintf("%x", checkAuthKey) {

			userID, err := strconv.Atoi(viewerID)
			if err != nil {
				json.NewEncoder(w).Encode(response{Success: false, Error: "Не удалось"})
			}

			vkUser := &VkUser{
				AccessToken: accessToken,
				UserID:      userID,
			}

			err = getScreenName(vkUser)
			if err != nil {
				json.NewEncoder(w).Encode(response{Success: false, Error: "Не удалось"})
			}

			vkAuth(w, r, vkUser)
		} else {
			json.NewEncoder(w).Encode(response{Success: false, Error: "Не удалось"})
		}
	}
}

func vkAuth(w http.ResponseWriter, r *http.Request, vkUser *VkUser) {
	// проверяем пользователя по UserID в бд по полю vk_user_id
	user := GetUsersByVkID(vkUser.UserID)
	if user.Id == 0 {

		vkUser.ScreenName = "vk_id" + strconv.Itoa(vkUser.UserID)

		SuccessRegistration(vkUser.ScreenName, "", "", vkUser.UserID)
		user = GetUsersByVkID(vkUser.UserID)
	}

	GetCookie(w, r, user)
}

func getScreenName(user *VkUser) error {
	// https://api.vk.com/method/users.get?user_ids=26106&fields=screen_name&access_token=78180cc2ef2cbb2adc2b7fd79927bbac37c2e2e32d90e8ef7f9d653ae3a2174765127d32a853890de3d42&v=5.126

	req, err := http.NewRequest("GET", "https://api.vk.com/method/users.get", nil)
	if err != nil {
		log.Print(err)
		return err
	}

	q := req.URL.Query()
	q.Add("user_ids", strconv.Itoa(user.UserID))
	q.Add("fields", "screen_name")
	q.Add("access_token", user.AccessToken)
	q.Add("v", "5.126")

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}

	//{"response":[{"first_name":"Роберт","id":7523894,"last_name":"Хромов","can_access_closed":true,"is_closed":false,"screen_name":"trashpony13"}]}
	type Response struct {
		ID         int    `json:"id"`
		LastName   string `json:"last_name"`
		FirstName  string `json:"first_name"`
		ScreenName string `json:"screen_name"`
	}

	type Responses struct {
		Responses []Response `json:"response"`
	}

	var r Responses
	err = json.NewDecoder(response.Body).Decode(&r)
	if err != nil {
		log.Print(err)
		return err
	}

	if len(r.Responses) > 0 {
		user.ScreenName = r.Responses[0].ScreenName
	}

	return nil
}

func GetVkUser(code string) (*VkUser, error) {
	req, err := http.NewRequest("GET", "https://oauth.vk.com/access_token", nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("redirect_uri", redirectUri)
	q.Add("code", code)

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	//fmt.Printf("Encoded URL is %q\n", req.URL.String())

	vkUser := &VkUser{}
	err = json.NewDecoder(response.Body).Decode(&vkUser)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return vkUser, nil
}
