package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	if r.Method == "GET" {
		if strings.HasSuffix(urlPath, "/") {
			files := getListOfFiles(urlPath)
			var listFilesHtml string
			listFilesHtml = "<div style='white-space: nowrap;'>"
			listFilesHtml += "<a class='col' href='../'>../</a> <br>"
			for _, file := range files {
				listFilesHtml += fmt.Sprintf("<a class='col' style='width: 50%%' href='%s'> %s </a> <span class='col'>%s</span> <span class='col'>%s</span><br>", file.filename, file.filename, file.formatSize(), file.lastMod)
			}
			listFilesHtml += "</div>"

			fmt.Fprintf(w, `<html>
			<head>
			  <title>Upload</title>
			  <style>
			  .col {
				  display: inline-block;
				  width: 25%%;
				  box-sizing: border-box;
				  margin: 0;
				  padding: 2px;
				  font-size: larger;
				}
				input[type="file"] {
				  font-size: 18px;
				  width: 45vw;
				}
			  </style>
			</head>
			<body>
				%s
				<br><br>
			  <form method="post" enctype="multipart/form-data">
			    <input type="file" name="files" multiple>
				 <br>
				 <button style='margin-top: 1vw; width: 17vw;' type="submit">Upload</button>
			  </form>
			</body>
			</html>`, listFilesHtml)
		} else {
			file, err := http.Dir(".").Open(urlPath)
			if err != nil {
				http.Error(w, fmt.Sprintf("File error: %v", err), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			urlParts := strings.Split(urlPath[1:], "/")
			urlPath := urlParts[len(urlParts)-1]
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", urlPath))
			w.Header().Set("Content-Type", "application/octet-stream")

			// Write the file to the response
			http.ServeContent(w, r, urlPath, time.Now(), file)
		}
	}

	if r.Method == "POST" {
		//r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // limit request body size to 10 MB
		err := r.ParseMultipartForm(10 << 22) // parse form data
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

			decodedFilename, err := url.QueryUnescape(file.Filename)
			if err != nil {
				fmt.Println("Error decoding filename:", err)
				return
			}
			fmt.Println(decodedFilename)
			file_path := path.Join(urlPath[1:], decodedFilename)
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
	http.ListenAndServe("0.0.0.0:8000", nil)
}
