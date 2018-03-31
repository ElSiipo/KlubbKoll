package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	goji "goji.io"
	mgo "gopkg.in/mgo.v2"
	// httpDeliver "github.com/ElSiipo/klubbkoll/club/delivery/http"
	// clubRepo "github.com/ElSiipo/klubbkoll/club/repository/mongoDB"
	// articleUcase "github.com/ElSiipo/klubbkoll/club/usecase"
)

func main() {
	session, err := mgo.Dial(getConnectionStringFromFile())

	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	mux := goji.NewMux()
	// mux.HandleFunc(pat.Get("/clubs"), mongodb_club.GetAll(session))
	// mux.HandleFunc(pat.Post("/clubs"), Store(session))
	// mux.HandleFunc(pat.Get("/clubs/:club_id"), GetByID(session))
	// mux.HandleFunc(pat.Put("/clubs/:club_id"), Update(session))
	// mux.HandleFunc(pat.Delete("/clubs/:club_id"), Delete(session))

	http.ListenAndServe("localhost:1234", mux)
}

func getConnectionStringFromFile() string {
	absPath, _ := filepath.Abs("./connectionString.txt")

	byteArray, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(byteArray[:])
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()
	c := session.DB("klubbkoll").C("clubs")
	index := mgo.Index{
		Key:        []string{"club_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
