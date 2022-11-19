package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	storageErr "refactoring/internal/storage/errors"
)

func (a *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.deleteUser START")
	defer log.Debug().Msg("api.deleteUser END")

	id := chi.URLParam(r, "id")
	idForFind, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("bad id")))
		return
	}

	err = a.storage.DeleteUser(idForFind)
	if err != nil {
		if errors.Is(err, storageErr.ErrUserIsNotFound) {
			_ = render.Render(w, r, ErrInvalidRequest(ErrUserIsNotFound))
		} else {
			_ = render.Render(w, r, ErrInternalServerError(err))
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.Data(w, r, []byte(""))
}
