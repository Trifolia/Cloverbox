package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/upload", uploadHandler)

	err := http.ListenAndServe(":" + strconv.Itoa(*port), nil)
	if err != nil {
		panic(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseMultipartForm(10 << 20) //maximum of 10MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get the file from the request
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	//determine what folder to put the file in
	var folder string
	ext := filepath.Ext(handler.Filename)
	switch ext {
	case ".png":
		folder = "images/"
	case ".ogg":
		folder = "sound/"
	default: 
		return //file type not supported
	}

	//check if the file exists
	if _, err := os.Stat("public/" + folder + handler.Filename); err == nil {
		http.Error(w, "File already exists", http.StatusInternalServerError)
		return
	}

	//write the file to the filesystem
	fileName := "public/" + folder + handler.Filename
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	w.Write([]byte("ok"))
}
