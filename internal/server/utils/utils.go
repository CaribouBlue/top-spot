package utils

import (
	"encoding/json"
	"net/http"

	"github.com/CaribouBlue/top-spot/internal/model"
	"github.com/CaribouBlue/top-spot/internal/spotify"
	"github.com/a-h/templ"
)

func HandleJsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}

func HandleHtmlResponse(r *http.Request, w http.ResponseWriter, component templ.Component) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	component.Render(r.Context(), w)
}

func AuthorizedSpotifyClient(user *model.UserModel) *spotify.Client {
	spotify := spotify.DefaultClient()
	spotify.SetAccessToken(user.Data.SpotifyAccessToken)
	return spotify
}
