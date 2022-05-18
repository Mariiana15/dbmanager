package dbmanager

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
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

type Car struct {
	Id          string `json:"id"`
	Brand       string `json:"brand" validate:"required,alphanum"`
	Model       string `json:"model" validate:"required,alphanum"`
	Horse_power uint32 `json:"horse_power" validate:"required,gte=0,lte=10000"`
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

type General struct {
	Gid  string `json:"gid"`
	Name string `json:"name"`
}

type Story struct {
	Gid  string `json:"gid"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type CustomField struct {
	Gid   string `json:"gid"`
	Name  string `json:"name"`
	Value string `json:"display_value"`
}

type Task struct {
	Id           string        `json:"id"`
	Hid          string        `json:"hid"`
	Gid          string        `json:"gid"`
	UserId       string        `json:"userId"`
	Name         string        `json:"name"`
	Notes        string        `json:"notes"`
	CustomField  []CustomField `json:"custom_fields"`
	Link         string        `json:"permalink_url"`
	Story        []Story       `json:"stories"`
	Dependecies  []General     `json:"dependencies"`
	State        string        `json:"state"`
	TypeTest     string        `json:"typeTest"`
	TypeTestId   string        `json:"typeTestId"`
	TypeUS       string        `json:"typeUS"`
	UserStory    string        `json:"userStory"`
	Priority     int           `json:"priority"`
	Alerts       int           `json:"alerts"`
	Scripts      int           `json:"scripts"`
	Date         int64         `json:"date"`
	UrlAlert     string        `json:"urlAlert"`
	UrlScript    string        `json:"urlScript"`
	AddInfo      int8          `json:"addInfo"`
	Test         General       `json:"test"`
	Result       Result        `json:"result"`
	Tecnologies  string        `json:"technologies"`
	Requirement  string        `json:"requirement"`
	Architecture string        `json:"architecture"`
}

type Result struct {
	Message   string `json:"message"`
	Alert     int    `json:"alert"`
	UrlAlert  string `json:"urlAlert"`
	Detail    string `json:"detail"`
	Script    int    `json:"script"`
	UrlScript string `json:"urlScript"`
}

type Section struct {
	Name      string  `json:"name"`
	ID        string  `json:"id"`
	Gid       string  `json:"gid"`
	Project   General `json:"project"`
	StoryUser []Task  `json:"storyUser"`
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

/// past

func (car *Car) insert() error {
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

func (car *Car) get(id string) error {
	var db = newConnect()
	err := db.QueryRow("SELECT id,brand,model,horse_power FROM cars WHERE id = ?", id).Scan(&car.Id, &car.Brand, &car.Model, &car.Horse_power)
	if err != nil {
		return err
	}
	defer db.Close()
	return err
}

func (car *Car) delete() error {
	var db = newConnect()
	_, err := db.Query(fmt.Sprintf("DELETE FROM cars WHERE id = '%s'", car.Id))
	if err != nil {
		return err
	}
	defer db.Close()
	return err
}

///

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

func (section *Section) SetSectionProject(Uid string) error {

	var db = newConnect()
	atCreate := time.Now().Unix()
	_, err := db.Query(fmt.Sprintf("INSERT section_project VALUES ('%v','%v','%v','%v','%v','%v', %v, '%v')", section.ID, Uid, section.Name, section.Gid, section.Project.Gid, section.Project.Name, atCreate, "active"))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (task *Task) SetUserStoryAsana(secId string) error {

	var db = newConnect()
	task.Id = uuid.NewV4().String()
	_, err := db.Query(fmt.Sprintf("INSERT user_story_asana VALUES ('%v','%v','%v','%v','%v','%v')", task.Id, task.Gid, secId, task.Name, task.Notes, task.Link))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (dependecie *General) SetUserStoryAsanaDependence(idTask string) error {

	var db = newConnect()
	id := uuid.NewV4().String()
	_, err := db.Query(fmt.Sprintf("INSERT user_story_asana_dependence VALUES ('%v','%v','%v','%v')", id, idTask, dependecie.Gid, dependecie.Name))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (task *Task) SetUserStoryAsanaDependence() error {

	for i := 0; i <= len(task.Dependecies)-1; i++ {
		err := task.Dependecies[i].SetUserStoryAsanaDependence(task.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cfield *CustomField) SetUserStoryAsanaCField(idTask string) error {

	var db = newConnect()
	id := uuid.NewV4().String()
	_, err := db.Query(fmt.Sprintf("INSERT user_story_asana_cfield VALUES ('%v','%v','%v','%v','%v')", idTask, id, cfield.Gid, cfield.Name, cfield.Value))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (task *Task) SetUserStoryAsanaCField() error {

	for i := 0; i <= len(task.CustomField)-1; i++ {
		err := task.CustomField[i].SetUserStoryAsanaCField(task.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (story *Story) SetUserStoryAsanaStories(idTask string) error {

	var db = newConnect()
	id := uuid.NewV4().String()
	_, err := db.Query(fmt.Sprintf("INSERT user_story_asana_stories VALUES ('%v','%v','%v','%v', '%v')", idTask, id, story.Gid, story.Text, story.Type))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (task *Task) SetUserStoryAsanaStories() error {

	for i := 0; i <= len(task.Story)-1; i++ {
		err := task.Story[i].SetUserStoryAsanaStories(task.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

/////

func (task *Task) SetUserStory() error {

	var db = newConnect()
	task.Hid = uuid.NewV4().String()
	_, err := db.Query(fmt.Sprintf("INSERT user_story VALUES ('%v',%v,'%v','%v',%v, %v, %v,%v,'%v','%v','%v', '%v','%v', '%v','%v','%v','%v','%v')", task.Hid, task.AddInfo, "1", task.Id, task.Date, task.Priority, task.Scripts, task.Alerts, task.State, task.TypeTestId, task.TypeUS, task.UrlAlert, task.UrlScript, task.UserStory, task.Tecnologies, task.Architecture, task.Requirement, task.UserId))
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	return nil
}

func (task *Task) SetUserStoryResult() error {

	var db = newConnect()
	id := ""
	response, err := db.Query(fmt.Sprintf("SELECT id FROM user_story_result where hid= '%v'", task.Hid))
	if err != nil {
		return err
	}
	for response.Next() {
		response.Scan(&id)
	}
	if id == "" {
		id = uuid.NewV4().String()
		_, err := db.Query(fmt.Sprintf("INSERT user_story_result VALUES ('%v','%v',%v,'%v','%v', %v,'%v','%v')", task.Hid, id, task.Result.Alert, task.Result.Detail, task.Result.Message, task.Result.Script, task.Result.UrlAlert, task.Result.UrlScript))
		if err != nil {
			return err
		}
	} else {
		_, err := db.Query(fmt.Sprintf("UPDATE user_story_result SET alert='%v', detail='%v', message='%v', script='%v', urlAlert='%v', urlScript='%v' where id='%v'", task.Result.Alert, task.Result.Detail, task.Result.Message, task.Result.Script, task.Result.UrlAlert, task.Result.UrlScript, id))
		if err != nil {
			return err
		}
	}
	_, err = db.Query(fmt.Sprintf("UPDATE user_story SET state = '%v' WHERE id = '%v';", "close", task.Hid))
	if err != nil {

		return err
	}
	defer db.Close()
	return nil
}

///

func GetTestHQA() (General, error) {

	var db = newConnect()
	var test General
	response, err := db.Query("SELECT id, name from test_hqa where id = '39c6f825-d086-11ec-ac5b-48e244d50ee5'")
	if err != nil {
		return test, err
	}
	for response.Next() {
		response.Scan(&test.Gid, &test.Name)
	}
	defer db.Close()
	return test, nil
}

func GetUserStoriesComplete(s *[]Section, userId string) error {

	err := GetSectionDB(userId, s)
	if err != nil {
		return err
	}
	return nil

}

func GetSectionDB(userId string, sections *[]Section) error {

	var db = newConnect()
	response, err := db.Query(fmt.Sprintf("SELECT id, name, sectionId, projectId, nameProject FROM section_project where state = '%v' AND userId = '%v';", "active", userId))
	if err != nil {
		return err
	}
	var sec Section
	for response.Next() {
		response.Scan(&sec.ID, &sec.Name, &sec.Gid, &sec.Project.Gid, &sec.Project.Name)
		var task []Task
		GetTaskDB(sec, &task)
		sec.StoryUser = task
		*sections = append(*sections, sec)
	}
	defer db.Close()
	return nil
}

func GetTaskDB(sec Section, task *[]Task) error {

	var db = newConnect()
	response, err := db.Query(fmt.Sprintf("SELECT id, gid, name, notes, permalink_url FROM user_story_asana where  sectionId = '%v';", sec.ID))
	if err != nil {
		return err
	}
	var t Task
	for response.Next() {
		response.Scan(&t.Id, &t.Gid, &t.Name, &t.Notes, &t.Link)

		var d []General
		var cf []CustomField
		var st []Story

		GetDependence(t, &d)
		GetCField(t, &cf)
		GetStories(t, &st)
		t.Dependecies = d
		t.CustomField = cf
		t.Story = st
		GetUserStoryFromAsana(&t, task)
	}
	defer db.Close()
	return nil
}

func GetDependence(task Task, dep *[]General) error {

	var db = newConnect()
	response, err := db.Query(fmt.Sprintf("SELECT gid, type FROM user_story_asana_dependence where usId = '%v';", task.Id))
	if err != nil {
		return err
	}
	var d General
	for response.Next() {
		response.Scan(&d.Gid, &d.Name)
		*dep = append(*dep, d)
	}
	defer db.Close()
	return nil
}

func GetCField(task Task, cfs *[]CustomField) error {

	var db = newConnect()
	response, err := db.Query(fmt.Sprintf("SELECT gid, name, display_value FROM user_story_asana_cfield where usId = '%v';", task.Id))
	if err != nil {
		return err
	}
	var cf CustomField
	for response.Next() {
		response.Scan(&cf.Gid, &cf.Name, &cf.Value)
		*cfs = append(*cfs, cf)
	}
	defer db.Close()
	return nil
}

func GetStories(task Task, sts *[]Story) error {

	var db = newConnect()
	response, err := db.Query(fmt.Sprintf("SELECT gid, text, type FROM user_story_asana_stories where usId = '%v';", task.Id))
	if err != nil {
		return err
	}
	var s Story
	for response.Next() {
		response.Scan(&s.Gid, &s.Text, &s.Type)
		*sts = append(*sts, s)
	}
	defer db.Close()
	return nil
}

func GetUserStoryFromAsana(t *Task, tasks *[]Task) error {

	var db = newConnect()
	response, err := db.Query(fmt.Sprintf("SELECT id, addInfo, date,priority, scripts, alerts,state, idTest, typeUs, urlAlert, urlScript, userStory, technologies, architecture, requirement FROM user_story where idIntegration = '%v' and state in ('%v','%v') and usId ='%v';", "1", "open", "close", t.Id))
	if err != nil {
		return err
	}
	for response.Next() {
		response.Scan(&t.Hid, &t.AddInfo, &t.Date, &t.Priority, &t.Scripts, &t.Alerts, &t.State, &t.TypeTestId, &t.TypeUS, &t.UrlAlert, &t.UrlScript, &t.UserStory, &t.Tecnologies, &t.Architecture, &t.Requirement)
		GetUserStoryResult(t)
		*tasks = append(*tasks, *t)
	}
	defer db.Close()
	return nil
}

func GetUserStoryResult(t *Task) error {

	var db = newConnect()
	var r Result
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

//////

func SetInfoTech(t string, a string, r string, id string) error {
	var db = newConnect()
	_, err := db.Query(fmt.Sprintf("UPDATE user_story SET technologies = '%v', architecture = '%v', requirement = '%v', addInfo = %v where id = '%v'", t, a, r, 0, id))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func SetChangeStateUserStory(c string, id string) error {
	var db = newConnect()
	_, err := db.Query(fmt.Sprintf("UPDATE user_story SET state = '%v' where id = '%v'", c, id))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func SetChangeStateSection(c string, id string) error {
	var db = newConnect()

	_, err := db.Query(fmt.Sprintf("UPDATE section_project SET state = '%v' where sectionId = '%v' and state='%v'", c, id, "active"))
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}
