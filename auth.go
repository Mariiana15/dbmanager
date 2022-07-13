package dbmanager

import (
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func (token *Token) Insert(access bool, source bool) error {

	var db = newConnect()
	var err error
	token.Id = uuid.NewV4().String()
	token.State = "active"
	table := "user_token"
	if !access {
		table = "user_refresh_token"
	}
	err = setTableToken(db, *token, table, source)
	if err != nil {
		return err
	}
	defer db.Close()
	return err
}
func setTableToken(db *sql.DB, token Token, table string, source bool) error {

	var o ObjectGenericDB

	response, err := db.Query(fmt.Sprintf("SELECT id FROM %v where state in ('%v', '%v') AND userId = '%v';", table, token.State, "refresh", token.UserId))
	if err != nil {
		return err
	}
	for response.Next() {
		response.Scan(&o.Id)
		if o.Id != "" {
			_, err = db.Query(fmt.Sprintf("UPDATE  %v SET state = '%v' WHERE id = '%v';", table, "inactive", o.Id))
			if err != nil {
				return err
			}
		}
	}
	if source {
		token.State = "refresh"
	}
	_, err = db.Query(fmt.Sprintf("INSERT INTO %v VALUES ( '%v', '%v', '%v',%v, '%v' );", table, token.Id, token.UserId, token.Token, token.ExpireAt, token.State))
	if err != nil {
		return err
	}
	return err
}

func (token *Token) GetToken(access bool) error {

	var db = newConnect()
	table := "user_token"
	if !access {
		table = "user_refresh_token"
	}
	response, err := db.Query(fmt.Sprintf("SELECT id,token FROM %v where state = '%v' AND userId = '%v';", table, "active", token.UserId))
	if err != nil {
		return err
	}
	for response.Next() {
		response.Scan(&token.Id, &token.Token)
	}
	if token.Id == "" {
		return fmt.Errorf("token no found")
	}
	defer db.Close()
	return nil
}

func (token *Token) DeleteToken(access bool) error {

	var db = newConnect()
	table := "user_token"
	if !access {
		table = "user_refresh_token"
	}
	_, err := db.Query(fmt.Sprintf("UPDATE  %v SET state = '%v' WHERE id = '%v';", table, "inactive", token.Id))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (auth *Auth) GetUserBasic(id string, email string) error {

	var db = newConnect()
	response, err := db.Query(fmt.Sprintf("SELECT id FROM user WHERE idGoogle = '%v' and email ='%v'", id, email))
	if err != nil {
		fmt.Println(err)
		return err
	}
	for response.Next() {
		response.Scan(&auth.Id)
	}
	response, err = db.Query(fmt.Sprintf("SELECT user, pass FROM user_basic WHERE userId = '%v'", auth.Id))
	if err != nil {
		fmt.Println(err)
		return err
	}
	for response.Next() {
		response.Scan(&auth.User, &auth.Pass)
	}
	defer db.Close()
	return err

}
