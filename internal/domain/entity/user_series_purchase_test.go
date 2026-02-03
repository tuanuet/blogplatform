package entity_test

import (
	"testing"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// TestUserSeriesPurchase_TableName verifies the table name is correctly set
func TestUserSeriesPurchase_TableName(t *testing.T) {
	// Arrange
	usp := entity.UserSeriesPurchase{}

	// Act
	tableName := usp.TableName()

	// Assert
	assert.Equal(t, "user_series_purchases", tableName, "TableName should return 'user_series_purchases'")
}

// TestUserSeriesPurchase_CompositePrimaryKey verifies the composite primary key structure
func TestUserSeriesPurchase_CompositePrimaryKey(t *testing.T) {
	// Arrange
	userID := uuid.New()
	seriesID := uuid.New()

	// Act
	usp := entity.UserSeriesPurchase{
		UserID:   userID,
		SeriesID: seriesID,
	}

	// Assert
	assert.Equal(t, userID, usp.UserID, "UserID should be set")
	assert.Equal(t, seriesID, usp.SeriesID, "SeriesID should be set")
}

// TestUserSeriesPurchase_AmountField verifies the amount field is correctly typed
func TestUserSeriesPurchase_AmountField(t *testing.T) {
	// Arrange
	amount := decimal.NewFromInt(50000)

	// Act
	usp := entity.UserSeriesPurchase{
		Amount: amount,
	}

	// Assert
	assert.True(t, usp.Amount.Equal(amount), "Amount should be set correctly")
	assert.Equal(t, "50000", usp.Amount.String(), "Amount value should match")
}

// TestUserSeriesPurchase_CreatedAtField verifies the createdAt field
func TestUserSeriesPurchase_CreatedAtField(t *testing.T) {
	// Arrange
	now := time.Now()

	// Act
	usp := entity.UserSeriesPurchase{
		CreatedAt: now,
	}

	// Assert
	assert.Equal(t, now, usp.CreatedAt, "CreatedAt should be set")
}

// TestUserSeriesPurchase_UserRelationship verifies the user relationship
func TestUserSeriesPurchase_UserRelationship(t *testing.T) {
	// Arrange
	userID := uuid.New()
	user := &entity.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Act
	usp := entity.UserSeriesPurchase{
		UserID: userID,
		User:   user,
	}

	// Assert
	assert.NotNil(t, usp.User, "User relationship should be set")
	assert.Equal(t, userID, usp.User.ID, "User ID should match")
	assert.Equal(t, "John Doe", usp.User.Name, "User Name should match")
}

// TestUserSeriesPurchase_SeriesRelationship verifies the series relationship
func TestUserSeriesPurchase_SeriesRelationship(t *testing.T) {
	// Arrange
	seriesID := uuid.New()
	series := &entity.Series{
		ID:    seriesID,
		Title: "Golang for Beginners",
	}

	// Act
	usp := entity.UserSeriesPurchase{
		SeriesID: seriesID,
		Series:   series,
	}

	// Assert
	assert.NotNil(t, usp.Series, "Series relationship should be set")
	assert.Equal(t, seriesID, usp.Series.ID, "Series ID should match")
	assert.Equal(t, "Golang for Beginners", usp.Series.Title, "Series Title should match")
}

// TestUserSeriesPurchase_FullEntity verifies all fields together
func TestUserSeriesPurchase_FullEntity(t *testing.T) {
	// Arrange
	userID := uuid.New()
	seriesID := uuid.New()
	amount := decimal.NewFromInt(100000)
	now := time.Now()
	user := &entity.User{
		ID:    userID,
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}
	series := &entity.Series{
		ID:    seriesID,
		Title: "Advanced Go Patterns",
	}

	// Act
	usp := entity.UserSeriesPurchase{
		UserID:    userID,
		SeriesID:  seriesID,
		Amount:    amount,
		CreatedAt: now,
		User:      user,
		Series:    series,
	}

	// Assert
	assert.Equal(t, userID, usp.UserID, "UserID should be set")
	assert.Equal(t, seriesID, usp.SeriesID, "SeriesID should be set")
	assert.True(t, usp.Amount.Equal(amount), "Amount should be set correctly")
	assert.Equal(t, now, usp.CreatedAt, "CreatedAt should be set")
	assert.NotNil(t, usp.User, "User relationship should be set")
	assert.Equal(t, userID, usp.User.ID, "User ID in relationship should match")
	assert.NotNil(t, usp.Series, "Series relationship should be set")
	assert.Equal(t, seriesID, usp.Series.ID, "Series ID in relationship should match")
}
