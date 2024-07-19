package db

import (
	"errors"
	"fmt"
	"user-permissions-api/pkg/types"

	// "user-permissions-api/pkg/api"

	"github.com/jackc/pgx/v5"
)

const (
	createTable = `CREATE TABLE Users (
	name text not null,
	user_acsess jsonb)`
	insertUser       = `INSERT INTO Users(name, user_acsess) VALUES($1, $2)`
	deleteUser       = `DELETE FROM Users WHERE name = $1`
	updateUser       = `UPDATE Users SET user_acsess=$1 WHERE name = $2`
	selectUserByName = `SELECT name, user_acsess FROM Users WHERE name = $1`
)

func CreateTable() error {
	_, err := DataBase.Exec(Ctx, createTable)
	return err
}
func IsExistsUser(name string) (bool, error) {
	var tmp int64
	err := DataBase.QueryRow(Ctx, selectUserByName, name).Scan(&tmp)
	// fmt.Println(tmp)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func InsertUser(user *types.User) error {
	tx, err := DataBase.Begin(Ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction : %w\n", err)
	}
	//Ensure the transaction is either committed or rolled back
	defer func() {
		if err != nil {
			tx.Rollback(Ctx)
		} else {
			tx.Commit(Ctx)
		}
	}()
	ok, err := IsExistsUser(user.Name)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if ok {
		return errors.New("User already exists")
	}
	// Insert the user into the database
	_, err = tx.Exec(Ctx, insertUser, &user.Name, &user.Access)
	if err != nil {
		return fmt.Errorf("error inserting user: %w\n", err)
	}
	return nil
}

func DeleteUser(user *types.User) error {
	_, err := DataBase.Exec(Ctx, deleteUser, user.Name)
	if err != nil {
		return fmt.Errorf("error deleting user %s: %w", user.Name, err)
	}
	return nil
}

func GetUserByName(name string) (*types.User, error) {
	var user types.User
	err := DataBase.QueryRow(Ctx, selectUserByName, name).Scan(&user.Name, &user.Access)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Если пользователь не найден, возвращаем ошибку
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
func merge(old_data, new_data []string) []string {
	result := make([]string, 0)
	unique := make(map[string]bool)
	for _, value := range old_data {
		if _, ok := unique[value]; !ok {
			result = append(result, value)
			unique[value] = true
		}
	}
	for _, value := range new_data {
		if _, ok := unique[value]; !ok {
			result = append(result, value)
			unique[value] = true
		}
	}
	return result
}
func AddUpdateUser(user, newUser *types.User) {
	newUser.Access.Archive.Records = merge(newUser.Access.Archive.Records, user.Access.Archive.Records)
	newUser.Access.Task.Agent = merge(newUser.Access.Task.Agent, user.Access.Task.Agent)
}
func UpdateUserRights(user *types.User) error {
	nameDb, err := GetUserByName(user.Name)
	if err != nil {
		return err
	}

	tx, err := DataBase.Begin(Ctx)
	if err != nil {
		return err
	}

	// Обновляем доступы пользователя согласно переданным данным
	AddUpdateUser(user, nameDb)

	_, err = tx.Exec(Ctx, updateUser, &nameDb.Access, &nameDb.Name)
	if err != nil {
		// Откатываем транзакцию в случае ошибки
		tx.Rollback(Ctx)
		return err
	}

	if err := tx.Commit(Ctx); err != nil {
		return err
	}

	return nil
}
func difference(old_data, new_data []string) []string {
	result := make([]string, 0)
	oldMap := make(map[string]bool, 0)
	newMap := make(map[string]bool, 0)
	for _, value := range new_data {
		newMap[value] = true
	}

	for _, value := range old_data {
		if _, ok := newMap[value]; !ok {
			if _, ok1 := oldMap[value]; !ok1 { // что бы не было дубликатов
				oldMap[value] = true
				result = append(result, value)
			}
		}
	}
	return result
}
func DeleteUpdateUser(user, newUser *types.User) {
	newUser.Access.Archive.Records = difference(newUser.Access.Archive.Records, user.Access.Archive.Records)
	newUser.Access.Task.Agent = difference(newUser.Access.Task.Agent, user.Access.Task.Agent)
}
func DeleteUserRights(user *types.User) error {
	nameDb, err := GetUserByName(user.Name)
	if err != nil {
		return err
	}

	tx, err := DataBase.Begin(Ctx)
	if err != nil {
		return err
	}

	DeleteUpdateUser(user, nameDb)

	_, err = tx.Exec(Ctx, updateUser, &nameDb.Access, &nameDb.Name)
	if err != nil {
		// Откатываем транзакцию в случае ошибки
		tx.Rollback(Ctx)
		return err
	}

	if err := tx.Commit(Ctx); err != nil {
		return err
	}

	return nil
}
