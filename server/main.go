package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	log "github.com/sirupsen/logrus"
)

// ClientCredentials структура для хранения данных клиентов OAuth2
type ClientCredentials struct {
	ID     string
	Secret string
	Domain string
}

// Словарь для хранения учетных данных клиентов
var clients = map[string]ClientCredentials{
	"111111": {"111111", "12345", "http://localhost/"},
	"222222": {"222222", "23456", "http://localhost/"},
	"333333": {"333333", "34567", "http://localhost/"},
}

// Словарь для соответствия идентификатора клиента и его сервиса
var services = map[string]string{
	"111111": "task_manager",
	"222222": "archive_manager",
	"333333": "my_service",
}

// example https://github.com/go-oauth2/oauth2/blob/0572260e96a86bb84c724204c9254d6ae739503d/example/server/server.go#L129
func setupHandlers(srv *server.Server) {
	http.HandleFunc("/get_token", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/check-token", func(w http.ResponseWriter, r *http.Request) {
		token, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := services[token.GetClientID()]
		response := map[string]interface{}{
			"service": data,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

func main() {
	manager := setupOAuthManager()
	srv := setupOAuthServer(manager)
	setupHandlers(srv)

	port := os.Getenv("SERVER_PORT")
	url := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(url, nil))
}

func setupOAuthManager() *manage.Manager {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	for id, creds := range clients {
		clientStore.Set(id, &models.Client{
			ID:     creds.ID,
			Secret: creds.Secret,
			Domain: creds.Domain,
		})
	}

	manager.MapClientStorage(clientStore)
	return manager
}

func setupOAuthServer(manager *manage.Manager) *server.Server {
	srvr := server.NewDefaultServer(manager)
	srvr.SetAllowGetAccessRequest(true)
	srvr.SetClientInfoHandler(server.ClientFormHandler)

	srvr.SetInternalErrorHandler(func(err error) *errors.Response {
		log.Println("Internal Error:", err.Error())
		return nil
	})

	srvr.SetResponseErrorHandler(func(err *errors.Response) {
		log.Println("Response Error:", err.Error.Error())
	})

	return srvr
}
