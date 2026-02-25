package is

import (
	"context"
)

type (
	Violation struct {
		Code    ViolationCode
		Message string
	}

	ViolationCode string

	Rule func(ctx context.Context, value any) *Violation
)

const (
	ViolationRequired  ViolationCode = "REQUIRED"
	ViolationMin       ViolationCode = "MIN"
	ViolationMax       ViolationCode = "MAX"
	ViolationLength    ViolationCode = "LENGTH"
	ViolationOneOf     ViolationCode = "ONE_OF"
	ViolationHasPrefix ViolationCode = "HAS_PREFIX"
	ViolationHasSuffix ViolationCode = "HAS_SUFFIX"
	ViolationContains  ViolationCode = "CONTAINS"
	ViolationEmail     ViolationCode = "EMAIL"
	ViolationURL       ViolationCode = "URL"
	ViolationUUID      ViolationCode = "UUID"
	ViolationNumeric   ViolationCode = "NUMERIC"
	ViolationAlpha     ViolationCode = "ALPHA"
	ViolationAlphaNum  ViolationCode = "ALPHANUMERIC"
	ViolationMatches   ViolationCode = "MATCHES"
	ViolationBetween   ViolationCode = "BETWEEN"
	ViolationPositive  ViolationCode = "POSITIVE"
	ViolationNonNeg    ViolationCode = "NON_NEGATIVE"
	ViolationGT        ViolationCode = "GT"
	ViolationGTE       ViolationCode = "GTE"
	ViolationLT        ViolationCode = "LT"
	ViolationLTE       ViolationCode = "LTE"
	ViolationEQ        ViolationCode = "EQ"
	ViolationMinLength ViolationCode = "MIN_LENGTH"
	ViolationMaxLength ViolationCode = "MAX_LENGTH"
	ViolationNotEmpty  ViolationCode = "NOT_EMPTY"
)

var Messages = map[ViolationCode]string{
	ViolationRequired:  "is required",
	ViolationMin:       "must be >= {min}",
	ViolationMax:       "must be <= {max}",
	ViolationLength:    "length must be between {min} and {max}",
	ViolationOneOf:     "must be one of {values}",
	ViolationHasPrefix: "must start with {prefix}",
	ViolationHasSuffix: "must end with {suffix}",
	ViolationContains:  "must contain {value}",
	ViolationEmail:     "must be a valid email",
	ViolationURL:       "must be a valid URL",
	ViolationUUID:      "must be a valid UUID",
	ViolationNumeric:   "must be numeric",
	ViolationAlpha:     "must contain only letters",
	ViolationAlphaNum:  "must contain only letters and digits",
	ViolationMatches:   "must match pattern {pattern}",
	ViolationBetween:   "must be between {min} and {max}",
	ViolationPositive:  "must be > 0",
	ViolationNonNeg:    "must be >= 0",
	ViolationGT:        "must be > {value}",
	ViolationGTE:       "must be >= {value}",
	ViolationLT:        "must be < {value}",
	ViolationLTE:       "must be <= {value}",
	ViolationEQ:        "must be = {value}",
	ViolationMinLength: "length must be >= {min}",
	ViolationMaxLength: "length must be <= {max}",
	ViolationNotEmpty:  "must not be empty",
}
