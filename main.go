package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gcfg.v1"
	"io/ioutil"
	"log"
	"o365-attack-toolkit/api"
	"o365-attack-toolkit/model"
	"o365-attack-toolkit/server"
	"os"
	"path/filepath"
	_ "database/sql"
	_"github.com/mattn/go-sqlite3"
)

func main(){

	model.GlbConfig = model.Config{}
	err := gcfg.ReadFileInto(&model.GlbConfig,"template.conf")
  var userToken string = ""


	if err != nil {
		log.Fatal(err.Error())
	}

  argsWithoutProg := os.Args[1:]
  if len(argsWithoutProg) > 0 {
    log.Println("testmode")
    userToken = argsWithoutProg[0]
    log.Println(userToken)
    //api.GetADUsers(userToken)
    go server.StartExtServer(model.GlbConfig)
    server.StartIntServer(model.GlbConfig)
    fmt.Println(model.GlbConfig)
  }  else  {
    initializeRules()
    go server.StartExtServer(model.GlbConfig)
    server.StartIntServer(model.GlbConfig)
    fmt.Println(model.GlbConfig)
  }
}


func initializeRules(){

	var ruleFiles []string
	var tempRule model.Rule

	root := "rules"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		ruleFiles = append(ruleFiles, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range ruleFiles {

		ruleFile, err := os.Open(file)

		if err != nil {
			log.Println(err)
		}

		defer ruleFile.Close()

		byteValue, _ := ioutil.ReadAll(ruleFile)

		json.Unmarshal(byteValue,&tempRule)

		model.GlbRules = append(model.GlbRules,tempRule)

	}

	log.Printf("Loaded %d rules successfully.",len(model.GlbRules))
}
