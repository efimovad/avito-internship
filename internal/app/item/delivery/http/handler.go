package item_handler

import (
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/efimovad/avito-internship/internal/app/general"
	"github.com/efimovad/avito-internship/internal/app/item"
	"github.com/efimovad/avito-internship/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	//"time"
)

type Handler struct {
	usecase		 item.Usecase
	//sanitizer    *bluemonday.Policy
	//logger       *zap.SugaredLogger
	sessionStore sessions.Store
	memcacheClient *memcache.Client
}

//func NewItemHandler(m *mux.Router, ucase item.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
func NewItemHandler(m *mux.Router, ucase item.Usecase, sessionStore sessions.Store) {
	handler := &Handler{
		usecase:   	  ucase,
		//sanitizer:    sanitizer,
		//logger:       logger,
		sessionStore: sessionStore,
		memcacheClient: memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212"),
	}
	handler.memcacheClient.Timeout = time.Second * 10

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
			allInfo = false
		}
	}

	it, err := h.memcacheClient.Get("item")//+ strconv.FormatInt(id, 10))
	if err != nil {
		log.Println(it)
	}

	myItem, err := h.usecase.Get(id, allInfo)
	if err != nil {
		err = errors.Wrapf(err, "GetItem<-usecase.Get()")
		general.Error(w, r, http.StatusNotFound, err)
		return
	}

	text, _ := json.Marshal(myItem)
	memcacheItem := &memcache.Item{
		Key:        "item", //+ strconv.FormatInt(id, 10),
		Value:      text,
		Expiration: 20,
	}

	//h.memcacheClient.Timeout = time.Second * 10
	if err := h.memcacheClient.Set(memcacheItem); err != nil {
		log.Println(err)
	}

	general.Respond(w, r, http.StatusOK, myItem)
}


func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params model.Params
	sort := r.URL.Query().Get("sort")
	if sort == "date" {
		params.Date = true
	} else if sort == "price" {
		params.Price = true
	}

	desc := r.URL.Query().Get("desc")
	if desc != "" {
		isDesc, err := strconv.ParseBool(desc)
		if err == nil {
			params.Desc = isDesc
		}
	}

	page := r.URL.Query().Get("page")
	if page != "" {
		pageInt, err := strconv.ParseInt(page, 10, 64)
		if err == nil {
			params.Page = pageInt
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
