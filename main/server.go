package main

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func getExtension(f *multipart.FileHeader) string {
	return strings.Split(f.Filename, ".")[1]
}

func handleFileUpload(index int, files []*multipart.FileHeader) {
	fmt.Println("server.uploadFile")
	file, err := files[index].Open()

	if err != nil {
		fmt.Println("Error retrieving the file")
		fmt.Println(err)
		return
	}
	defer file.Close()
	extension := getExtension(files[index])

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	filepath := files[index].Filename
	f, err := os.Create(filepath)
	f.Write(fileBytes)
	f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Fprintf(w, "Successfully Uploaded File\n")
	stat, err := os.Stat(files[index].Filename)
	if err != nil {
		fmt.Println("Error reading file Stat")
		fmt.Println(err)
	}
	dirname := strings.Split(stat.ModTime().String(), " ")[0]
	fmt.Println(dirname)
	os.Mkdir(dirname, os.ModePerm)
	destination := dirname + "/" + extension
	os.Mkdir(destination, os.ModePerm)
	err = os.Rename(filepath, destination+"/"+filepath)
	fmt.Println(err)
}

func Upload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm

	files := m.File["myFile"]

	for i := range files {
		handleFileUpload(i, files)
	}
	_, err = fmt.Fprintf(w, "Successfully Uploaded all files")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})
	r.HandleFunc("/upload", Upload).Methods("POST")
	r.HandleFunc("/uploadSingle", Upload).Methods("POST")

	http.ListenAndServe(":8080", r)
}
