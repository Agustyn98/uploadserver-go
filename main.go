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
		//fmt.Println("POST request")
		//formID := r.FormValue("form_id")
		//fmt.Println("choosing form:")
		if false {
			r.ParseForm()
			value := r.FormValue("dirName")
			dirPath := path.Join(urlPath, value)
			if len(value) > 0 {
				err := os.Mkdir(dirPath[1:], 0777)
				if err != nil {
					fmt.Println("Error creating directory: ", err)
					return
				}
			}

		} else {

			sz := r.Header.Get("Content-Length")
			fmt.Println(sz)
			mr, err := r.MultipartReader()
			if err != nil {
				fmt.Printf("\nmultipart reader erro %s", err)
				return
			}
			for {
				part, err := mr.NextPart()
				fmt.Println("reading next part")
				if err == io.EOF {
					break
				}
				var read int64
				//var p float32
				filePath := path.Join(urlPath, part.FileName())
				fmt.Printf("\nopening part %s", filePath[1:])
				dst, err := os.OpenFile(filePath[1:], os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					fmt.Printf("\nError reading %s", err)
					continue
				}
				fmt.Println("for loop")
				for {
					buffer := make([]byte, 1<<20)
					cBytes, err := part.Read(buffer)
					read = read + int64(cBytes)
					//p = float32(read) / float32(part.Header.Get("Content-Length")) * 100
					//fmt.Printf("progress: %v \n", p)
					fmt.Printf("\nProgress: %d  /  %s ", read, part.Header.Get("Content-Length"))
					fmt.Println(part.Header.Values("Content-Length"))
					dst.Write(buffer[0:cBytes])
					if err == io.EOF {
						fmt.Printf("\nError when reading from buffer: %s\n", err)
						break
					}
				}
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
	fmt.Printf("Listening on %s:%s", getIp(), port)
	err := http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		fmt.Printf("\nError: %s", err)
	}
}
