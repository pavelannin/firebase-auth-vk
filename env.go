package main

import (
	"github.com/anninpavel/firebase-auth-vk/firebase"
	"github.com/anninpavel/firebase-auth-vk/vk"
)

// Environment -- Структура содержит окружение приложения
type Environment struct {
	conf     *Configuration
	vk       *vk.API
	firebase *firebase.Environment
}

// NewEnvironment -- Создает новое окружение
func NewEnvironment() *Environment {
	conf := NewConfiguration()
	return &Environment{
		conf:     conf,
		vk:       vk.New(conf.VkClientID, conf.VkClientSecret),
		firebase: firebase.New(conf.FireBaseCredentialsFile),
	}
}
