package is

import (
	"fmt"
	"strings"
)

func formatMessage(code ViolationCode, params map[string]any) string {
	template, ok := Messages[code]
	if !ok {
		return "invalid value"
	}

	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		template = strings.ReplaceAll(template, placeholder, fmt.Sprintf("%v", value))
	}

	return template
}
