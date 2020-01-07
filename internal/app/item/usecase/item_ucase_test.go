package item_ucase

import (
	"github.com/efimovad/avito-internship/internal/app/item"
	"github.com/efimovad/avito-internship/internal/model"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func testUcase(t *testing.T) (*MockItemRepository, item.Usecase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := NewMockItemRepository(ctrl)
	ucase := NewItemUsecase(repo)
	return repo, ucase
}

func compare(l model.Item, r model.Item, t *testing.T) bool {
	t.Helper()
	return l.Title == r.Title && l.Price == r.Price && l.MainImage == r.MainImage
}

func compareAll(l model.Item, r model.Item, t *testing.T) bool {
	t.Helper()
	return l.Title == r.Title &&
		l.Price == r.Price &&
		l.Images[0] == r.Images[0] &&
		l.Images[1] == r.Images[1] &&
		l.Images[2] == r.Images[2] &&
		l.Description == r.Description
}

func TestUsecase_Create(t *testing.T) {
	repo, ucase := testUcase(t)

	testItem := &model.Item{
		Title:       "title",
		Description: "description",
		Price:       100,
		Images:      []string{"image/1","image/2","image/3"},
	}

	testCases := []struct {
		name		string
		item		*model.Item
		expectError	error
	}{
		{
			name: "valid",
			item: testItem,
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo.
				EXPECT().
				Create(testItem).Do(func(arg *model.Item) {
				arg.ID = 1
			}).
				Return(tc.expectError)
			err := ucase.Create(testItem)

			if err != tc.expectError {
				t.Fatal("actual error:", err , "and expected:", tc.expectError, "do not match")
				return
			}
		})
	}
}

func TestUsecase_Get(t *testing.T) {
	repo, ucase := testUcase(t)

	testItem := &model.Item{
		ID:			 1,
		Title:       "title",
		Description: "description",
		Price:       100,
		Images:      []string{"image/1","image/2","image/3"},
	}

	testCases := []struct {
		name          	string
		item 			*model.Item
		allInfo			bool
		expectError		error
	}{
		{
			name:			"valid",
			item:			testItem,
			expectError:	nil,
			allInfo:		false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo.
				EXPECT().
				Get(tc.item.ID).
				Return(testItem, tc.expectError)
			myItem, err := ucase.Get(tc.item.ID, tc.allInfo)

			if tc.expectError != err {
				t.Fatal()
			}

			if tc.allInfo && err == nil {
				assert.Equal(t, compareAll(*myItem, *tc.item, t), true)
			} else if err == nil {
				assert.Equal(t, compare(*myItem, *tc.item, t), true)
			}
		})
	}
}

func TestUsecase_List(t *testing.T) {
	repo, ucase := testUcase(t)
	dateDesc, priceDesc, dateAsc, priceAsc := itemLists(t)

	testCases := []struct {
		name          	string
		list 			[]model.Item
		params			model.Params
		expectError		error
	}{
		{
			name:			"date desc",
			list:			dateDesc,
			expectError:	nil,
			params:			model.Params{
				Date:  true,
				Price: false,
				Desc:  true,
				Page:  0,
			},
		},
		{
			name:			"price desc",
			list:			priceDesc,
			expectError:	nil,
			params:			model.Params{
				Date:  false,
				Price: true,
				Desc:  true,
				Page:  0,
			},
		},
		{
			name:			"date asc",
			list:			dateAsc,
			expectError:	nil,
			params:			model.Params{
				Date:  true,
				Price: false,
				Desc:  false,
				Page:  0,
			},
		},
		{
			name:			"price asc",
			list:			priceAsc,
			expectError:	nil,
			params:			model.Params{
				Date:  false,
				Price: true,
				Desc:  false,
				Page:  0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo.
				EXPECT().
				List(tc.params).
				Return(tc.list, tc.expectError)

			list, err := ucase.List(tc.params)

			if tc.expectError != err {
				t.Fatal()
			}

			for i, item := range list {
				assert.Equal(t, compare(item, tc.list[i], t), true)
			}
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