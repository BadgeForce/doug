package webhookhttp

import (
	"encoding/json"
	"io"
	"log"
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

	hc, err := doug.ParseHook([]byte(doug.Configs.Github.Secret), r)

	w.Header().Set("Content-type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed processing hook! ('%s')", err)
		io.WriteString(w, "{}")
		return
	}

	evt := github.ReleaseEvent{}
	if err := json.Unmarshal(hc.Payload, &evt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed processing hook! ('%s')", err)
		io.WriteString(w, "{}")
		return
	}

	errors := doug.UploadArtifacts(evt)
	if errors != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed processing hook! ('%s')", errors)
		io.WriteString(w, "{}")
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "{}")
	return
}

var routes = Routes{
	Route{
		"Artifact release hook handler", "POST", "/artifact", ArtifactRelease,
	},
}
