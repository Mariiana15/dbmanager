package dbmanager

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

type ObjectGenericDB struct {
	Id string `json:"id"`
}

type dbManager interface {
	insert() error
	get(id string) error
	delete() error
}

type Token struct {
	UserId   string `json:"userId"`
	Token    string `json:"token"`
	Id       string `json:"id"`
	ExpireAt int64  `json:"expireAt"`
	State    string `json:"state"`
}

type Auth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Id   string `json:"id"`
}

func newConnect() *sql.DB {
	db, err := sql.Open("mysql", "tfm:tfm@tcp(localhost:3306)/hqa")
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}
func insertDB(db dbManager) error {
	err := db.insert()
	return err
}

func getDB(db dbManager, id string) error {
	err := db.get(id)
	return err
}

func deleteDB(db dbManager) error {
	err := db.delete()
	return err
}

func NewToken() interfaceToken {
	return &Token{}
}

type interfaceToken interface {
	Insert(access bool, source bool) error
	GetToken(access bool) error
	DeleteToken(access bool) error
}

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

//////////////////////

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

///////TFM

func SetUserCodeAsana(userId string, codeVerifier string, code string) error {

	var db = newConnect()
	id := uuid.NewV4().String()
	_, err := db.Query(fmt.Sprintf("INSERT user_code_asana VALUES ('%v','%v','%v','%v', '%v')", id, userId, codeVerifier, code, time.Now().Unix()))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func GetUserCodeAsana(userId string) (string, string, error) {

	var db = newConnect()
	var code string
	var code_verifier string

	response, err := db.Query(fmt.Sprintf("SELECT code_verifier, code FROM user_code_asana where userId= '%v'", userId))
	if err != nil {
		return "", "", err
	}
	for response.Next() {
		response.Scan(&code_verifier, &code)
	}
	defer db.Close()
	return code, code_verifier, nil
}

/////
