package kv

import (
	"strconv"
	"strings"
)

// 智能分割键值对（处理引号包裹的空格）
func splitKV(input string) []string {
	var (
		buffer    strings.Builder
		pairs     []string
		inQuotes  bool
		quoteChar byte
	)

	for i := 0; i < len(input); i++ {
		c := input[i]

		switch {
		case c == '\\' && i+1 < len(input): // 处理转义字符
			buffer.WriteByte(input[i+1])
			i++
		case (c == '"' || c == '\'') && !inQuotes:
			inQuotes = true
			quoteChar = c
		case c == quoteChar && inQuotes:
			inQuotes = false
			quoteChar = 0
		case c == ' ' && !inQuotes:
			if buffer.Len() > 0 {
				pairs = append(pairs, buffer.String())
				buffer.Reset()
			}
		default:
			buffer.WriteByte(c)
		}
	}

	if buffer.Len() > 0 {
		pairs = append(pairs, buffer.String())
	}
	return pairs
}

func ParseKV(input string) (map[string]string, error) {
	result := make(map[string]string)
	pairs := splitKV(input) // 步骤1：分割键值对

	for _, pair := range pairs {
		// 步骤2：分割键值
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue // 跳过无效格式
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 步骤3：处理带引号的值
		if len(value) > 0 && (value[0] == '"' || value[0] == '\'') {
			unquoted, err := strconv.Unquote(value)
			if err == nil {
				value = unquoted
			}
		}

		result[key] = value
	}
	return result, nil
}
