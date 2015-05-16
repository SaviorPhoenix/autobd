package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type File struct {
	Name     string           `json:"name"`
	Size     int64            `json:"size"`
	ModTime  time.Time        `json:"lastModified"`
	Mode     os.FileMode      `json:"fileMode"`
	IsDir    bool             `json:"isDir"`
	Manifest map[string]*File `json:"manifest,omitempty"`
}

var (
	apiVersion string = "v0"
	version    string = "0.1"
	commit     string
)

func NewFile(name string, size int64, modtime time.Time, mode os.FileMode, isDir bool) *File {
	return &File{name, size, modtime, mode, isDir, nil}
}

func GetManifest(dirPath string) (map[string]*File, error) {
	list, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	manifest := make(map[string]*File)
	for _, child := range list {
		childPath := path.Join(dirPath, child.Name())
		manifest[childPath] = NewFile(childPath, child.Size(), child.ModTime(), child.Mode(), child.IsDir())
		if child.IsDir() == true {
			childContent, err := GetManifest(childPath)
			if err != nil {
				return nil, err
			}
			manifest[childPath].Manifest = childContent
		}
	}
	return manifest, nil
}

func LogHttp(r *http.Request) {
	log.Printf("%s %s %s %s", r.Method, r.URL, r.RemoteAddr, r.UserAgent())
}

func LogHttpErr(w http.ResponseWriter, r *http.Request, err error, status int) {
	log.Printf("Returned error \"%s\" (HTTP %s) to %s", err.Error(), http.StatusText(status), r.RemoteAddr)
	serialErr, _ := json.Marshal(err.Error())
	http.Error(w, string(serialErr), status)
}

func GetQueryValue(name string, w http.ResponseWriter, r *http.Request) string {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		LogHttpErr(w, r, err, http.StatusInternalServerError)
		return ""
	}
	value := query.Get(name)
	if len(value) == 0 || value == "" {
		LogHttpErr(w, r, fmt.Errorf("Must specify %s", name), http.StatusBadRequest)
		return ""
	}
	return value
}

func ServeManifest(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
	dir := GetQueryValue("dir", w, r)
	if dir == "" {
		return
	}
	manifest, err := GetManifest(dir)
	if err != nil {
		LogHttpErr(w, r, err, http.StatusInternalServerError)
		return
	}
	serial, _ := json.MarshalIndent(&manifest, "  ", "  ")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(serial))
}

func versionInfo() {
	if commit == "" {
		commit = "unknown"
	}
	fmt.Printf("Autobd version %s (API %s) (git commit %s)\n", version, apiVersion, commit)
}

func main() {
	versionInfo()
	http.HandleFunc("/"+apiVersion+"/manifest", ServeManifest)
	log.Panic(http.ListenAndServe(":8080", nil))
}
