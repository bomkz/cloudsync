package main

import (
	"log"
	"os"

	"github.com/bomkz/cloudsync/cliauth"
)

func checkAuth() {

	appData, err := os.UserConfigDir()

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(appData + "bomkz")
	if err != nil && err == os.ErrNotExist {
		os.Mkdir(appData+"/bomkz", 0755)
	}

	err = nil

	_, err = os.Stat(appData + "/bomkz/cloudsync")
	if err != nil && err == os.ErrNotExist {
		os.Mkdir(appData+"/bomkz/cloudsync", 0755)
	}

	err = nil

	_, err = os.Stat(appData + "/bomkz/cloudsync/token.dat")
	if err != nil && err == os.ErrNotExist {

		err = nil
		err = cliauth.AuthenticateUser()
		if err != nil {
			log.Fatal(err)
		}

		token, err := cliauth.GetAuthToken()

		if err != nil {
			log.Fatal(err)
		}

		os.WriteFile(appData+"/bomkz/cloudsync/token.dat", []byte(token), 0755)
	}

	os.ReadFile(appData + "/bomkz/cloudsync/token.dat")

}
