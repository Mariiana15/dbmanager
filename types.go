package dbmanager

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
