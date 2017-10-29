package vk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const urlAuthGetToken = "https://oauth.vk.com/access_token"

// OAuthResponse -- Структура ответа сервера на запрос "https://oauth.vk.com/access_token"
type OAuthResponse struct {
	Token            string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// OAuth -- Запрос авторизации.
// Если запрос успешно выполнен полученый токен будет добавлен в API.AccessToken
func (vk *API) OAuth() error {
	params := url.Values{
		"client_id":     {vk.ClietnID},
		"client_secret": {vk.ClientSecret},
		"grant_type":    {"client_credentials"},
		"v":             {apiVersion},
	}
	if vk.debug {
		log.Println("[VK] OAuth <--", params)
	}
	resp, err := http.PostForm(urlAuthGetToken, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	data, err := parseOAuthResponse(body)
	if err != nil {
		return err
	}
	if vk.debug {
		log.Printf("[VK] OAuth --> %+v\n", data)
	}
	if len(data.Error) > 0 {
		return fmt.Errorf("OAuth error code = %s (%s)", data.Error, data.ErrorDescription)
	}

	vk.AccessToken = data.Token
	return nil
}

func parseOAuthResponse(response []byte) (OAuthResponse, error) {
	var model OAuthResponse
	if err := json.Unmarshal(response, &model); err != nil {
		return model, err
	}
	return model, nil
}
