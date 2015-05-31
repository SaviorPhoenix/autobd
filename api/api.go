//Package api implements the endpoints and utility necessary to present
//a consistent API to autobd-nodes
package api

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/tywkeene/autobd/logging"
	"github.com/tywkeene/autobd/manifest"
	"github.com/tywkeene/autobd/packing"
	"github.com/tywkeene/autobd/version"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

//Check and make sure the client wants or can handle gzip, and replace the writer if it
//can, if not, simply use the normal http.ResponseWriter
func GzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

//GetQueryValue() takes a name of a key:value pair to fetchf rom a URL encoded query,
//a http.ResponseWriter 'w', and a http.Request 'r'. In the event that an error is encountered
//the error will be returned to the client via logging facilities that use 'w' and 'r'
func GetQueryValue(name string, w http.ResponseWriter, r *http.Request) string {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		logging.LogHttpErr(w, r, fmt.Errorf("Query parse error"), http.StatusInternalServerError)
		return ""
	}
	value := query.Get(name)
	if len(value) == 0 || value == "" {
		logging.LogHttpErr(w, r, fmt.Errorf("Must specify %s", name), http.StatusBadRequest)
		return ""
	}
	return value
}

//ServeManifest() is the http handler for the "/manifest" API endpoint. It will extract the requested
//directory to be manifested by calling GetQueryValue(), then returns writes it to the client as a
//map[string]*manifest.Manifest encoded in json
func ServeManifest(w http.ResponseWriter, r *http.Request) {
	logging.LogHttp(r)
	dir := GetQueryValue("dir", w, r)
	if dir == "" {
		return
	}
	dirManifest, err := manifest.GetManifest(dir)
	if err != nil {
		logging.LogHttpErr(w, r, fmt.Errorf("Error getting manifest"), http.StatusInternalServerError)
		return
	}
	serial, _ := json.MarshalIndent(&dirManifest, "  ", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "Autobd v"+version.Server())
	io.WriteString(w, string(serial))
}

//ServeServerVer() is the http handler for the "/version" http API endpoint.
//It writes the json encoded struct version.VersionInfo to the client
func ServeServerVer(w http.ResponseWriter, r *http.Request) {
	logging.LogHttp(r)
	serialVer, _ := json.MarshalIndent(&version.VersionInfo{version.Server(), version.API(),
		version.Commit()}, "  ", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "Autobd v"+version.Server())
	io.WriteString(w, string(serialVer))
}

//ServeSync() is the http handler for the "/sync" http API endpoint.
//It takes the requested directory or file name passed as a url parameter "grab" i.e "/sync?grab=file1"
//If the requested file is a directory, it will be tarballed and the "Content-Type" http-header will be
//set to "application/x-tar".
//If the file is a normal file, it will be served with http.ServeContent(), with the Content-Type http-header
//set by http.ServeContent()
func ServeSync(w http.ResponseWriter, r *http.Request) {
	logging.LogHttp(r)
	grab := GetQueryValue("grab", w, r)
	if grab == "" {
		return
	}
	fd, err := os.Open(grab)
	if err != nil {
		logging.LogHttpErr(w, r, fmt.Errorf("Error getting file"), http.StatusInternalServerError)
		return
	}
	defer fd.Close()
	info, err := fd.Stat()
	if err != nil {
		logging.LogHttpErr(w, r, fmt.Errorf("Error getting file"), http.StatusInternalServerError)
		return
	}
	if info.IsDir() == true {
		w.Header().Set("Content-Type", "application/x-tar")
		if err := packing.PackDir(grab, w); err != nil {
			logging.LogHttpErr(w, r, fmt.Errorf("Error packing directory"), http.StatusInternalServerError)
			return
		}
		return
	}
	http.ServeContent(w, r, grab, info.ModTime(), fd)
}

func SetupRoutes() {
	http.HandleFunc("/"+version.API()+"/manifest", GzipHandler(ServeManifest))
	http.HandleFunc("/"+version.API()+"/sync", GzipHandler(ServeSync))
	http.HandleFunc("/version", GzipHandler(ServeServerVer))
}
