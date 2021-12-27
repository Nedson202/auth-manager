package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bxcodec/faker/v3"
	"github.com/gocarina/gocsv"
)

type Auth struct {
	email    string `faker:"email,unique"`
	password string `faker:"password,unique"`
}

func createFile(authRows []Auth) {
	authDataFilename := "auth-record.csv"
	authDataFile, err := os.Create(authDataFilename)
	if err != nil {
		log.Println(err)
		return
	}

	err = gocsv.MarshalFile(authRows, authDataFile)
}

func generateAuthRecord(recordToGeneratePerThread int) (authRows []Auth) {
	for i := 0; i < recordToGeneratePerThread; i++ { // Generate 100000 structs having a unique word
		auth := Auth{}
		err := faker.FakeData(&auth)
		if err != nil {
			fmt.Println(err)
		}

		log.Println(fmt.Sprintf("Appending %v record", i))
		authRows = append(authRows, auth)
	}
	return authRows
}

func main() {
	authRows := generateAuthRecord(100000)
	createFile(authRows)
}
