package author

import "gorm.io/gorm"

var App *Application

type Application struct {
	MySQL *gorm.DB
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	MySQL, err := InitMySQL()
	if err != nil {
		panic(err.Error())
	}
	App.MySQL = MySQL
}
