package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bomkz/cloudserver/srvauth"
	"github.com/bomkz/cloudsync/cliauth"
)

func checkAuth() {

	appData, err := os.UserConfigDir()

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(appData + "/bomkz")
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(appData+"/bomkz", 0755)
	} else if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(appData + "/bomkz/cloudsync")
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(appData+"/bomkz/cloudsync", 0755)
	} else if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(appData + "/bomkz/cloudsync/token.dat")
	if err != nil && os.IsNotExist(err) {
		err = cliauth.AuthenticateUser()
		if err != nil {
			log.Fatal(err)
		}

		token, err := cliauth.GetAuthToken()

		if err != nil {
			log.Fatal(err)
		}

		os.WriteFile(appData+"/bomkz/cloudsync/token.dat", []byte(token), 0755)
	} else if err != nil {
		log.Fatal(err)
	}

	token, err := os.ReadFile(appData + "/bomkz/cloudsync/token.dat")
	if err != nil {
		log.Fatal(err)
	}

	uuid, err := srvauth.CheckAuth(string(token))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(uuid)
}
