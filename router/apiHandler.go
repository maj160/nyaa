package router

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util"
	"github.com/gorilla/mux"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}

	nbTorrents := 0
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}

	torrents, nbTorrents, err := torrentService.GetAllTorrents(maxPerPage, maxPerPage*(pagenum-1))
	if err != nil {
		util.SendError(w, err, 400)
		return
	}

	b := model.ApiResultJson{
		Torrents: model.TorrentsToJSON(torrents),
	}
	b.QueryRecordCount = maxPerPage
	b.TotalRecordCount = nbTorrents
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ApiViewHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	b := model.ApiResultJson{Torrents: []model.TorrentsJson{}}

	torrent, err := torrentService.GetTorrentById(id)
	res := torrent.ToJson()
	b.Torrents = append(b.Torrents, res)

	b.QueryRecordCount = 1
	b.TotalRecordCount = 1
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ApiUploadHandler(w http.ResponseWriter, r *http.Request) {
	if config.UploadsDisabled == 1 {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	//verify token
	//token := r.Header.Get("Authorization")

	decoder := json.NewDecoder(r.Body)
	b := model.TorrentsJson{}
	err := decoder.Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	category, sub_category, err := ValidateJson(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //406?
		return
	}

	torrent := model.Torrents{
		Name:         b.Name,
		Category:     category,
		Sub_Category: sub_category,
		Status:       1,
		Hash:         b.Hash,
		Date:         time.Now(),
		Filesize:     0,
		Description:  string(b.Description)}
	db.ORM.Create(&torrent)
	fmt.Printf("%+v\n", torrent)
}
