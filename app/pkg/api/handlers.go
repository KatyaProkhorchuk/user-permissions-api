package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"user-permissions-api/pkg/db"
	"user-permissions-api/pkg/types"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var mu sync.Mutex

func handleError(ctx *gin.Context, statusCode int, message string, err error, log *log.Logger) {
	ctx.AbortWithStatusJSON(statusCode, gin.H{"error": message})
	log.Warning(fmt.Sprintf("%s: %v\n", message, err))
}
func getUserData(ctx *gin.Context, log *log.Logger) (*types.User, error) {
	var user types.User
	// Attempt to bind the JSON body of the request to the user struct
	if err := ctx.ShouldBindJSON(&user); err != nil {
		message := "Bad input"
		handleError(ctx, http.StatusBadRequest, message, err, log)
		return nil, err
	}
	return &user, nil
}
func getClientData(ctx *gin.Context, log *log.Logger) (*types.Request, error) {
	var request types.Request
	// Attempt to bind the JSON body of the request to the user struct
	if err := ctx.ShouldBindJSON(&request); err != nil {
		message := "Bad input"
		handleError(ctx, http.StatusBadRequest, message, err, log)
		return nil, err
	}
	return &request, nil
}

func InsertHandler(log *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := getUserData(ctx, log)
		if err != nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if err := db.InsertUser(user); err != nil {
			message := fmt.Sprintf("Error add user %v\n", err)
			log.WithFields(logrus.Fields{"error": err.Error()}).Error(message)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}
		messageOk := "User successfuly created"
		ctx.JSON(http.StatusOK, messageOk)
		log.Info(messageOk)
	}
}

func DeleteHandler(log *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := getUserData(ctx, log)
		if err != nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if err := db.DeleteUser(user); err != nil {
			message := fmt.Sprintf("Error delete user %v\n", err)
			log.WithFields(logrus.Fields{"error": err.Error()}).Error(message)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}
		messageOk := "User successfuly deleted"
		ctx.JSON(http.StatusOK, messageOk)
		log.Info(messageOk)
	}
}

func AddUserRightsHandler(log *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := getUserData(ctx, log)
		if err != nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if err := db.UpdateUserRights(user); err != nil {
			message := fmt.Sprintf("Error update user %v\n", err)
			log.WithFields(logrus.Fields{"error": err.Error()}).Error(message)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}
		messageOk := "User successfuly update"
		ctx.JSON(http.StatusOK, messageOk)
		log.Info(messageOk)
	}
}

func DeleteUserRightsHandler(log *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := getUserData(ctx, log)
		if err != nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if err := db.DeleteUserRights(user); err != nil {
			message := fmt.Sprintf("Error update user %v\n", err)
			log.WithFields(logrus.Fields{"error": err.Error()}).Error(message)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}
		messageOk := "User successfuly update"
		ctx.JSON(http.StatusOK, messageOk)
		log.Info(messageOk)
	}
}

func checkToken(token string, log *log.Logger) (service string, err error) {
	host := "server"
	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	clientId := "333333"
	secret := "34567"
	// https://www.nic.ru/help/oauth-server_3642.html#token
	url := fmt.Sprintf(
		"http://%s:%d/check-token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		host,
		port,
		clientId,
		secret,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Error creating HTTP request: %v", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error executing HTTP request: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Errorf("HTTP request returned non-200 status: %d, body: %s", resp.StatusCode, body)
		return "", fmt.Errorf("HTTP request returned non-200 status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %v", err)
		return "", err
	}
	var responseStruct map[string]string
	if err = json.Unmarshal(body, &responseStruct); err != nil {
		log.Errorf("Error unmarshalling JSON response: %v", err)
		return "", err
	}
	service, ok := responseStruct["service"]
	if !ok {
		return "", fmt.Errorf("service field missing in response")
	}
	return service, nil
}

func CheckAccessHandler(log *log.Logger) gin.HandlerFunc {
	// проверка прав пользователя
	return func(ctx *gin.Context) {
		req, err := getClientData(ctx, log)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request information"})
			return
		}
		log.Warning(fmt.Sprintf("Token: %s\n", req.Token))
		name, err := checkToken(req.Token, log)
		if err != nil {
			handleError(ctx, http.StatusBadRequest, "Access Denied, by token", err, log)
			return
		}
		user, err := db.GetUserByName(req.Name)
		if err != nil {
			handleError(ctx, http.StatusBadRequest, "Could not check user", err, log)
			return
		}

		var access []string

		// Определение, какой менеджер доступа использовать на основе имени сервиса
		switch name {
		case "archive_manager":
			access = user.Access.Archive.Records
		case "task_manager":
			access = user.Access.Task.Agent
		default:
			fmt.Errorf("unknown service: %s", name)
		}

		if len(access) == 0 {
			fmt.Errorf("no access to service %s", name)
			return
		}

		if err != nil {
			handleError(ctx, http.StatusBadRequest, "Error retrieving user access", err, log)
			return
		}
		responce := struct {
			Access []string `json:"access"`
		}{
			Access: access,
		}
		log.Info("User is successfully checked.")
		ctx.JSON(http.StatusOK, responce)

	}
}
