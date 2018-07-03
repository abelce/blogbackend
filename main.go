package main

import (
	"net/http"
	"fmt"
	"io"
	"os"
	"github.com/satori/go.uuid"
	"file/utils"
	"file/port"
	"strconv"
	"encoding/json"
)

type handler struct {
	next http.Handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "x-requested-with")
	w.Header().Add("Content-Type", "application/x-msdownload")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
	} else {
		// w.Header.
		// w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		// w.Header().Set("Content-Type", "application/json;charset=utf-8")
		h.next.ServeHTTP(w, r)

	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "x-requested-with")
	w.Header().Add("Content-Type", "application/x-msdownload")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")

		if err != nil {
			utils.HandleHTTPError(w, err);
		}
		defer file.Close()
		path, _ := uuid.NewV4()
		f, err := os.OpenFile("/data/upload_files/" + path.String(), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			utils.HandleServerError(w, err);
		}
		defer f.Close()
		io.Copy(f, file)

		id, _ := uuid.NewV4();
		_, err = port.InsertFile(id.String(), path.String(), handler.Filename, "--", utils.GetTimestamp())
		if err != nil {
			utils.HandleServerError(w, err);
			fmt.Fprintln(w, "upload failed123!")
			return;
		}
		fmt.Fprintln(w, "upload success!")
	}
}

func GetFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "x-requested-with")
	w.Header().Add("Content-Type", "application/x-msdownload")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
	} else {
		r.ParseForm();
		var pageSize int = 10;
		var pageNum int = 1;
		if r.FormValue("pageSize") != "" {
			pageSize, _ = strconv.Atoi(r.FormValue("pageSize"))
		}
		if r.FormValue("pageNum") != "" {
			pageNum, _ = strconv.Atoi(r.FormValue("pageNum"))
		}
		files, err := port.GetFiles(pageSize * (pageNum - 1), pageNum * pageSize)
		if err != nil {
			utils.HandleServerError(w, err);
		}
		tmp, err := json.Marshal(files);

		wc := 0
		for wc < len(tmp) {
			n, err := w.Write(tmp)
			if err != nil {
				utils.HandleServerError(w, err);	
			}
			wc += n;
		}
    }
}

func main() {
	//var routeHandler http.Handler = &handler{
	//		next: application.NewRouter(),
	//}

	fs := http.FileServer(http.Dir("/data/upload_files"))
	// http.Handle("/", routeHandler)
	http.Handle("/", http.StripPrefix("/", fs))
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/files", GetFiles)
	http.ListenAndServe(":9010", nil)
}
