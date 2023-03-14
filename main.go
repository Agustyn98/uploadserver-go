package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	fmt.Fprintf(w, "<html><head><title>Hello</title></head><body><h1>Hello, World!</h1></body></html>")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	url_path := r.URL.Path
	fmt.Printf("Current path: %s\n", url_path)
	if r.Method == "GET" {
		if strings.HasSuffix(url_path, "/") {
			//http.ServeFile(w, r, "static/upload.html")
			files := getListOfFiles(url_path)
			var listFilesHtml string
			listFilesHtml = "<div>"
			listFilesHtml += "<a class='link' href='../'>../</a>"
			for _, file := range files {
				listFilesHtml += fmt.Sprintf("<a class='link' href='%s'> %s </a>", file.filename, file.filename)
			}
			listFilesHtml += "</div>"

			fmt.Fprintf(w, `<html>
			<head>
			  <title>Upload</title>
			  <style>
			  .link {
				display: inline-block;
				width: 500px;
				}
			  </style>
			</head>
			<body>
				%s
			  <form method="post" enctype="multipart/form-data">
			    <input type="file" name="files" multiple>
			    <input type="submit">
			  </form>
			</body>
			</html>`, listFilesHtml)
		} else {
			file, err := http.Dir(".").Open(url_path)
			if err != nil {
				http.Error(w, fmt.Sprintf("File error: %v", err), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			// Set the response headers to indicate that the response is a file
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url_path[1:]))
			w.Header().Set("Content-Type", "application/octet-stream")
			//w.Header().Set("Content-Length", fmt.Sprintf("%d", ))

			// Write the file to the response
			http.ServeContent(w, r, url_path, time.Now(), file)
		}
	}

	if r.Method == "POST" {
		//r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // limit request body size to 10 MB
		err := r.ParseMultipartForm(10 << 20) // parse form data
		if err != nil {
			fmt.Fprintln(w, "Error parsing form data:", err)
			return
		}

		files := r.MultipartForm.File["files"] // get files from form data
		for _, file := range files {           // iterate over files
			f, err := file.Open() // open file for reading
			if err != nil {
				fmt.Fprintln(w, "Error opening file:", err)
				continue
			}
			defer f.Close()

			file_path := path.Join(url_path[1:], file.Filename)
			dst, err := os.Create(file_path) // create destination file for writing
			if err != nil {
				fmt.Fprintln(w, "Error creating file:", err)
				continue
			}
			defer dst.Close()

			nBytes, err := io.Copy(dst, f) // copy bytes from source to destination
			if err != nil {
				fmt.Fprintln(w, "Error copying file:", err)
				continue
			}
			fmt.Fprintf(w, "File %s uploaded successfully with %d bytes\n", file.Filename, nBytes)
		}
	}
}

func main() {
	http.HandleFunc("/", uploadHandler)
	fmt.Println("Listening on :8000...")
	http.ListenAndServe(":8000", nil)
	//var dir string
	//dir = "static"
	//l := getListOfFiles(dir)
	//for _, file := range l {
	//	fmt.Println(file.filename)
	//}

}
