package dbmanager

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func newConnect() *sql.DB {

	db, err := sql.Open("mysql", "tfm:tfm@tcp(localhost:3306)/hqa")
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

///
