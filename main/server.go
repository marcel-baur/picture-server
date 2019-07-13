package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func UploadSeveral(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm

	files := m.File["myFile"]
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// stat, err := os.Stat(files[i].Filename)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// subfolder := strings.Split( stat.ModTime().String(), " ")[0]
		err = os.Mkdir("img", os.ModePerm)
		if err != nil {
			fmt.Println("here1")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = os.Rename(files[i].Filename, "img/"+files[i].Filename)
		if err != nil {
			fmt.Println("here")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// dst, err := os.Create("/img/upload/"+files[i].Filename)
		// defer dst.Close()
		// if err != nil {
		// 	fmt.Println("here")
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// if _, err := io.Copy(dst, file); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

	}
	fmt.Fprintf(w, "Successfully Uploaded %v File(s)\n", len(files))
}

func Upload(w http.ResponseWriter, r *http.Request) {

	fmt.Println("server.uploadFile")
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error retrieving the file")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Println("Uploaded File: %v\n", handler.Filename)
	fmt.Println("File Size: %v\n", handler.Size)
	fmt.Println("MIME Header: %v\n", handler.Header)

	extension := strings.Split(handler.Filename, ".")[1]
	fmt.Println(extension)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	filepath := handler.Filename
	f, err := os.Create(filepath)
	f.Write(fileBytes)
	f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
	stat, err := os.Stat(handler.Filename)
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})
	r.HandleFunc("/upload", UploadSeveral).Methods("POST")
	r.HandleFunc("/uploadSingle", Upload).Methods("POST")

	http.ListenAndServe(":8080", r)
}
