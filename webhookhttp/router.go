package webhookhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/BadgeForce/doug"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

//Route . . .
//struct used to define a route,
//Name: description of the route
//Method: http method
//Pattern: actual route endpoint
//HandlerFunc: function to handle the route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes slice of routes that are registered with the mux router. All new routes can be defined here.
type Routes []Route

//NewRouter ... register each route and returns a new mux router instance
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

//ArtifactRelease . . . http handler that will handle github release events, particularly your truffle projects
func ArtifactRelease(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var headers map[string]string
	b, err := json.Marshal(r.Header)
	if err != nil {
		writeError(w, fmt.Sprintf("Failed processing hook! ('%+v')", err), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(b, &headers); err != nil {
		writeError(w, fmt.Sprintf("Failed processing hook! ('%+v')", err), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, fmt.Sprintf("Failed processing hook! ('%+v')", err), http.StatusInternalServerError)
		return
	}

	hc, err := doug.ParseHook([]byte(doug.Configs.Github.Secret), headers, string(body))
	if err != nil {
		writeError(w, fmt.Sprintf("Failed processing hook! ('%+v')", err), http.StatusBadRequest)
		return
	}

	evt := github.ReleaseEvent{}
	if err := json.Unmarshal(hc.Payload, &evt); err != nil {
		writeError(w, fmt.Sprintf("Failed processing hook! ('%+v')", err), http.StatusBadRequest)
		return
	}

	errors := doug.UploadArtifacts(evt)
	if errors != nil {
		writeError(w, fmt.Sprintf("Failed processing hook! ('%+v')", errors), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "{}")
	return
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	io.WriteString(w, message)
}

var routes = Routes{
	Route{
		"Artifact release hook handler", "POST", "/artifact", ArtifactRelease,
	},
}
