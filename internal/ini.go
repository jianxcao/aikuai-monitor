package aikuaimonitor

import (
	"log"
	"os"
	"path"

	"github.com/go-ini/ini"
)

type User struct {
	Name       string
	Password   string
	Url        string
	ClientName string
}

func InitIni() ([]*User, error) {
	dir, _ := os.Getwd()
	log.Println("current path", dir)
	cfg, err := ini.Load(path.Join(dir, "config/conf.ini"))
	userPath := os.Getenv("AIKUAI_MONITOR_CONFIG_PATH")
	if len(userPath) > 0 {
		cfg, err = ini.Load(path.Join(userPath, "conf.ini"))
	}

	if err != nil {
		log.Fatal(err)
		log.Fatal("请确认在app下有/config/conf.ini文件")
		return nil, err
	}

	Sections := cfg.Sections()
	var users = []*User{}
	for _, s := range Sections {
		user := s.Key("user").String()
		password := s.Key("password").String()
		url := s.Key("url").String()

		if len(user) != 0 && len(password) != 0 && len(url) != 0 {
			users = append(users, &User{
				Name:       user,
				Password:   password,
				Url:        url,
				ClientName: s.Name(),
			})
		}
	}
	return users, err
}
