package ishelper

import (
	"math"
	"math/big"
)

type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type floating interface {
	~float32 | ~float64
}

// Number is the type constraint for numeric values accepted by Min and Max.
type Number interface {
	signed | unsigned | floating
}

// ToRat converts any numeric value to *big.Rat for exact comparison.
// Returns (nil, false) if the value is not a supported numeric type.
func ToRat(value any) (*big.Rat, bool) {
	switch v := value.(type) {
	case int:
		return big.NewRat(int64(v), 1), true
	case int8:
		return big.NewRat(int64(v), 1), true
	case int16:
		return big.NewRat(int64(v), 1), true
	case int32:
		return big.NewRat(int64(v), 1), true
	case int64:
		return big.NewRat(v, 1), true
	case uint:
		return new(big.Rat).SetUint64(uint64(v)), true
	case uint8:
		return new(big.Rat).SetUint64(uint64(v)), true
	case uint16:
		return new(big.Rat).SetUint64(uint64(v)), true
	case uint32:
		return new(big.Rat).SetUint64(uint64(v)), true
	case uint64:
		return new(big.Rat).SetUint64(v), true
	case uintptr:
		return new(big.Rat).SetUint64(uint64(v)), true
	case float32:
		f := float64(v)
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return nil, false
		}
		r := new(big.Rat)
		r.SetFloat64(f)
		return r, true
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, false
		}
		r := new(big.Rat)
		r.SetFloat64(v)
		return r, true
	default:
		return nil, false
	}
}
