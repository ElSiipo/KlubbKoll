package app

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
)

// TODO:
// Get all clubs.
// Get single club.
// Upvote/downvote? Or just "update"? Maybe update is dangerous.

// GetClubs returns all clubs
func GetClubs(c *gin.Context) {
	connectionString := GetConnectionStringFromFile()

	session, err := mgo.Dial(connectionString)
	if err != nil {
		panic(err)
	}
	defer session.Close()

}

// GetConnectionStringFromFile Gets connection string from file
func GetConnectionStringFromFile() string {
	absPath, _ := filepath.Abs("../rest-api-server/app/connectionString.txt")

	byteArray, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(byteArray[:])
}
