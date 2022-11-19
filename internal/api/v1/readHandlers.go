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

var ErrUserIsNotFound = errors.New("user is not found")

func (a *API) searchUsers(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.searchUsers START")
	defer log.Debug().Msg("api.searchUsers END")

	allUsers, err := a.storage.GetAllUsers()
	if err != nil {
		_ = render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if len(allUsers.List) == 0 {
		_ = render.Render(w, r, ErrNoContent(errors.New("there are no users")))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, allUsers)
}

func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.getUser START")
	defer log.Debug().Msg("api.getUser END")

	id := chi.URLParam(r, "id")

	idForFind, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("bad id")))
		return
	}

	user, err := a.storage.GetUser(idForFind)
	if err != nil {
		if errors.Is(err, storageErr.ErrUserIsNotFound) {
			_ = render.Render(w, r, ErrInvalidRequest(ErrUserIsNotFound))
		} else {
			_ = render.Render(w, r, ErrInternalServerError(err))
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}
