package vk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const urlAPIMethod = "https://api.vk.com/method/"

const (
	// LangRU -- Русский язык
	LangRU Lang = iota
	// LangUA -- Украинский язык
	LangUA Lang = iota
	// LangBE -- Белорусский язык
	LangBE Lang = iota
	// LangEN -- Английский язык
	LangEN Lang = iota
	// LangES -- Испанский язык
	LangES Lang = iota
	// LangFI -- Финский язык
	LangFI Lang = iota
	// LangDE -- Немецкий язык
	LangDE Lang = iota
	// LangIT -- Итальянский язык
	LangIT Lang = iota
)

// Lang -- Язык запроса
type Lang int

// CheckUserTokenResponse -- Структура ответа сервера на запос проверки токена пользователя
type CheckUserTokenResponse struct {
	Response struct {
		Success int   `json:"success"`
		UserID  int64 `json:"user_id"`
		Date    int64 `json:"date"`
		Expire  int64 `json:"expire"`
	} `json:"response"`
	Error struct {
		Code    int    `json:"error_code"`
		Message string `json:"error_msg"`
	} `json:"error"`
}

// CheckUserToken -- Запрос проверка токена пользователя
func (vk *API) CheckUserToken(userTokenAccess string) (CheckUserTokenResponse, error) {
	var data CheckUserTokenResponse
	params := map[string]string{
		"access_token": vk.AccessToken,
		"token":        userTokenAccess,
	}
	body, err := vk.Request("secure.checkToken", params, false)
	if err != nil {
		return data, err
	}
	data, err = parseCheckUserTokenResponse(body)
	if err != nil {
		return data, err
	}
	if len(data.Error.Message) > 0 {
		return data, fmt.Errorf("Error [secure.checkToken] code = %d (%s)", data.Error.Code, data.Error.Message)
	}
	if data.Response.Success != 1 {
		return data, fmt.Errorf("Error [secure.checkToken] success = %d", data.Response.Success)
	}
	return data, nil
}

func parseCheckUserTokenResponse(response []byte) (CheckUserTokenResponse, error) {
	var model CheckUserTokenResponse
	if err := json.Unmarshal(response, &model); err != nil {
		return model, err
	}
	return model, nil
}

type (
	// GetUserProfileResponse -- Структура ответа сервера на запос получения профайоа пользователя
	GetUserProfileResponse struct {
		Response []UserInfo `json:"response"`
		Error    struct {
			Code    int    `json:"error_code"`
			Message string `json:"error_msg"`
		} `json:"error"`
	}

	// UserInfo -- Модель данных пользователя
	UserInfo struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Photo     string `json:"photo_200"`
	}
)

// GetUserProfile -- Запрос на получение профайла пользователя
func (vk *API) GetUserProfile(userID int64, userToken string, lang Lang) (UserInfo, error) {
	var data UserInfo
	params := map[string]string{
		"user_ids":  strconv.FormatInt(userID, 10),
		"fields":    "photo_200",
		"name_case": "Nom",
		"lang":      strconv.Itoa(int(lang)),
	}
	body, err := vk.Request("users.get", params, true)
	if err != nil {
		return data, err
	}
	model, err := parseGetUserProfileResponse(body)
	if err != nil {
		return data, err
	}
	if len(model.Error.Message) > 0 {
		return data, fmt.Errorf("Error [users.get] code = %d (%s)", model.Error.Code, model.Error.Message)
	}
	if len(model.Response) == 0 {
		return data, fmt.Errorf("Error [users.get] answer empty")
	}
	data = model.Response[0]
	return data, nil
}

func parseGetUserProfileResponse(response []byte) (GetUserProfileResponse, error) {
	var model GetUserProfileResponse
	if err := json.Unmarshal(response, &model); err != nil {
		return model, err
	}
	return model, nil
}

// Request -- Запрос к VK API
func (vk *API) Request(methodName string, extra map[string]string, anonymousRequest bool) ([]byte, error) {
	params := url.Values{
		"v": {apiVersion},
	}
	if !anonymousRequest {
		params.Add("client_id", vk.ClietnID)
		params.Add("client_secret", vk.ClientSecret)
	}
	for key, value := range extra {
		params.Add(key, value)
	}
	if vk.debug {
		log.Println("[VK]", methodName, "<--", params)
	}
	resp, err := http.PostForm(urlAPIMethod+methodName, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if vk.debug {
		log.Println("[VK]", methodName, "-->", string(body))
	}
	return body, nil
}
