package http

import (
	"github.com/efimovad/avito-internship/internal/app/general"
	"github.com/efimovad/avito-internship/internal/app/item"
	"github.com/efimovad/avito-internship/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	usecase		 item.Usecase
	sanitizer    *bluemonday.Policy
	logger       *zap.SugaredLogger
	sessionStore sessions.Store
}

func NewItemHandler(m *mux.Router, ucase item.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &Handler{
		usecase:   	  ucase,
		sanitizer:    sanitizer,
		logger:       logger,
		sessionStore: sessionStore,
	}

	m.HandleFunc("/create/item", handler.CreateItem).Methods(http.MethodPost)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "CreateItem<-r.Body.Close()")
			general.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "CreateItem<-ioutil.ReadAll()")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	myItem := new(model.Item)
	if err := myItem.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "CreateItem<-myItem.UnmarshalJSON()")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.usecase.Create(myItem); err != nil {
		err := errors.Wrap(err, "CreateItem<-usecase.Create()")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	general.Respond(w, r, http.StatusCreated, myItem.ID)
}
