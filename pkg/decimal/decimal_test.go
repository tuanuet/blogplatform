package decimal

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDecimal_BasicOperations(t *testing.T) {
	t.Run("should create decimal from string", func(t *testing.T) {
		// Arrange
		value := "10.50"

		// Act
		dec, err := decimal.NewFromString(value)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "10.5", dec.String())
	})

	t.Run("should add two decimals", func(t *testing.T) {
		// Arrange
		a := decimal.NewFromInt(10)
		b := decimal.NewFromFloat(5.5)

		// Act
		result := a.Add(b)

		// Assert
		assert.Equal(t, "15.5", result.String())
	})

	t.Run("should compare decimals", func(t *testing.T) {
		// Arrange
		smaller := decimal.NewFromInt(10)
		larger := decimal.NewFromInt(20)

		// Act & Assert
		assert.True(t, smaller.LessThan(larger))
		assert.True(t, larger.GreaterThan(smaller))
	})
}
