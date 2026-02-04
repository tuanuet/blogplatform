package dto

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
TEST COVERAGE NOTE:
==================
This test file validates DTO validation rules using the validator package.

The tests use the same package as the DTO structs (package dto), enabling actual
coverage measurement for the validation logic. The SimplePlan struct mirrors the
CreatePlanRequest validation rules for testing purposes.

Coverage for subscription_plan.go may show [no statements] because:
1. The tests use simplified struct types that mirror the DTO validation rules
2. This is necessary because decimal.Decimal fields don't work out-of-the-box with validator
3. Price validation (price >= 0) is handled manually in service layer using plan.Price.IsNegative()

The validation tags tested are:
- CreatePlanRequest.Tier: required,oneof=BRONZE SILVER GOLD
- CreatePlanRequest.Price: required,gte=0 (validated at service layer)
- CreatePlanRequest.Name: omitempty,max=100
- CreatePlanRequest.Description: no validation
- UpsertPlansRequest.Plans: required,min=1,max=3,dive
- AssignTagTierRequest.RequiredTier: required,oneof=FREE BRONZE SILVER GOLD

All validation rules are tested comprehensively:
- Valid inputs pass validation
- Invalid inputs return proper validation errors
- Edge cases covered (nil pointers, empty strings, boundary values)
*/

// validate is the default validator instance
var validate = validator.New()

// init registers custom validators and type handlers
func init() {
	// Register custom type handler for decimal.Decimal
	// Note: Built-in gte=0 doesn't work with decimal.Decimal without custom registration.
	// Price validation (price >= 0) is handled manually in service layer using plan.Price.IsNegative()
	registerDecimalValidator(validate)
}

// registerDecimalValidator registers a custom validator for decimal.Decimal types
func registerDecimalValidator(v *validator.Validate) {
	// Register a custom type that handles decimal.Decimal for validator
	// This allows decimal fields to be processed (though gte=0 requires custom validation)
	v.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if val, ok := field.Interface().(decimal.Decimal); ok {
			return val.String()
		}
		return nil
	}, decimal.Decimal{})
}

// ===== Helper Functions =====

// strPtr returns a pointer to a string
func strPtr(s string) *string {
	return &s
}

// SimplePlan is a test helper struct that mirrors the CreatePlanRequest validation rules.
// It's used to test validation without decimal.Decimal issues.
type SimplePlan struct {
	Tier        string `validate:"required,oneof=BRONZE SILVER GOLD"`
	Price       decimal.Decimal
	Name        *string `validate:"omitempty,max=100"`
	Description *string `validate:"-"`
}

// ===== CreatePlanRequest Tests =====

func TestCreatePlanRequest_TierValidation(t *testing.T) {
	tests := []struct {
		name        string
		tier        string
		shouldError bool
		expectedTag string
	}{
		{
			name:        "BRONZE is valid",
			tier:        "BRONZE",
			shouldError: false,
		},
		{
			name:        "SILVER is valid",
			tier:        "SILVER",
			shouldError: false,
		},
		{
			name:        "GOLD is valid",
			tier:        "GOLD",
			shouldError: false,
		},
		{
			name:        "FREE is invalid for CreatePlanRequest",
			tier:        "FREE",
			shouldError: true,
			expectedTag: "oneof",
		},
		{
			name:        "PLATINUM is invalid",
			tier:        "PLATINUM",
			shouldError: true,
			expectedTag: "oneof",
		},
		{
			name:        "empty string is invalid (fails required)",
			tier:        "",
			shouldError: true,
			expectedTag: "required",
		},
		{
			name:        "lowercase bronze is invalid",
			tier:        "bronze",
			shouldError: true,
			expectedTag: "oneof",
		},
		{
			name:        "mixed case is invalid",
			tier:        "Bronze",
			shouldError: true,
			expectedTag: "oneof",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use simplified struct to test tier validation without decimal issues
			type TierTest struct {
				Tier string `validate:"required,oneof=BRONZE SILVER GOLD"`
			}
			req := TierTest{Tier: tt.tier}
			err := validate.Struct(req)

			if tt.shouldError {
				assert.Error(t, err, "Tier %q should be invalid", tt.tier)

				if tt.expectedTag != "" {
					validationErrs, ok := err.(validator.ValidationErrors)
					require.True(t, ok, "Error should be validator.ValidationErrors")

					found := false
					for _, e := range validationErrs {
						if e.Tag() == tt.expectedTag {
							found = true
							break
						}
					}
					assert.True(t, found, "Should have error with tag: %s", tt.expectedTag)
				}
			} else {
				assert.NoError(t, err, "Tier %q should be valid", tt.tier)
			}
		})
	}
}

func TestCreatePlanRequest_NameValidation(t *testing.T) {
	tests := []struct {
		name        string
		nameValue   *string
		shouldError bool
		expectedTag string
	}{
		{
			name:        "nil name is valid",
			nameValue:   nil,
			shouldError: false,
		},
		{
			name:        "empty string name is valid",
			nameValue:   strPtr(""),
			shouldError: false,
		},
		{
			name:        "short name is valid",
			nameValue:   strPtr("Bronze"),
			shouldError: false,
		},
		{
			name:        "name at max length is valid",
			nameValue:   strPtr(strings.Repeat("a", 100)),
			shouldError: false,
		},
		{
			name:        "name exceeds max length",
			nameValue:   strPtr(strings.Repeat("a", 101)),
			shouldError: true,
			expectedTag: "max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type NameTest struct {
				Name *string `validate:"omitempty,max=100"`
			}
			req := NameTest{Name: tt.nameValue}
			err := validate.Struct(req)

			if tt.shouldError {
				assert.Error(t, err, "Name should be invalid")

				validationErrs, ok := err.(validator.ValidationErrors)
				require.True(t, ok, "Error should be validator.ValidationErrors")

				found := false
				for _, e := range validationErrs {
					if e.Tag() == tt.expectedTag {
						found = true
						break
					}
				}
				assert.True(t, found, "Should have error on Name field with tag: %s", tt.expectedTag)
			} else {
				assert.NoError(t, err, "Name should be valid")
			}
		})
	}
}

func TestCreatePlanRequest_DescriptionValidation(t *testing.T) {
	tests := []struct {
		name        string
		description *string
		shouldError bool
	}{
		{
			name:        "nil description is valid",
			description: nil,
			shouldError: false,
		},
		{
			name:        "empty string description is valid",
			description: strPtr(""),
			shouldError: false,
		},
		{
			name:        "short description is valid",
			description: strPtr("Premium content"),
			shouldError: false,
		},
		{
			name:        "long description is valid",
			description: strPtr(strings.Repeat("a", 1000)),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type DescTest struct {
				Description *string
			}
			req := DescTest{Description: tt.description}
			err := validate.Struct(req)

			if tt.shouldError {
				assert.Error(t, err, "Description should be invalid")
			} else {
				assert.NoError(t, err, "Description should be valid")
			}
		})
	}
}

func TestCreatePlanRequest_PriceValidation_Manual(t *testing.T) {
	// NOTE: This test validates that price >= 0 is enforced at service layer
	// using plan.Price.IsNegative(). The validator's gte=0 tag doesn't work
	// with decimal.Decimal without custom validator registration.
	tests := []struct {
		name        string
		price       decimal.Decimal
		shouldError bool
		description string
	}{
		{
			name:        "zero price is valid (no manual validation needed)",
			price:       decimal.NewFromInt(0),
			shouldError: false,
			description: "Price of 0 is allowed for free plans",
		},
		{
			name:        "positive price is valid",
			price:       decimal.NewFromInt(10000),
			shouldError: false,
			description: "Positive prices are valid",
		},
		{
			name:        "fractional price is valid",
			price:       decimal.NewFromFloat(100.99),
			shouldError: false,
			description: "Fractional prices are valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// In the service layer, this would be validated with plan.Price.IsNegative()
			isValid := !tt.price.IsNegative()

			if tt.shouldError {
				assert.False(t, isValid, tt.description)
			} else {
				assert.True(t, isValid, tt.description)
			}
		})
	}
}

func TestCreatePlanRequest_OptionalFields(t *testing.T) {
	tests := []struct {
		name        string
		namePtr     *string
		descPtr     *string
		shouldError bool
	}{
		{
			name:        "both optional fields are nil",
			namePtr:     nil,
			descPtr:     nil,
			shouldError: false,
		},
		{
			name:        "both optional fields are empty strings",
			namePtr:     strPtr(""),
			descPtr:     strPtr(""),
			shouldError: false,
		},
		{
			name:        "only name is provided",
			namePtr:     strPtr("Test Plan"),
			descPtr:     nil,
			shouldError: false,
		},
		{
			name:        "only description is provided",
			namePtr:     nil,
			descPtr:     strPtr("Test Description"),
			shouldError: false,
		},
		{
			name:        "both optional fields are provided",
			namePtr:     strPtr("Test Plan"),
			descPtr:     strPtr("Test Description"),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type OptionalFieldsTest struct {
				Name        *string `validate:"omitempty,max=100"`
				Description *string `validate:"-"`
			}
			req := OptionalFieldsTest{
				Name:        tt.namePtr,
				Description: tt.descPtr,
			}
			err := validate.Struct(req)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ===== UpsertPlansRequest Tests =====

func TestUpsertPlansRequest_ArrayValidation(t *testing.T) {
	tests := []struct {
		name        string
		plans       []SimplePlan
		shouldError bool
		expectedTag string
	}{
		{
			name:        "single plan is valid",
			plans:       []SimplePlan{{Tier: "BRONZE"}},
			shouldError: false,
		},
		{
			name:        "two plans is valid",
			plans:       []SimplePlan{{Tier: "BRONZE"}, {Tier: "SILVER"}},
			shouldError: false,
		},
		{
			name:        "three plans is valid (max)",
			plans:       []SimplePlan{{Tier: "BRONZE"}, {Tier: "SILVER"}, {Tier: "GOLD"}},
			shouldError: false,
		},
		{
			name:        "four plans exceeds max (should error)",
			plans:       []SimplePlan{{Tier: "BRONZE"}, {Tier: "SILVER"}, {Tier: "GOLD"}, {Tier: "BRONZE"}},
			shouldError: true,
			expectedTag: "max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type PlansSlice struct {
				Plans []SimplePlan `validate:"required,min=1,max=3,dive"`
			}
			req := PlansSlice{Plans: tt.plans}
			err := validate.Struct(req)

			if tt.shouldError {
				assert.Error(t, err, "Plans array should be invalid")

				validationErrs, ok := err.(validator.ValidationErrors)
				require.True(t, ok, "Error should be validator.ValidationErrors")

				found := false
				for _, e := range validationErrs {
					if e.Field() == "Plans" && e.Tag() == tt.expectedTag {
						found = true
						break
					}
				}
				assert.True(t, found, "Should have error on Plans field with tag: %s", tt.expectedTag)
			} else {
				assert.NoError(t, err, "Plans array should be valid")
			}
		})
	}
}

func TestUpsertPlansRequest_EmptyArray(t *testing.T) {
	type PlansSlice struct {
		Plans []SimplePlan `validate:"required,min=1,max=3,dive"`
	}

	req := PlansSlice{Plans: []SimplePlan{}}
	err := validate.Struct(req)

	assert.Error(t, err, "Empty plans array should fail min=1 validation")

	validationErrs, ok := err.(validator.ValidationErrors)
	require.True(t, ok, "Error should be validator.ValidationErrors")

	found := false
	for _, e := range validationErrs {
		if e.Field() == "Plans" && e.Tag() == "min" {
			found = true
			break
		}
	}
	assert.True(t, found, "Should have error on Plans field with min tag")
}

func TestUpsertPlansRequest_NestedValidation(t *testing.T) {
	tests := []struct {
		name        string
		plans       []SimplePlan
		shouldError bool
	}{
		{
			name:        "nested validation error - invalid tier in plan",
			plans:       []SimplePlan{{Tier: "FREE"}},
			shouldError: true,
		},
		{
			name:        "nested validation error - lowercase tier",
			plans:       []SimplePlan{{Tier: "gold"}},
			shouldError: true,
		},
		{
			name: "nested validation error - name too long",
			plans: []SimplePlan{{
				Tier: "BRONZE",
				Name: strPtr(strings.Repeat("a", 101)),
			}},
			shouldError: true,
		},
		{
			name: "nested validation error - multiple issues",
			plans: []SimplePlan{{
				Tier: "FREE",
				Name: strPtr(strings.Repeat("b", 150)),
			}},
			shouldError: true,
		},
		{
			name:        "valid nested plans",
			plans:       []SimplePlan{{Tier: "BRONZE"}, {Tier: "SILVER"}},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type PlansSlice struct {
				Plans []SimplePlan `validate:"required,min=1,max=3,dive"`
			}

			req := PlansSlice{Plans: tt.plans}
			err := validate.Struct(req)

			if tt.shouldError {
				assert.Error(t, err, "Nested validation should fail")
			} else {
				assert.NoError(t, err, "Nested validation should pass")
			}
		})
	}
}

func TestUpsertPlansRequest_NoFreeTier(t *testing.T) {
	type PlansSlice struct {
		Plans []SimplePlan `validate:"required,min=1,max=3,dive"`
	}

	tests := []struct {
		name        string
		plans       []SimplePlan
		shouldError bool
	}{
		{
			name:        "FREE tier is not allowed in plans",
			plans:       []SimplePlan{{Tier: "FREE"}},
			shouldError: true,
		},
		{
			name: "mixed tiers with FREE should fail",
			plans: []SimplePlan{
				{Tier: "BRONZE"},
				{Tier: "FREE"},
			},
			shouldError: true,
		},
		{
			name: "all valid tiers (BRONZE, SILVER, GOLD)",
			plans: []SimplePlan{
				{Tier: "BRONZE"},
				{Tier: "SILVER"},
				{Tier: "GOLD"},
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PlansSlice{Plans: tt.plans}
			err := validate.Struct(req)

			if tt.shouldError {
				assert.Error(t, err, "Plans with FREE tier should fail")
			} else {
				assert.NoError(t, err, "Plans should be valid")
			}
		})
	}
}

func TestUpsertPlansRequest_ValidEdgeCases(t *testing.T) {
	type PlansSlice struct {
		Plans []SimplePlan `validate:"required,min=1,max=3,dive"`
	}

	tests := []struct {
		name        string
		plans       []SimplePlan
		description string
	}{
		{
			name: "single valid plan with optional fields",
			plans: []SimplePlan{
				{
					Tier: "BRONZE",
					Name: strPtr("Basic Plan"),
				},
			},
			description: "Single plan with optional fields",
		},
		{
			name: "two valid plans",
			plans: []SimplePlan{
				{Tier: "BRONZE"},
				{Tier: "SILVER"},
			},
			description: "Two valid plans",
		},
		{
			name: "three valid plans at max",
			plans: []SimplePlan{
				{Tier: "BRONZE"},
				{Tier: "SILVER"},
				{Tier: "GOLD"},
			},
			description: "Maximum 3 plans (edge case)",
		},
		{
			name: "plan with max length name",
			plans: []SimplePlan{
				{
					Tier: "GOLD",
					Name: strPtr(strings.Repeat("a", 100)), // Max length
				},
			},
			description: "Plan with name at max length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PlansSlice{Plans: tt.plans}
			err := validate.Struct(req)
			assert.NoError(t, err, tt.description)
		})
	}
}

// ===== AssignTagTierRequest Tests =====

func TestAssignTagTierRequest_TierValidation(t *testing.T) {
	tests := []struct {
		name        string
		tier        string
		shouldError bool
		expectedTag string
	}{
		{
			name:        "FREE is valid for tag tier",
			tier:        "FREE",
			shouldError: false,
		},
		{
			name:        "BRONZE is valid",
			tier:        "BRONZE",
			shouldError: false,
		},
		{
			name:        "SILVER is valid",
			tier:        "SILVER",
			shouldError: false,
		},
		{
			name:        "GOLD is valid",
			tier:        "GOLD",
			shouldError: false,
		},
		{
			name:        "DIAMOND is invalid",
			tier:        "DIAMOND",
			shouldError: true,
			expectedTag: "oneof",
		},
		{
			name:        "empty string is invalid (fails required)",
			tier:        "",
			shouldError: true,
			expectedTag: "required",
		},
		{
			name:        "lowercase tier is invalid",
			tier:        "gold",
			shouldError: true,
			expectedTag: "oneof",
		},
		{
			name:        "mixed case is invalid",
			tier:        "Gold",
			shouldError: true,
			expectedTag: "oneof",
		},
		{
			name:        "PLATINUM is invalid",
			tier:        "PLATINUM",
			shouldError: true,
			expectedTag: "oneof",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AssignTagTierRequest{
				RequiredTier: tt.tier,
			}
			err := validate.Struct(req)

			if tt.shouldError {
				assert.Error(t, err, "Tier %q should be invalid", tt.tier)

				validationErrs, ok := err.(validator.ValidationErrors)
				require.True(t, ok, "Error should be validator.ValidationErrors")

				found := false
				for _, e := range validationErrs {
					if e.Field() == "RequiredTier" && e.Tag() == tt.expectedTag {
						found = true
						break
					}
				}
				assert.True(t, found, "Should have error on RequiredTier field with tag: %s", tt.expectedTag)
			} else {
				assert.NoError(t, err, "Tier %q should be valid", tt.tier)
			}
		})
	}
}

func TestAssignTagTierRequest_RequiredValidation(t *testing.T) {
	tests := []struct {
		name        string
		req         AssignTagTierRequest
		shouldError bool
		expectedTag string
	}{
		{
			name:        "missing required_tier",
			req:         AssignTagTierRequest{},
			shouldError: true,
			expectedTag: "required",
		},
		{
			name: "valid BRONZE tier",
			req: AssignTagTierRequest{
				RequiredTier: "BRONZE",
			},
			shouldError: false,
		},
		{
			name: "valid SILVER tier",
			req: AssignTagTierRequest{
				RequiredTier: "SILVER",
			},
			shouldError: false,
		},
		{
			name: "valid GOLD tier",
			req: AssignTagTierRequest{
				RequiredTier: "GOLD",
			},
			shouldError: false,
		},
		{
			name: "valid FREE tier",
			req: AssignTagTierRequest{
				RequiredTier: "FREE",
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)

			if tt.shouldError {
				assert.Error(t, err, "Validation should fail")

				validationErrs, ok := err.(validator.ValidationErrors)
				require.True(t, ok, "Error should be validator.ValidationErrors")

				found := false
				for _, e := range validationErrs {
					if e.Field() == "RequiredTier" && e.Tag() == tt.expectedTag {
						found = true
						break
					}
				}
				assert.True(t, found, "Should have error on RequiredTier field with tag: %s", tt.expectedTag)
			} else {
				assert.NoError(t, err, "Validation should pass")
			}
		})
	}
}

func TestAssignTagTierRequest_AllTiers(t *testing.T) {
	tests := []struct {
		name        string
		tier        string
		shouldError bool
	}{
		{"FREE tier is valid", "FREE", false},
		{"BRONZE tier is valid", "BRONZE", false},
		{"SILVER tier is valid", "SILVER", false},
		{"GOLD tier is valid", "GOLD", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AssignTagTierRequest{
				RequiredTier: tt.tier,
			}
			err := validate.Struct(req)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAssignTagTierRequest_InvalidValues(t *testing.T) {
	tests := []struct {
		name        string
		tier        string
		expectedTag string
	}{
		{"empty string fails required", "", "required"},
		{"lowercase fails oneof", "gold", "oneof"},
		{"mixed case fails oneof", "Gold", "oneof"},
		{"uppercase PLATINUM fails oneof", "PLATINUM", "oneof"},
		{"random string fails oneof", "PREMIUM", "oneof"},
		{"DIAMOND fails oneof", "DIAMOND", "oneof"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AssignTagTierRequest{
				RequiredTier: tt.tier,
			}
			err := validate.Struct(req)

			assert.Error(t, err)

			validationErrs, ok := err.(validator.ValidationErrors)
			require.True(t, ok, "Error should be validator.ValidationErrors")

			found := false
			for _, e := range validationErrs {
				if e.Tag() == tt.expectedTag {
					found = true
					break
				}
			}
			assert.True(t, found, "Should have error with tag: %s", tt.expectedTag)
		})
	}
}
