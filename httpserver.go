package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (lp *LocalPeer) ListenHTTP(addr string) {
	http.HandleFunc("/ping", HandleMethod("GET", lp.HandlePing))
	http.HandleFunc("/announce", HandleMethod("GET", lp.HandleAnnounce))
	http.HandleFunc("/search", HandleMethod("POST", lp.HandleSearch))
	http.HandleFunc("/recent", HandleMethod("GET", lp.HandleRecent))
	http.HandleFunc("/popular", HandleMethod("GET", lp.HandlePopular))
	http.HandleFunc("/mirror", HandleMethod("GET", lp.HandleMirror))
	http.HandleFunc("/mirror/progress", HandleMethod("GET", lp.HandleMirrorProgress))
	http.HandleFunc("/index", HandleMethod("GET", lp.HandleIndex))
	http.HandleFunc("/add/post", HandleMethod("POST", lp.HandleAddPost))
	http.HandleFunc("/resolve", HandleMethod("GET", lp.HandleResolve))
	http.HandleFunc("/bootstrap", HandleMethod("GET", lp.HandleBootstrap))
	http.HandleFunc("/peers", HandleMethod("GET", lp.HandlePeers))
	http.HandleFunc("/add/peer", HandleMethod("GET", lp.HandleAddPeer))
	http.HandleFunc("/set", HandleMethod("POST", lp.HandleSet)) // key, value
	http.HandleFunc("/get", HandleMethod("GET", lp.HandleGet))  // key
	http.HandleFunc("/explore", HandleMethod("GET", lp.HandleExplore))
	http.HandleFunc("/map", HandleMethod("GET", lp.HandleMap))

	log.Fatal(http.ListenAndServe(addr, nil))
}

type HTTPResponse struct {
	Status int
	Data   interface{}
	Error  error
}

func WriteHTTPResponse(w http.ResponseWriter, hr HTTPResponse) {
	if hr.Status != 0 || hr.Status != http.StatusOK {
		w.WriteHeader(hr.Status)
		io.WriteString(w, hr.Error.Error())
		return
	}

	b, err := json.Marshal(hr.Data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write(b)
}

func HandleMethod(method string, handle func(w http.ResponseWriter, r *http.Request) HTTPResponse) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Must use " + method))
		} else {
			WriteHTTPResponse(w, handle(w, r))
		}
	}
}

var (
	ErrNoAddress  = errors.New("A Public Address must be specified")
	ErrBadAddress = errors.New("Could not decode Public Address")

	ErrNoQuery = errors.New("A query string must be specified")

	ErrNoKey = errors.New("A key string must be specified")
)

func (lp *LocalPeer) HandlePing(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address := r.FormValue("address")

	if address == "" {
		return HTTPResponse{http.StatusBadRequest, nil, ErrNoAddress}
	}

	addr, err := DecodeAddress()
	if err != nil {
		return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
	}

	return lp.Ping(addr)
}

func (lp *LocalPeer) HandleAnnounce(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address := r.FormValue("address")

	if address == "" {
		return HTTPResponse{http.StatusBadRequest, nil, ErrNoAddress}
	}

	addr, err := DecodeAddress()
	if err != nil {
		return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
	}

	return lp.Announce(addr)
}

func (lp *LocalPeer) HandleSearch(w http.ResponseWriter, r *http.Request) HTTPResponse {
	query = r.FormValue("query")

	if query == "" {
		return HTTPResponse{http.StatusBadRequest, nil, ErrNoQuery}
	}

	address = r.FormValue("address")
	var addr *Address

	if address != "" {
		var err error
		addr, err = DecodeAddress(address)

		if err != nil {
			return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
		}
	}

	page, _ := strconv.Atoi(r.FormValue("page"))

	recursive, _ := strconv.ParseBool(r.FormValue("recursive"))

	return lp.Search(query, addr, page, recursive)
}

func (lp *LocalPeer) HandleRecent(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address = r.FormValue("address")
	var addr *Address

	if address != "" {
		var err error
		addr, err = DecodeAddress(address)

		if err != nil {
			return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
		}
	}

	page, _ := strconv.Atoi(r.FormValue("page"))

	return lp.Recent(addr, page)
}

func (lp *LocalPeer) HandlePopular(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address = r.FormValue("address")
	var addr *Address

	if address != "" {
		var err error
		addr, err = DecodeAddress(address)

		if err != nil {
			return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
		}
	}

	page, _ := strconv.Atoi(r.FormValue("page"))

	return lp.Popular(addr, page)
}

func (lp *LocalPeer) HandleMirror(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address := r.FormValue("address")

	if address == "" {
		return HTTPResponse{http.StatusBadRequest, nil, ErrNoAddress}
	}

	addr, err := DecodeAddress()
	if err != nil {
		return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
	}

	return lp.Mirror(addr)
}

func (lp *LocalPeer) HandleMirrorProgress(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address := r.FormValue("address")

	if address == "" {
		return HTTPResponse{http.StatusBadRequest, nil, ErrNoAddress}
	}

	addr, err := DecodeAddress()
	if err != nil {
		return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
	}

	return lp.MirrorProgress(addr)
}

func (lp *LocalPeer) HandleIndex(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address := r.FormValue("address")
	var addr *Address

	if address != "" {
		var err error
		addr, err = DecodeAddress(address)

		if err != nil {
			return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
		}
	}

	since, _ := strconv.Atoi(r.FormValue("since"))

	return lp.Index(addr, since)
}

func (lp *LocalPeer) HandleAddPost(w http.ResponseWriter, r *http.Request) HTTPResponse {
	index := strconv.ParseBool(r.FormValue("index"))

	var post Post
	err := json.Unmarshal([]byte(r.FormValue("post")), &post)

	if err == nil {
		return lp.AddPost(&post, index)
	}

	if r.FormValue("post") != "" {
		return HTTPResponse{http.StatusBadRequest, nil, err}
	}

	post.Id, _ = strconv.Atoi(r.FormValue("id"))
	post.URI = r.FormValue("uri")

	post.Title = r.FormValue("title")
	post.Size = strconv.Atoi(r.FormValue("size"))
	post.Files = strings.Split(r.FormValue("files"), ",")

	post.Time = strconv.Itoa(r.FormValue("time"))

	post.Tags = strings.Split(r.FormValue("tags"), ",")

	json.Unmarshal([]byte(r.FormValue("meta")), &post.Meta)

	return lp.AddPost(&post, index)
}

func (lp *LocalPeer) HandleResolve(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address := r.FormValue("address")

	if address == "" {
		return HTTPResponse{http.StatusBadRequest, nil, ErrNoAddress}
	}

	addr, err := DecodeAddress()
	if err != nil {
		return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
	}

	return lp.Resolve(addr)
}

func (lp *LocalPeer) HandleBootstrap(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address := r.FormValue("address")

	if address == "" {
		return HTTPResponse{http.StatusBadRequest, nil, ErrNoAddress}
	}

	addr, err := DecodeAddress()
	if err != nil {
		return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
	}

	return lp.Bootstrap(addr)
}

func (lp *LocalPeer) HandlePeers(w http.ResponseWriter, r *http.Request) HTTPResponse {
	return lp.Peers()
}

func (lp *LocalPeer) HandleAddPeer(w http.ResponseWriter, r *http.Request) HTTPResponse {
	address := r.FormValue("address")

	if address == "" {
		return HTTPResponse{http.StatusBadRequest, nil, ErrNoAddress}
	}

	addr, err := DecodeAddress()
	if err != nil {
		return HTTPResponse{http.StatusBadRequest, nil, ErrBadAddress}
	}

	return lp.AddPeer(addr)
}

func (lp *LocalPeer) HandleSet(w http.ResponseWriter, r *http.Request) HTTPResponse {
	var identity Identity

	data := r.FormValue("data")

	err := json.Unmarshal([]byte(data), &identity)

	if err == nil {
		return lp.Set(identity)
	}

	if data != "" {
		return HTTPResponse{http.StatusBadRequest, nil, err}
	}

	identity.Name = r.FormValue("name")
	identity.Desc = r.FormValue("desc")
	identity.Public = r.FormValue("public")

	return lp.Set(identity)
}

func (lp *LocalPeer) HandleGet(w http.ResponseWriter, r *http.Request) HTTPResponse {
	key := r.FormValue("key")
	if key != "" {
		return lp.Get(key)
	}

	keys := strings.Split(r.FormValue("keys"), ",")
	return lp.Get(keys...)
}

func (lp *LocalPeer) HandleExplore(w http.ResponseWriter, r *http.Request) HTTPResponse {
	return lp.Explore()
}

func (lp *LocalPeer) HandleMap(w http.ResponseWriter, r *http.Request) HTTPResponse {
	return lp.Map()
}
