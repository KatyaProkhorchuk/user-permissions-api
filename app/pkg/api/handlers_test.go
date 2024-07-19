package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"user-permissions-api/pkg/db"
	"user-permissions-api/pkg/types"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestInsertHandler(t *testing.T) {
	//connect to db
	log := logrus.New()
	log.Out = os.Stdout
	db.ConnectDB(log)
	defer db.DataBase.Close(db.Ctx)
	db.CreateTable()
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/insertUser", InsertHandler(log))

	user := types.User{Name: "testuser", Access: types.AccessServices{}}
	jsonData, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPost, "/insertUser", bytes.NewReader(jsonData))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteHandler(t *testing.T) {
	log := logrus.New()
	log.Out = os.Stdout
	db.ConnectDB(log)
	defer db.DataBase.Close(db.Ctx)
	db.CreateTable()
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/deleteUser", DeleteHandler(log))

	user := types.User{Name: "testuser", Access: types.AccessServices{}}
	jsonData, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPost, "/deleteUser", bytes.NewReader(jsonData))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
