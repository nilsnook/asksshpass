package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type app struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func newApp(f *os.File) *app {
	ilog := log.New(f, "INFO\t", log.LstdFlags)
	eLog := log.New(f, "ERROR\t", log.LstdFlags|log.Lshortfile)
	return &app{
		infoLog:  ilog,
		errorLog: eLog,
	}
}

// pass gets the specified ssh key password from
// cmd line unix utility called "pass".
// Also there is an assumption that you have
// your ssh key password stored under 'ssh' branch
// in this "pass" utility
func (a *app) pass(key string) string {
	key = fmt.Sprintf("ssh/%s", key)
	passCmd := exec.Command("pass", key)
	passCmd.Stderr = a.errorLog.Writer()
	pass, err := passCmd.Output()
	if err != nil {
		a.errorLog.Fatalln(err)
	}
	return string(pass)
}

func main() {
	LOGFILE := path.Join(os.TempDir(), "asksshpass.log")
	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// create new new app
	a := newApp(f)

	args := os.Args
	if len(args) > 1 {
		a.infoLog.Println(args[1])
		// extract the ssh key from the argument string
		// split by "/" and trim "spaces", ":"
		qarr := strings.Split(args[1], "/")
		key := qarr[len(qarr)-1]
		key = strings.TrimRight(key, ": ")
		a.infoLog.Println("Key:", key)
		// get the password from "pass" cmd utility
		// and print it
		fmt.Print(a.pass(key))
	} else {
		a.errorLog.Fatalln("No args received!")
	}
}
