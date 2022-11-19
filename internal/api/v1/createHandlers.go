package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"refactoring/internal/model"
)

type createUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *createUserRequest) Bind(r *http.Request) error { return nil }

func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.createUser START")
	defer log.Debug().Msg("api.createUser END")

	reqUser := createUserRequest{}
	if err := render.Bind(r, &reqUser); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	newUser := model.User{CreatedAt: time.Now(), DisplayName: reqUser.DisplayName, Email: reqUser.Email}
	id, err := a.storage.CreateUser(&newUser)
	if err != nil {
		_ = render.Render(w, r, ErrInternalServerError(err))
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{"user_id": id})
}
