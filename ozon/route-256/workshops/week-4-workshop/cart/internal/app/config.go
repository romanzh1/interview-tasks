package app

import (
	"fmt"

	"week-4-workshop/cart/internal/app/defenitions"
)

type (
	Options struct {
		Addr, ProductToken, ProductAddr, DbConnStr string
	}

	configProductService struct {
		productToken, productAddr string
	}
	path struct {
		index, cartItemAdd, cartItemDelete, cartDelete, cartList string
	}
	config struct {
		addr      string
		dbConnStr string
		configProductService
		path path
	}
)

func NewConfig(opts Options) config {
	return config{
		addr: opts.Addr,
		configProductService: configProductService{
			productToken: opts.ProductToken,
			productAddr:  opts.ProductAddr,
		},
		dbConnStr: opts.DbConnStr,
		path: path{
			index:       "/",
			cartItemAdd: fmt.Sprintf("POST /user/{%s}/cart/{%s}", defenitions.ParamUserID, defenitions.ParamSkuID),
		},
	}
}
