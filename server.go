package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/anninpavel/firebase-auth-vk/vk"
)

func main() {
	log.Println("FireBase Auth Vk Starting")

	env := NewEnvironment()
	server := &http.Server{
		Addr:           env.conf.Port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if env.vk.OAuth() != nil {
		log.Fatal("[Vk] error auth")
	}

	http.HandleFunc(env.conf.Path, env.authVk)
	log.Fatal(server.ListenAndServe())
}

// AuthVkResponce -- Структура ответа сервера на запрос авторизации через VK
type AuthVkResponce struct {
	Token    string `json:"token"`
	FullName string `json:"userName"`
	Avatar   string `json:"userAvatar"`
}

func (env *Environment) authVk(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Println(r.RemoteAddr, "--> AuthVk", "("+r.Method+")")
	defer func() {
		log.Println(r.RemoteAddr, "<-- AuthVk", "("+time.Now().Sub(startTime).String()+")")
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Required parameters: userId, accessToken", http.StatusBadRequest)
		return
	}
	vkUserID, err := strconv.ParseInt(r.Form.Get("userId"), 10, 64)
	vkUserToken := r.Form.Get("accessToken")
	userLanguage := r.Form.Get("lang")
	if err != nil || len(vkUserToken) == 0 {
		http.Error(w, "Required parameters: userId, accessToken", http.StatusBadRequest)
		return
	}

	checkToken, err := env.vk.CheckUserToken(vkUserToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if vkUserID != checkToken.Response.UserID {
		http.Error(w, "UserId does not match token", http.StatusBadRequest)
		return
	}

	user, err := env.vk.GetUserProfile(vkUserID, vkUserToken, findLanguage(userLanguage))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := env.firebase.CreateCustomToken(strconv.FormatInt(vkUserID, 10))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(AuthVkResponce{
		Token:    token,
		FullName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Avatar:   user.Photo,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func findLanguage(lang string) vk.Lang {
	switch lang {
	case "ru":
		return vk.LangRU
	case "ua":
		return vk.LangUA
	case "be":
		return vk.LangBE
	case "en":
		return vk.LangEN
	case "es":
		return vk.LangES
	case "fi":
		return vk.LangFI
	case "de":
		return vk.LangDE
	case "it":
		return vk.LangIT
	default:
		return vk.LangEN
	}
}
