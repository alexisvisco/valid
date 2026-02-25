package valid_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexisvisco/valid"
	"github.com/alexisvisco/valid/is"
	"testing"

	"github.com/goforj/godump"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- helpers ----------------------------------------------------------------

func alwaysFail(code string) is.Rule {
	return func(_ context.Context, _ any) *is.Violation {
		return &is.Violation{Code: is.ViolationCode(code)}
	}
}

func alwaysPass() is.Rule {
	return func(_ context.Context, _ any) *is.Violation { return nil }
}

// ---- Field ------------------------------------------------------------------

func TestField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		value     any
		rules     []is.Rule
		wantPaths []string
		wantCodes []string
	}{
		{
			name: "no rules",
			path: "X", value: "v", rules: nil,
		},
		{
			name: "one rule pass",
			path: "X", value: "v", rules: []is.Rule{alwaysPass()},
		},
		{
			name: "one rule fail",
			path: "X", value: "",
			rules:     []is.Rule{alwaysFail("ERR")},
			wantPaths: []string{"X"}, wantCodes: []string{"ERR"},
		},
		{
			name: "two rules both fail — short-circuit on first",
			path: "Name", value: "",
			rules:     []is.Rule{alwaysFail("A"), alwaysFail("B")},
			wantPaths: []string{"Name"}, wantCodes: []string{"A"},
		},
		{
			name: "required on empty string",
			path: "Email", value: "",
			rules:     []is.Rule{is.Required},
			wantPaths: []string{"Email"}, wantCodes: []string{string(is.ViolationRequired)},
		},
		{
			name: "required on non-empty string",
			path: "Email", value: "a@b.com",
			rules: []is.Rule{is.Required},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := valid.Field(tt.path, tt.value, tt.rules...)(context.Background())
			require.Len(t, got, len(tt.wantCodes), "Field(%q)", tt.path)
			for i, code := range tt.wantCodes {
				assert.Equal(t, code, got[i].Code, "Field error[%d].Code", i)
				assert.Equal(t, tt.wantPaths[i], got[i].Path, "Field error[%d].Path", i)
			}
		})
	}
}

// ---- Struct -----------------------------------------------------------------

func TestStruct(t *testing.T) {
	t.Parallel()

	t.Run("all groups empty → nil", func(t *testing.T) {
		t.Parallel()
		err := valid.Struct(context.Background())
		require.NoError(t, err)
	})

	t.Run("no groups → nil", func(t *testing.T) {
		t.Parallel()
		err := valid.Struct(context.Background())
		require.NoError(t, err)
	})

	t.Run("groups with errors → *Error flattened", func(t *testing.T) {
		t.Parallel()
		g1 := valid.Field("A", "", is.Required)
		g2 := valid.Field("B", 0, is.Required)
		err := valid.Struct(context.Background(), g1, g2)
		require.Error(t, err)
		ve := valid.As(err)
		require.NotNil(t, ve)
		require.Len(t, ve.Fields, 2)
		require.Equal(t, "A", ve.Fields[0].Path)
		require.Equal(t, "B", ve.Fields[1].Path)
	})

	t.Run("Error() message format", func(t *testing.T) {
		t.Parallel()
		err := valid.Struct(context.Background(), valid.Field("Name", "", is.Required))
		require.Error(t, err)
		msg := err.Error()
		require.NotEmpty(t, msg)
		assert.Contains(t, msg, "Name")
	})
}

// ---- Slice ------------------------------------------------------------------

func TestSlice(t *testing.T) {
	t.Parallel()

	t.Run("nil slice → nil", func(t *testing.T) {
		t.Parallel()
		var items []string
		got := valid.Slice("Items", items, func(ctx context.Context, i int, item string) error {
			return valid.Struct(ctx, valid.Field("Name", item, is.Required))
		})(context.Background())
		require.Nil(t, got)
	})

	t.Run("all pass → nil", func(t *testing.T) {
		t.Parallel()
		items := []string{"a", "b", "c"}
		got := valid.Slice("Tags", items, func(ctx context.Context, i int, item string) error {
			return valid.Struct(ctx, valid.Field("Value", item, is.Required))
		})(context.Background())
		require.Nil(t, got)
	})

	t.Run("index 0 fails → prefixed path", func(t *testing.T) {
		t.Parallel()
		items := []string{"", "ok"}
		got := valid.Slice("Items", items, func(ctx context.Context, i int, item string) error {
			return valid.Struct(ctx, valid.Field("Name", item, is.Required))
		})(context.Background())
		require.Len(t, got, 1)
		assert.Equal(t, "Items.0.Name", got[0].Path)
	})

	t.Run("index 2 fails → prefixed path", func(t *testing.T) {
		t.Parallel()
		items := []string{"a", "b", ""}
		got := valid.Slice("Items", items, func(ctx context.Context, i int, item string) error {
			return valid.Struct(ctx, valid.Field("Name", item, is.Required))
		})(context.Background())
		require.Len(t, got, 1)
		assert.Equal(t, "Items.2.Name", got[0].Path)
	})

	t.Run("fn returns non-valid error → path.i code=invalid", func(t *testing.T) {
		t.Parallel()
		items := []string{"x"}
		got := valid.Slice("Items", items, func(ctx context.Context, i int, item string) error {
			return errors.New("some internal error")
		})(context.Background())
		require.Len(t, got, 1)
		assert.Equal(t, "Items.0", got[0].Path)
		assert.Equal(t, "invalid", got[0].Code)
	})
}

// ---- Each -------------------------------------------------------------------

func TestEach(t *testing.T) {
	t.Parallel()

	t.Run("nil slice → nil", func(t *testing.T) {
		t.Parallel()
		var items []string
		got := valid.Each("Tags", items, is.Required)(context.Background())
		require.Nil(t, got)
	})

	t.Run("all pass → nil", func(t *testing.T) {
		t.Parallel()
		items := []string{"a", "b", "c"}
		got := valid.Each("Tags", items, is.Required)(context.Background())
		require.Nil(t, got)
	})

	t.Run("one element fails → indexed path", func(t *testing.T) {
		t.Parallel()
		items := []string{"card", "bank_transfer", "crypto"}
		got := valid.Each("Types", items, is.OneOf("card", "bank_transfer"))(context.Background())
		require.Len(t, got, 1)
		assert.Equal(t, "Types.2", got[0].Path)
		assert.Equal(t, string(is.ViolationOneOf), got[0].Code)
	})

	t.Run("multiple elements fail → all indexed paths", func(t *testing.T) {
		t.Parallel()
		items := []string{"card", "crypto", "cash"}
		got := valid.Each("Types", items, is.OneOf("card", "bank_transfer"))(context.Background())
		require.Len(t, got, 2)
		assert.Equal(t, "Types.1", got[0].Path)
		assert.Equal(t, "Types.2", got[1].Path)
	})

	t.Run("short-circuits per element on first rule failure", func(t *testing.T) {
		t.Parallel()
		items := []string{""}
		got := valid.Each("Tags", items, is.Required, alwaysFail("SECOND"))(context.Background())
		require.Len(t, got, 1)
		assert.Equal(t, string(is.ViolationRequired), got[0].Code)
	})
}

// ---- As ---------------------------------------------------------------------

func TestAs(t *testing.T) {
	t.Parallel()

	t.Run("nil error → nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, valid.As(nil))
	})

	t.Run("*Error direct → returned", func(t *testing.T) {
		t.Parallel()
		ve := &valid.Error{Fields: []valid.FieldError{{Path: "X", Code: "ERR"}}}
		got := valid.As(ve)
		require.NotNil(t, got)
		require.Same(t, ve, got)
	})

	t.Run("plain error → nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, valid.As(errors.New("plain")))
	})

	t.Run("wrapped *Error → returned", func(t *testing.T) {
		t.Parallel()
		ve := &valid.Error{Fields: []valid.FieldError{{Path: "X", Code: "ERR"}}}
		wrapped := fmt.Errorf("wrap: %w", ve)
		got := valid.As(wrapped)
		require.NotNil(t, got)
	})
}

// ---- Rename -----------------------------------------------------------------

func TestRename(t *testing.T) {
	t.Parallel()

	t.Run("nil error → nil", func(t *testing.T) {
		t.Parallel()
		var e *valid.Error
		require.Nil(t, e.Rename(map[string]string{"A": "a"}))
	})

	t.Run("exact match → translated", func(t *testing.T) {
		t.Parallel()
		ve := &valid.Error{Fields: []valid.FieldError{
			{Path: "Reference", Code: "VALIDATION_REQUIRED"},
		}}
		got := ve.Rename(map[string]string{"Reference": "reference"})
		require.NotNil(t, got)
		require.Len(t, got.Fields, 1)
		assert.Equal(t, "reference", got.Fields[0].Path)
	})

	t.Run("wildcard match → indices substituted", func(t *testing.T) {
		t.Parallel()
		ve := &valid.Error{Fields: []valid.FieldError{
			{Path: "Items.0.Name", Code: "VALIDATION_REQUIRED"},
			{Path: "Items.2.Quantity", Code: "VALIDATION_MIN"},
		}}
		got := ve.Rename(map[string]string{
			"Items.*.Name":     "items.*.name",
			"Items.*.Quantity": "items.*.quantity",
		})
		require.NotNil(t, got)
		require.Len(t, got.Fields, 2)
		assert.Equal(t, "items.0.name", got.Fields[0].Path)
		assert.Equal(t, "items.2.quantity", got.Fields[1].Path)
	})

	t.Run("no match → fallback to original path", func(t *testing.T) {
		t.Parallel()
		ve := &valid.Error{Fields: []valid.FieldError{
			{Path: "Internal.Field", Code: "VALIDATION_REQUIRED"},
		}}
		got := ve.Rename(map[string]string{"Other": "other"})
		require.NotNil(t, got)
		require.Len(t, got.Fields, 1)
		assert.Equal(t, "Internal.Field", got.Fields[0].Path)
	})

	t.Run("partial match — unmatched fields keep original path", func(t *testing.T) {
		t.Parallel()
		ve := &valid.Error{Fields: []valid.FieldError{
			{Path: "Reference", Code: "VALIDATION_REQUIRED"},
			{Path: "Internal.Only", Code: "VALIDATION_REQUIRED"},
		}}
		got := ve.Rename(map[string]string{"Reference": "reference"})
		require.NotNil(t, got)
		require.Len(t, got.Fields, 2)
		assert.Equal(t, "reference", got.Fields[0].Path)
		assert.Equal(t, "Internal.Only", got.Fields[1].Path)
	})
}

// ---- Integration: validateCreateOrder ---------------------------------------

type (
	CreateOrderParams struct {
		Reference      string
		BillingAddress BillingAddressParams
		Tags           []string
		Items          []OrderItemParams
		Payment        *PaymentParams
		Discounts      []*DiscountParams
	}

	DiscountParams struct {
		Type   string // e.g. "percentage" or "fixed"
		Amount uint32
	}

	PaymentParams struct {
		Method        string // "card", "bank_transfer", etc.
		TransactionID string
	}

	BillingAddressParams struct {
		Street string
		City   string
	}

	OrderItemParams struct {
		Name     string
		Quantity int
	}
)

func (o PaymentParams) Valid(ctx context.Context) error {
	return valid.Struct(ctx,
		valid.Field("Method", o.Method, is.Required, is.OneOf("card", "bank_transfer")),
		valid.Field("TransactionID", o.TransactionID, is.Required, is.HasPrefix("txn_")),
	)
}

func (o DiscountParams) Valid(ctx context.Context) error {
	return valid.Struct(ctx,
		valid.Field("Type", o.Type, is.Required, is.OneOf("percentage", "fixed")),
		valid.Field("Amount", o.Amount, is.Required, is.GreaterThan(0)),
	)
}

func validateCreateOrder(ctx context.Context, params CreateOrderParams) error {
	return valid.Struct(ctx,
		valid.Field("Payment", params.Payment, is.Required),
		valid.Nested("Payment", params.Payment),
		valid.Nested("Discounts", params.Discounts),
		valid.Field("Reference", params.Reference, is.Required),
		valid.Field("BillingAddress.Street", params.BillingAddress.Street, is.Required),
		valid.Field("BillingAddress.City", params.BillingAddress.City, is.Required),
		valid.Field("Tags", params.Tags, is.Required),
		valid.Slice("Items", params.Items, func(ctx context.Context, i int, item OrderItemParams) error {
			return valid.Struct(ctx,
				valid.Field("Name", item.Name, is.Required),
				valid.Field("Quantity", item.Quantity, is.Min(1)),
			)
		}),
	)
}

func TestIntegration(t *testing.T) {
	t.Parallel()

	mapping := map[string]string{
		"Payment":               "payment",
		"Payment.TransactionID": "payment.transactionId",
		"Payment.Method":        "payment.method",
		"Reference":             "reference",
		"BillingAddress.Street": "billingAddress.street",
		"BillingAddress.City":   "billingAddress.city",
		"Discounts.*.Type":      "discounts.*.type",
		"Discounts.*.Amount":    "discounts.*.amount",
		"Tags":                  "tags",
		"Items.*.Name":          "items.*.name",
		"Items.*.Quantity":      "items.*.quantity",
	}

	t.Run("valid params → nil", func(t *testing.T) {
		t.Parallel()
		params := CreateOrderParams{
			Reference:      "REF-001",
			BillingAddress: BillingAddressParams{Street: "1 Main St", City: "Springfield"},
			Tags:           []string{"urgent"},
			Items:          []OrderItemParams{{Name: "Widget", Quantity: 2}},
			Payment:        &PaymentParams{Method: "card", TransactionID: "txn_123"},
			Discounts:      []*DiscountParams{{Type: "percentage", Amount: 1}},
		}
		err := validateCreateOrder(context.Background(), params)
		require.NoError(t, err)
	})

	t.Run("empty params → multiple errors transformed", func(t *testing.T) {
		t.Parallel()
		params := CreateOrderParams{
			Items: []OrderItemParams{{Name: "", Quantity: 0}},
		}
		err := validateCreateOrder(context.Background(), params)
		require.Error(t, err)
		ve := valid.As(err)
		require.NotNil(t, ve)

		transformed := ve.Rename(mapping)
		require.NotNil(t, transformed)

		paths := make(map[string]bool)
		for _, fe := range transformed.Fields {
			paths[fe.Path] = true
		}

		godump.Dump(transformed)

		expected := []string{
			"payment",
			"reference",
			"billingAddress.street",
			"billingAddress.city",
			"tags",
			"items.0.name",
			"items.0.quantity",
		}
		require.Len(t, paths, 7)
		for _, p := range expected {
			assert.True(t, paths[p], "expected path %q in transformed errors, got: %v", p, transformed.Fields)
		}
	})

	t.Run("custom rule receives ctx value", func(t *testing.T) {
		t.Parallel()

		type ctxKey struct{}
		allowedCurrency := "EUR"

		var mustMatchCtxCurrency is.Rule = func(ctx context.Context, value any) *is.Violation {
			allowed, _ := ctx.Value(ctxKey{}).(string)
			s, ok := value.(string)
			if !ok || s != allowed {
				return &is.Violation{Code: "CURRENCY_MISMATCH", Message: "currency must match " + allowed}
			}
			return nil
		}

		ctx := context.WithValue(context.Background(), ctxKey{}, allowedCurrency)

		err := valid.Struct(ctx,
			valid.Field("Currency", "USD", mustMatchCtxCurrency),
		)
		require.Error(t, err)
		ve := valid.As(err)
		require.NotNil(t, ve)
		require.Len(t, ve.Fields, 1)
		assert.Equal(t, "Currency", ve.Fields[0].Path)
		assert.Equal(t, "CURRENCY_MISMATCH", ve.Fields[0].Code)

		err = valid.Struct(ctx,
			valid.Field("Currency", "EUR", mustMatchCtxCurrency),
		)
		require.NoError(t, err)
	})

	t.Run("partial errors — only failing fields appear", func(t *testing.T) {
		t.Parallel()
		params := CreateOrderParams{
			Reference:      "REF-002",
			BillingAddress: BillingAddressParams{Street: "2 Elm St", City: ""},
			Tags:           []string{"ok"},
			Items:          []OrderItemParams{{Name: "Gadget", Quantity: 1}},
			Payment:        &PaymentParams{Method: "card", TransactionID: "txn_456"},
		}
		err := validateCreateOrder(context.Background(), params)
		require.Error(t, err)
		ve := valid.As(err)
		transformed := ve.Rename(mapping)
		require.NotNil(t, transformed)
		require.Len(t, transformed.Fields, 1)
		assert.Equal(t, "billingAddress.city", transformed.Fields[0].Path)
	})

	t.Run("nested slice — non-adjacent failures all reported", func(t *testing.T) {
		t.Parallel()
		items := []DiscountParams{
			{Type: "percentage", Amount: 1},
			{Type: "percentage", Amount: 0},
			{Type: "percentage", Amount: 2},
			{Type: "percentage", Amount: 0},
		}
		err := valid.Struct(context.Background(),
			valid.Nested("Discounts", items),
		)
		require.Error(t, err)
		ve := valid.As(err)
		require.NotNil(t, ve)
		require.Len(t, ve.Fields, 2)
		assert.Equal(t, "Discounts.1.Amount", ve.Fields[0].Path)
		assert.Equal(t, string(is.ViolationRequired), ve.Fields[0].Code)
		assert.Equal(t, "Discounts.3.Amount", ve.Fields[1].Path)
		assert.Equal(t, string(is.ViolationRequired), ve.Fields[1].Code)
	})

	t.Run("each — per-element rule violations reported with indexed paths", func(t *testing.T) {
		t.Parallel()
		items := []string{"card", "bank_transfer", "crypto"}
		err := valid.Struct(context.Background(),
			valid.Field("PaymentTypes", items, is.Required),
			valid.Each("PaymentTypes", items, is.OneOf("card", "bank_transfer")),
		)
		require.Error(t, err)
		ve := valid.As(err)
		require.NotNil(t, ve)
		require.Len(t, ve.Fields, 1)
		assert.Equal(t, "PaymentTypes.2", ve.Fields[0].Path)
		assert.Equal(t, string(is.ViolationOneOf), ve.Fields[0].Code)
	})
}
