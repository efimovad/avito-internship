package item_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/efimovad/avito-internship/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const SessionKey  = "test-handler-session"

func TestHandler_CreateItem(t *testing.T) {
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "correct item",
			payload: map[string]interface{}{
				"title":                    "new item",
				"price":     				100,
				"description":		        "description",
				"images":		            []string{"image/1", "image/2", "image/3"},
			},
			expectedCode: 201,
		},
		{
			name: "incorrect item: empty title",
			payload: map[string]interface{}{
				"price":     				100,
				"description":		        "description",
				"images":		            []string{"image/1", "image/2", "image/3"},
			},
			expectedCode: 400,
		},
		{
			name: "incorrect item: long title",
			payload: map[string]interface{}{
				"title":					"dghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsf" +
					"gdghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgdg" +
					"hfgdgsfgdghfgdgsfgdghfgdgsfgdghfgdgsfgaaa",
				"price":     				100,
				"description":		        "description",
				"images":		            []string{"image/1", "image/2", "image/3"},
			},
			expectedCode: 400,
		},
		{
			name: "incorrect item: too many images",
			payload: map[string]interface{}{
				"title":					"my title",
				"price":     				100,
				"description":		        "description",
				"images":		            []string{"image/1", "image/2", "image/3", "image/4"},
			},
			expectedCode: 400,
		},
	}

	ctrl := gomock.NewController(t)
	ucase := NewMockItemUsecase(ctrl)

	router := mux.NewRouter()
	sessionStore := sessions.NewCookieStore([]byte(SessionKey))
	NewItemHandler(router, ucase, sessionStore)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal(err)
			}

			var newItem *model.Item
			err := mapstructure.Decode(tc.payload, &newItem)
			if err != nil {
				t.Fatal(err)
			}

			if tc.expectedCode == 201 {
				ucase.EXPECT().Create(newItem).Do(func(arg *model.Item) {
					arg.ID = 1
				}).Return(nil)
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/item", b)

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestHandler_GetItem(t *testing.T) {
	testItem := &model.Item{
		Title:       "title",
		Description: "description",
		Price:       100,
		Images:      []string{"image/1","image/2","image/3"},
	}

	testCases := []struct {
		name         string
		id			 int64
		allInfo		 bool
		expectedItem *model.Item
		expectedCode int
		expectedErr	 error
	}{
		{
			id:	  1,
			allInfo: false,
			name: "existed item: part info",
			expectedCode: 200,
			expectedErr:  nil,
			expectedItem: testItem,
		},
		{
			id:	  1,
			allInfo:true,
			name: "existed item: all info",
			expectedCode: 200,
			expectedErr: nil,
			expectedItem: testItem,
		},
		{
			id:	  99,
			allInfo: true,
			name: "not existed item",
			expectedCode: 404,
			expectedErr:  errors.New("not found"),
			expectedItem: nil,
		},
	}

	ctrl := gomock.NewController(t)
	ucase := NewMockItemUsecase(ctrl)

	router := mux.NewRouter()
	sessionStore := sessions.NewCookieStore([]byte(SessionKey))
	NewItemHandler(router, ucase, sessionStore)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ucase.EXPECT().
				Get(tc.id, tc.allInfo).
				Return(tc.expectedItem, tc.expectedErr)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet,
				"/item/" + strconv.FormatInt(tc.id, 10) + "?fields=" + strconv.FormatBool(tc.allInfo),
				nil)

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestHandler_GetItems(t *testing.T) {
	//dateDesc, priceDesc, dateAsc, priceAsc := itemLists(t)
	_, _, dateAsc, _ := itemLists(t)

	testCases := []struct {
		name         string
		params		 model.Params
		expectedCode int
		expectedErr	 error
	}{
		{
			name: "correct params",
			params:	model.Params{
				Date:  true,
				Price: false,
				Desc:  false,
				Page:  0,
			},
			expectedCode: 200,
			expectedErr:  nil,
		},
		/*{
			name: "incorrect params",
			params:	model.Params{
				Date:  true,
				Price: true,
				Desc:  false,
				Page:  0,
			},
			expectedCode: 400,
			expectedErr: nil,
		},*/
	}

	ctrl := gomock.NewController(t)
	ucase := NewMockItemUsecase(ctrl)

	router := mux.NewRouter()
	sessionStore := sessions.NewCookieStore([]byte(SessionKey))
	NewItemHandler(router, ucase, sessionStore)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ucase.EXPECT().
				List(tc.params).
				Return(dateAsc, tc.expectedErr)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet,
				"/items" + "?sort=date&desc=" + strconv.FormatBool(tc.params.Desc),
				nil)

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func itemLists(t *testing.T) ([]model.Item, []model.Item, []model.Item, []model.Item){
	t.Helper()

	item1 := model.Item{
		Title:       "item1",
		Price:       100,
		MainImage:   "image/1",
	}

	item2 := model.Item{
		Title:       "item2",
		Price:       200,
		MainImage:   "image/2",
	}

	item3 := model.Item{
		Title:       "item3",
		Price:       50,
		MainImage:   "image/3",
	}

	return []model.Item{item3, item2, item1}, []model.Item{item2, item1, item3}, []model.Item{item1, item2, item3}, []model.Item{item3, item1, item2}
}