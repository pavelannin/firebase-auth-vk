package vk

const apiVersion = "5.68"

// API --
type API struct {
	AccessToken  string
	ClietnID     string
	ClientSecret string
	debug        bool
}

// New -- Создает экземпляр API
func New(clientID string, clientSecret string) *API {
	return &API{
		ClietnID:     clientID,
		ClientSecret: clientSecret,
		debug:        false,
	}
}

// NewDebug -- Создает экземпляр API, с логирование операций
func NewDebug(clientID string, clientSecret string) *API {
	return &API{
		ClietnID:     clientID,
		ClientSecret: clientSecret,
		debug:        true,
	}
}
