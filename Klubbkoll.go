package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	goji "goji.io"
	"goji.io/pat"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ErrorWithJSON handles errorResponse
func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

// ResponseWithJSON handles regular responses
func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func main() {
	session, err := mgo.Dial(getConnectionStringFromFile())

	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/clubs"), getAllClubs(session))
	mux.HandleFunc(pat.Post("/clubs"), addClub(session))
	mux.HandleFunc(pat.Get("/clubs/:club_id"), clubByClubID(session))
	mux.HandleFunc(pat.Put("/clubs/:club_id"), updateClub(session))
	mux.HandleFunc(pat.Delete("/clubs/:club_id"), deleteClub(session))

	http.ListenAndServe("localhost:1234", mux)
}

// getConnectionStringFromFile Gets connection string from file
func getConnectionStringFromFile() string {
	absPath, _ := filepath.Abs("./app/connectionString.txt")

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

func getAllClubs(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("klubbkoll").C("clubs")
		var clubs []Club

		err := c.Find(bson.M{}).All(&clubs)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all books: ", err)
			return
		}

		respBody, err := json.MarshalIndent(clubs, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(http.StatusOK)
		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func addClub(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		session := s.Copy()
		defer session.Close()

		var club Club
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&club)
		if err != nil {
			fmt.Println("incorrect body", err)
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB("klubbkoll").C("clubs")
		err = c.Insert(club)
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "Club with this Club Id already exists", http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert book: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+club.ClubID)
		w.WriteHeader(http.StatusCreated)
	}
}

func clubByClubID(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		session := s.Copy()
		defer session.Close()

		clubID := pat.Param(r, "club_id")

		c := session.DB("klubbkoll").C("clubs")

		var club Club
		err := c.Find(bson.M{"club_id": clubID}).One(&club)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed to find club: ", err)
			return
		}

		if club.ClubID == "" {
			ErrorWithJSON(w, "Club not found", http.StatusNotFound)
			return
		}

		respBody, err := json.MarshalIndent(club, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)

	}
}

func updateClub(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		session := s.Copy()
		defer session.Close()

		clubID := pat.Param(r, "club_id")

		var club Club
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&club)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB("klubbkoll").C("clubs")

		err = c.Update(bson.M{"club_id": clubID}, &club)
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed update club: ", err)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Club not found", http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteClub(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		session := s.Copy()
		defer session.Close()

		clubID := pat.Param(r, "club_id")
		c := session.DB("klubbkoll").C("clubs")

		err := c.Remove(bson.M{"club_id": clubID})
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed remove club: ", err)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Club not found", http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)

	}
}

//Club is the main component of this project
type Club struct {
	ClubID       string `json:"club_id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	OpeningHours string `json:"openinghours"`
	Price        string `json:"price"`
}
