package dfi

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ListenHTTP(addr string) {
	http.HandleFunc("/ping", HandleMethod("GET", HandlePing))           // address
	http.HandleFunc("/announce", HandleMethod("GET", HandleAnnounce))   // address
	http.HandleFunc("/search", HandleMethod("POST", HandleSearch))      // query, address?, page?, recursive?
	http.HandleFunc("/recent", HandleMethod("GET", HandleRecent))       // address?, page?
	http.HandleFunc("/popular", HandleMethod("GET", HandlePopular))     // address?, page?
	http.HandleFunc("/mirror", HandleMethod("GET", HandleMirror))       // address, progress?
	http.HandleFunc("/index", HandleMethod("GET", HandleIndex))         // address?, since?
	http.HandleFunc("/addpost", HandleMethod("POST", HandleAddPost))    // post, index?
	http.HandleFunc("/resolve", HandleMethod("GET", HandleResolve))     // address
	http.HandleFunc("/bootstrap", HandleMethod("GET", HandleBootstrap)) // address
	http.HandleFunc("/peers", HandleMethod("GET", HandlePeers))
	http.HandleFunc("/peers/add", HandleMethod("GET", HandlePeersAdd)) // address
	http.HandleFunc("/set", HandleMethod("POST", HandleSet))           // key, value
	http.HandleFunc("/get", HandleMethod("GET", HandleGet))            // key
	http.HandleFunc("/explore", HandleMethod("GET", HandleExplore))
	http.HandleFunc("/map", HandleMethod("GET", HandleMap))

	log.Fatal(http.ListenAndServe(addr, nil))
}

func HandleMethod(method string, handle func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Must use " + method))
		} else {
			handle(w, r)
		}
	}
}

func HandlePing(w http.ResponseWriter, r *http.Request) {

}

func HandleAnnounce(w http.ResponseWriter, r *http.Request) {

}

func HandleRecent(w http.ResponseWriter, r *http.Request) {

}

func HandlePopular(w http.ResponseWriter, r *http.Request) {

}

func HandleMirror(w http.ResponseWriter, r *http.Request) {

}

func HandleIndex(w http.ResponseWriter, r *http.Request) {

}

func HandleAddPost(w http.ResponseWriter, r *http.Request) {

}

func HandleResolve(w http.ResponseWriter, r *http.Request) {

}

func HandleBootstrap(w http.ResponseWriter, r *http.Request) {

}

func HandlePeers(w http.ResponseWriter, r *http.Request) {

}

func HandlePeersAdd(w http.ResponseWriter, r *http.Request) {

}

func HandleSet(w http.ResponseWriter, r *http.Request) {

}

func HandleGet(w http.ResponseWriter, r *http.Request) {

}

func HandleExplore(w http.ResponseWriter, r *http.Request) {

}

func HandleMap(w http.ResponseWriter, r *http.Request) {

}
