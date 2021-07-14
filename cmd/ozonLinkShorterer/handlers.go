package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"ozonLinkShorterer/cmd/ozonLinkShorterer/models"
	"strconv"
)

var codeDict = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

var base = len(codeDict)

//Парсинг JSON-а
func decode(decodeStruct interface{}, requestBody io.ReadCloser) error {
	decoder := json.NewDecoder(requestBody)
	decoder.DisallowUnknownFields()
	return decoder.Decode(decodeStruct)
}

//Отправка ошибки
func respondError(w http.ResponseWriter, err error, status int) {
	log.Println(err.Error(), status)
	http.Error(w, strconv.Itoa(status), status)
}

//Отправка ответа
func sendResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		respondError(w, errors.New("Encoding error!"), http.StatusInternalServerError)
		return
	}
}

//Шифрование id БД в короткую ссылку
func encodeId(id int64) string {
	var buffer bytes.Buffer
	base64 := int64(base)
	for id > 0 {
		buffer.WriteByte(codeDict[id%base64])
		id /= base64
	}
	bufferBytes := buffer.Bytes()
	for i, j := 0, len(bufferBytes)-1; i < j; i, j = i+1, j-1 {
		bufferBytes[i], bufferBytes[j] = bufferBytes[j], bufferBytes[i]
	}
	return string(bufferBytes)
}

//Дешифровка короткой ссылки в id в базе данных
func decodeShortUrl(url string) int64 {
	var id int64 = 0
	urlLen := len(url) - 1
	for i, c := range url {
		index := int64(bytes.IndexByte(codeDict, byte(c)))
		id += int64(math.Pow(float64(base), float64(urlLen-i))) * index
	}
	return id
}

//Получение короткой ссылки
func shortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondError(w, errors.New("Only POST method supported"), http.StatusBadRequest)
		return
	}
	var urlReq models.JsonUrl
	if err := decode(&urlReq, r.Body); err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}
	var checkId int64
	var response models.JsonUrl
	checkErr := dbPointer.QueryRow("SELECT id FROM urls WHERE originalUrl = ? ", urlReq.Url).Scan(&checkId)
	if checkErr != nil {
		if checkErr == sql.ErrNoRows {
			r, insertErr := dbPointer.Exec("INSERT INTO urls (originalUrl) VALUES (?)", urlReq.Url)
			if insertErr != nil {
				respondError(w, insertErr, http.StatusBadRequest)
				return
			}
			lastId, _ := r.LastInsertId()
			response.Url = encodeId(lastId)
			sendResponse(w, response)
			return
		}
		respondError(w, checkErr, http.StatusInternalServerError)
		return
	}
	response.Url = encodeId(checkId)
	sendResponse(w, response)
}

//Получение оригинальной ссылки
func retrieveUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondError(w, errors.New("Only POST method supported"), http.StatusBadRequest)
		return
	}
	var urlReq models.JsonUrl
	if err := decode(&urlReq, r.Body); err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}
	urlId := decodeShortUrl(urlReq.Url)
	var retUrl string
	checkErr := dbPointer.QueryRow("SELECT originalUrl FROM urls WHERE id = ? ", urlId).Scan(&retUrl)
	if checkErr != nil {
		respondError(w, checkErr, http.StatusBadRequest)
		return
	}
	response := models.JsonUrl{retUrl}
	sendResponse(w, response)
}
