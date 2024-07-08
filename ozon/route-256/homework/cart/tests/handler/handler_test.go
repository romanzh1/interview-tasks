package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"route256/cart/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

func (suite *TestSuite) TestDeleteItemFromUserCart() {
	ctx := context.Background()

	userID := int64(1)
	skuID := int64(1001)

	req := models.CartRequest{
		UserID: userID,
		SkuID:  skuID,
		Count:  1,
	}

	err := suite.repo.AddItemToUserCart(ctx, req)
	assert.NoError(suite.T(), err)

	url := suite.server.URL + "/user/" + strconv.FormatInt(userID, 10) + "/cart/" + strconv.FormatInt(skuID, 10)
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	assert.NoError(suite.T(), err)

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close() //nolint
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNoContent, response.StatusCode)

	items, err := suite.repo.ListUserCart(ctx, 1)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), items)
}

func (suite *TestSuite) TestListUserCart() {
	ctx := context.Background()

	req := []models.CartRequest{
		{UserID: 1, SkuID: 2958025, Count: 2},
		{UserID: 1, SkuID: 773297411, Count: 1},
	}

	err := suite.repo.AddItemToUserCart(ctx, req[0])
	assert.NoError(suite.T(), err)
	err = suite.repo.AddItemToUserCart(ctx, req[1])
	assert.NoError(suite.T(), err)

	url := suite.server.URL + "/user/" + strconv.FormatInt(1, 10) + "/cart"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	assert.NoError(suite.T(), err)

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close() //nolint

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.StatusCode)

	var items struct {
		Items      []models.CartItem `json:"items"`
		TotalPrice uint32            `json:"total_price"`
	}
	expectedItems := struct {
		Items      []models.CartItem `json:"items"`
		TotalPrice uint32            `json:"total_price"`
	}{
		Items: []models.CartItem{
			{SkuID: req[0].SkuID, Count: req[0].Count, Name: "Roxy Music. Stranded. Remastered Edition", Price: 1028},
			{SkuID: req[1].SkuID, Count: req[1].Count, Name: "Кроссовки Nike JORDAN", Price: 2202},
		},
		TotalPrice: 4258,
	}

	err = json.NewDecoder(response.Body).Decode(&items)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), expectedItems, items)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
