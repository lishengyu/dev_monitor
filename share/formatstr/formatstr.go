package formatstr

import "strings"

func JoinMapWithSep(m map[string]string, sep string) string {
	if len(m) == 0 {
		return ""
	}
	var result strings.Builder
	for _, v := range m {
		result.WriteString(v)
		result.WriteString(sep)
	}
	return result.String()[:result.Len()-len(sep)]
}
