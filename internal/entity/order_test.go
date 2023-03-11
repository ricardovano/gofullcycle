package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfAGetAnErrorIfIdIsBlank(t *testing.T) {
	order := Order{}
	assert.Error(t, order.Validate(), "invalid id")
}

func TestWithAllValidParams(t *testing.T) {
	order := Order{Id: "100A", Price: 10.0, Tax: 1.1}
	order.CalculateFinalPrice()

	assert.NoError(t, order.Validate())
	assert.Equal(t, "100A", order.Id)
	assert.Equal(t, 10.0, order.Price)
	assert.Equal(t, 1.1, order.Tax)
	assert.Equal(t, 11.0, order.FinalPrice)
}
