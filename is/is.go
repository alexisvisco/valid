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
	ViolationRequired  ViolationCode = "VALIDATION_REQUIRED"
	ViolationMin       ViolationCode = "VALIDATION_MIN"
	ViolationMax       ViolationCode = "VALIDATION_MAX"
	ViolationLength    ViolationCode = "VALIDATION_LENGTH"
	ViolationOneOf     ViolationCode = "VALIDATION_ONE_OF"
	ViolationHasPrefix ViolationCode = "VALIDATION_HAS_PREFIX"
	ViolationHasSuffix ViolationCode = "VALIDATION_HAS_SUFFIX"
	ViolationContains  ViolationCode = "VALIDATION_CONTAINS"
	ViolationEmail     ViolationCode = "VALIDATION_EMAIL"
	ViolationURL       ViolationCode = "VALIDATION_URL"
	ViolationUUID      ViolationCode = "VALIDATION_UUID"
	ViolationNumeric   ViolationCode = "VALIDATION_NUMERIC"
	ViolationAlpha     ViolationCode = "VALIDATION_ALPHA"
	ViolationAlphaNum  ViolationCode = "VALIDATION_ALPHANUMERIC"
	ViolationMatches   ViolationCode = "VALIDATION_MATCHES"
	ViolationBetween   ViolationCode = "VALIDATION_BETWEEN"
	ViolationPositive  ViolationCode = "VALIDATION_POSITIVE"
	ViolationNonNeg    ViolationCode = "VALIDATION_NON_NEGATIVE"
	ViolationGT        ViolationCode = "VALIDATION_GT"
	ViolationGTE       ViolationCode = "VALIDATION_GTE"
	ViolationLT        ViolationCode = "VALIDATION_LT"
	ViolationLTE       ViolationCode = "VALIDATION_LTE"
	ViolationEQ        ViolationCode = "VALIDATION_EQ"
	ViolationMinLength ViolationCode = "VALIDATION_MIN_LENGTH"
	ViolationMaxLength ViolationCode = "VALIDATION_MAX_LENGTH"
	ViolationNotEmpty  ViolationCode = "VALIDATION_NOT_EMPTY"
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
