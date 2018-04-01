package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	models "github.com/ElSiipo/Klubbkoll/club"

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

func GetAll(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		session := s.Copy()
		defer session.Close()

		c := session.DB("klubbkoll").C("clubs")
		var clubs []models.Club

		err := c.Find(bson.M{}).All(&clubs)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all clubs: ", err)
			return
		}

		respBody, err := json.MarshalIndent(clubs, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
		fmt.Println("Get all clubs - Status:", http.StatusOK, ":", time.Since(t1))
	}
}

func Store(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		session := s.Copy()
		defer session.Close()

		var club models.Club
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
			log.Println("Failed insert club: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+club.ClubID)
		w.WriteHeader(http.StatusCreated)

		fmt.Println("Add club - Status:", http.StatusOK, ":", time.Since(t1))
	}
}

func GetByID(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		session := s.Copy()
		defer session.Close()

		clubID := pat.Param(r, "club_id")

		c := session.DB("klubbkoll").C("clubs")

		var club models.Club
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
		fmt.Println("Get club by id - Status:", http.StatusOK, ":", time.Since(t1))
	}
}

func Update(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		session := s.Copy()
		defer session.Close()

		clubID := pat.Param(r, "club_id")

		var club models.Club
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
		fmt.Println("Update club - Status:", http.StatusOK, ":", time.Since(t1))
	}
}

// Delete is right now responsible for handling request AND doing DB call..
func Delete(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

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
		fmt.Println("Delete club - Status:", http.StatusOK, ":", time.Since(t1))
	}
}
