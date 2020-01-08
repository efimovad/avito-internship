package item_handler

import (
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/efimovad/avito-internship/internal/app/general"
	"github.com/efimovad/avito-internship/internal/app/item"
	"github.com/efimovad/avito-internship/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	MemCachedService = "memcached:11211"
)

type Handler struct {
	usecase			item.Usecase
	sanitizer		*bluemonday.Policy
	sessionStore	sessions.Store
	memcacheClient	*memcache.Client
}

func NewItemHandler(m *mux.Router, ucase item.Usecase, sessionStore sessions.Store, sanitizer *bluemonday.Policy) {
	handler := &Handler{
		usecase:   	  ucase,
		sanitizer:    sanitizer,
		sessionStore: sessionStore,
		memcacheClient: nil,
	}

	m.HandleFunc("/item", handler.CreateItem).Methods(http.MethodPost)
	m.HandleFunc("/item/{id:[0-9]+}", handler.GetItem).Methods(http.MethodGet)
	m.HandleFunc("/items", handler.GetItems).Methods(http.MethodGet)

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
	myItem.Images = make([]string, 0)
	if err := myItem.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "CreateItem<-myItem.UnmarshalJSON()")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	v := validator.New()
	if err := v.Struct(myItem); err != nil {
		err := errors.Wrap(err, "CreateItem<-validator.Struct()")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	myItem.Sanitize(h.sanitizer)

	if err := h.usecase.Create(myItem); err != nil {
		err := errors.Wrap(err, "CreateItem<-usecase.Create()")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	general.Respond(w, r, http.StatusCreated, myItem.ID)
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "GetItem<-Atoi(wrong id)")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	var allInfo bool
	fields := r.URL.Query().Get("fields")
	if fields != "" {
		allInfo, err = strconv.ParseBool(fields)
		if err != nil {
			general.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "wrong fields param"))
			return
		}
	}

	if cachedItem, err := h.getFromCache(id, allInfo); err == nil {
		general.Respond(w, r, http.StatusOK, cachedItem)
		return
	}

	myItem, err := h.usecase.Get(id, allInfo)
	if err != nil {
		err = errors.Wrapf(err, "GetItem<-usecase.Get()")
		general.Error(w, r, http.StatusNotFound, err)
		return
	}

	_ = h.setToCache(myItem, id, allInfo)

	general.Respond(w, r, http.StatusOK, myItem)
}

func (h * Handler) setToCache(myItem *model.Item, id int64, allInfo bool) error {
	if h.memcacheClient == nil {
		h.memcacheClient = memcache.New(MemCachedService)
	}

	jsonItem, err := json.Marshal(myItem)
	if err != nil {
		return errors.Wrap(err, "json.Marshal()")
	}

	cachedItem := &memcache.Item{
		Key:        "item" + strconv.FormatInt(id, 10) + strconv.FormatBool(allInfo),
		Value:      jsonItem,
		Expiration: 20,
	}

	err = h.memcacheClient.Set(cachedItem)
	if err != nil {
		return errors.Wrap(err, "memcacheClient.Set()")
	}

	return nil
}

func (h *Handler) getFromCache(id int64, allInfo bool) (*model.Item, error) {
	if h.memcacheClient == nil {
		h.memcacheClient = memcache.New(MemCachedService)
	}

	it, err := h.memcacheClient.Get("item" + strconv.FormatInt(id, 10) + strconv.FormatBool(allInfo))
	if err != nil {
		return nil, errors.Wrap(err, "memcacheClient.Get()")
	}

	cachedItem := &model.Item{}
	if err := json.Unmarshal(it.Value, &cachedItem); err != nil {
		return nil, errors.Wrap(err, "unmarshal json")
	}
	return cachedItem, nil
}


func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params model.Params
	sort := r.URL.Query().Get("sort")
	if sort == "date" {
		params.Date = true
	} else if sort == "price" {
		params.Price = true
	} else if sort != "" {
		general.Error(w, r, http.StatusBadRequest, errors.New("wrong sort param"))
		return
	}

	desc := r.URL.Query().Get("desc")
	if desc != "" {
		isDesc, err := strconv.ParseBool(desc)
		if err == nil {
			params.Desc = isDesc
		} else {
			general.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "wrong desc param"))
			return
		}
	}

	page := r.URL.Query().Get("page")
	if page != "" {
		pageInt, err := strconv.ParseInt(page, 10, 64)
		if err == nil {
			params.Page = pageInt
		} else {
			general.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "wrong page param"))
			return
		}
	}

	items, err := h.usecase.List(params)
	if err != nil {
		err = errors.Wrapf(err, "GetItems<-usecase.List()")
		general.Error(w, r, http.StatusNotFound, err)
		return
	}
	general.Respond(w, r, http.StatusOK, items)
}
