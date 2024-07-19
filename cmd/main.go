package app

import (
	"fmt"
	"os"
	"sync"

	"github.com/KatyaProkhorchuk/user-permissions-api/pkg/api"
	"github.com/gin-gonic/gin" // api
	"github.com/sirupsen/logrus"
)

var (
	wg      sync.WaitGroup
	log     = logrus.New()
	logFile = "logs/server_log.log"
)

// (создание/удаление пользователя, добавить/убрать права пользователя, проверка прав пользователя).
func adminRouters() {
	admin := gin.Default()
	admin.POST("/insertUser", api.InsertHandler(log))
}

func main() {
	// выводим логи по умолчанию в stdout вместо stderr
	log.Out = os.Stdout
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0755)
	// eсли файл есть то логи пишем в него
	if err != nil {
		log.Warn(fmt.Sprintf("Could not open log file: %v\n", err))
	} else {
		defer file.Close()
		log.Out = file
	}
	// бд
}
