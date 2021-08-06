package conf

import (
	"bytes"
	"cmdtools/tools"
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	Env        = "dev"
	ExecDir    = "" // current exec file path
	IntranetIp = ""
)

//init config
func InitConf() {
	// init Intranet Ip
	if runtime.GOOS == "windows" {
		// todo windows intranet ip get optimize
		IntranetIp = "127.0.0.1"
	} else {
		IntranetIp, _ = tools.GetClientIp()
	}
	log.Println("intranet ip is " + IntranetIp)
	dir, err := GetCurrentPath()
	ExecDir = dir
	if err != nil {
		log.Fatalln(err)
	}

	// read DRAGON env first, if empty str them run dev env
	env := os.Getenv("DRAGON")
	log.Println("os.Getenv:", env)
	if env == "" {
		env = "dev"
	}
	Env = env

	var config []byte
	pathFile := dir + "../core/dragon/conf/config/" + Env + ".yml"
	config, _ = ioutil.ReadFile(pathFile)
	err = viper.ReadConfig(bytes.NewReader(config))
	if err != nil {
		// read yml config fail, return fail
		log.Fatalln("init config fail: core/dragon/conf/config/"+Env+".yml not found", err)
		return
	}
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	err = viper.ReadConfig(bytes.NewReader(config))
	if err != nil {
		log.Fatalln("viper.ReadConfig fail", err)
	}
}

//get current exec file path
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return path[0 : i+1], nil
}

// according to operating system to change path slash, default use linux path slash
func FmtSlash(path string) string {
	sys := runtime.GOOS
	if sys == `windows` {
		return strings.Replace(path, "/", "\\", -1)
	}
	return path
}
