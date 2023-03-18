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
	ip := r.RemoteAddr
	method := r.Method
	request := r.URL.Path
	fmt.Printf("Request from %s: %s %s\n", ip, method, request)
	if r.Method == "GET" {
		if strings.HasSuffix(urlPath, "/") {
			files := getListOfFiles(urlPath)
			var listFilesHtml string
			listFilesHtml = "<div style='white-space: nowrap;'>"
			listFilesHtml += "<a class='col' href='../'>../</a> <br>"
			for _, file := range files {
				listFilesHtml += fmt.Sprintf("<a class='col' style='width: 50%%;' href='%s'> %s </a> <span class='col'>%s</span> <span class='col'>%s</span><br>", file.filename, file.filename, file.formatSize(), file.lastMod)
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
				  overflow: hidden;
				  word-wrap: break-word;
				}
				input[type="file"] {
				  font-size: 18px;
				  width: 45vw;
				}
				input[type="text"] {
				  font-size: 18px;
				  width: 15vw;
				}
			  </style>
			</head>
			<body>
				%s
				<br><br>
			<div style="display: flex; gap: 30px;">
			  <form method="post">
			    <input type="text" name="dirName" multiple>
				<input type="hidden" name="form_id" value="folder">
				 <br>
				 <button style='margin-top: 1vw; width: 17vw;' type="submit">Create folder</button>
			  </form>
				<br><br>
			  <form method="post" enctype="multipart/form-data">
			    <input type="file" name="files" multiple>
				<input type="hidden" name="form_id" value="files">
				 <br>
				 <button style='margin-top: 1vw; width: 17vw;' type="submit">Upload</button>
			  </form>
			</div>
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
		formID := r.FormValue("form_id")
		if formID == "folder" {
			r.ParseForm()
			value := r.FormValue("dirName")
			if len(value) > 0 {
				err := os.Mkdir(value, 0755)
				if err != nil {
					fmt.Println("Error creating directory: ", err)
					return
				}
			}

		} else {

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
				fmt.Printf("File %s uploaded successfully with %d bytes\n", file.Filename, nBytes)
			}
		}
		http.Redirect(w, r, urlPath, http.StatusSeeOther)
	}
}

func main() {
	http.HandleFunc("/", uploadHandler)

	var port string
	if len(os.Args) < 2 {
		port = "8000"
	} else {
		port = os.Args[1]
	}
	fmt.Println("Listening on :" + port)
	err := http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
