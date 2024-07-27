package cmd

import (
	"github.com/stretchr/testify/suite"
	"testing"
	item_suite "week-2-workshop/cart/test/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(item_suite.ItemS))
}
