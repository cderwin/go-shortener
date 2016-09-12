package main

import (
	"encoding/json"
	"github.com/gocraft/web"
	"log"
	"net/http"
)

func (s *Server) healthcheck(w http.ResponseWriter, r *web.Request) {
	status := struct{ Status string }{"ok"}
	jsonBlob, _ := json.Marshal(status)
	w.Write(jsonBlob)
}

type UrlData struct{ Url string }

func (s *Server) addUrl(w http.ResponseWriter, r *web.Request) {
	var data UrlData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Could not parse body as json", http.StatusBadRequest)
		return
	}

	shortUrl, err := s.Redis.SaveURL(data.Url)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not save url", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(UrlData{Url: shortUrl})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not encode url as json", http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func (s *Server) fetchUrl(w http.ResponseWriter, r *web.Request) {
	shortUrl := r.PathParams["path"]
	longUrl, err := s.Redis.GetURL(shortUrl)
	if err == NilValue {
		http.Error(w, "Shortlink does not exist", 404)
		return
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Url could not be retrieved", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r.Request, longUrl, http.StatusMovedPermanently)
}

func (s *Server) urlStats(w http.ResponseWriter, r *web.Request) {
	shortUrl := r.PathParams["path"]
	stats, err := s.Redis.GetHits(shortUrl)
	if err == NilValue {
		http.Error(w, "Stats do not exist", 404)
		return
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not fetch stats", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(stats)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not encode stats as json", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
