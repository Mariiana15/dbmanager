package dbmanager

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func newConnect() *sql.DB {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	db, err := sql.Open("mysql", "tfm:"+os.Getenv("CONNECTION_STRING"))
	if err != nil {
		fmt.Println(err.Error())
	}
	return db

}
func InsertDB(db dbManager) error {

	err := db.Insert()
	return err
}

func GetDB(db dbManager, id string) error {
	err := db.Get(id)
	return err
}

func DeleteDB(db dbManager) error {
	err := db.Delete()
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

/// past

func (car *Car) Insert() error {
	var db = newConnect()
	car.Id = car.Brand[1:3] + strconv.Itoa(rand.Intn(1000)) + car.Model[1:3]
	_, err := db.Query(fmt.Sprintf("INSERT INTO cars VALUES ( '%s' ,'%s','%s',%d );", car.Id, car.Brand, car.Model, car.Horse_power))
	if err != nil {
		_, err = db.Query(fmt.Sprintf("UPDATE cars SET brand = '%s' , model = '%s', horse_power = %d  WHERE  id = '%s';", car.Brand, car.Model, car.Horse_power, car.Id))
		if err != nil {
			return err
		}
	}
	defer db.Close()
	return err
}

func (car *Car) Get(id string) error {
	var db = newConnect()
	err := db.QueryRow("SELECT id,brand,model,horse_power FROM cars WHERE id = ?", id).Scan(&car.Id, &car.Brand, &car.Model, &car.Horse_power)
	if err != nil {
		return err
	}
	defer db.Close()
	return err
}

func (car *Car) Delete() error {
	var db = newConnect()
	_, err := db.Query(fmt.Sprintf("DELETE FROM cars WHERE id = '%s'", car.Id))
	if err != nil {
		return err
	}
	defer db.Close()
	return err
}

func GetIndustryHQA(s string) error {

	var db = newConnect()
	var r Industry
	t.Result = r
	response, err := db.Query(fmt.Sprintf("SELECT alert, detail, message, script, urlAlert, urlScript FROM user_story_result where hid= '%v'", t.Hid))
	if err != nil {
		return err
	}
	for response.Next() {
		response.Scan(&t.Result.Alert, &t.Result.Detail, &t.Result.Message, &t.Result.Script, &t.Result.UrlAlert, &t.Result.UrlScript)
	}
	defer db.Close()
	return nil

}

///
