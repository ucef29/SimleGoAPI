package handlers

import (
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func UsersRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	if path == "/users" {
		switch r.Method {
		case http.MethodGet:
			usersGetAll(w, r)
			return
		case http.MethodPost:
			usersPostOne(w, r)
			return
		case http.MethodOptions:
			postOptionsResponse(w, []string{http.MethodGet, http.MethodPost, http.MethodOptions}, nil)
			return
		default:
			postError(w, http.StatusMethodNotAllowed)
			return
		}
	}
	path = strings.TrimPrefix(path, "/users/")
	if !bson.IsObjectIdHex(path) {
		postError(w, http.StatusNotFound)
		return
	}

	id := bson.ObjectIdHex(path)

	switch r.Method {
	case http.MethodGet:
		userGetOne(w, r, id)
		return
	case http.MethodPut:
		usersPutOne(w, r, id)
		return
	case http.MethodPatch:
		usersPatchOne(w, r, id)
		return
	case http.MethodDelete:
		usersDeleteOne(w, r, id)
		return
	case http.MethodOptions:
		postOptionsResponse(w, []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}, nil)
		return
	default:
		postError(w, http.StatusMethodNotAllowed)
		return
	}
}
