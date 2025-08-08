package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func Validate(input string) bool {
	// 首先检查是否是空字符串
	if input == "" {
		fmt.Println("错误：输入为空字符串")
		return false
	}

	// 尝试解析JSON
	var js any
	decoder := json.NewDecoder(bytes.NewReader([]byte(input)))
	decoder.DisallowUnknownFields() // 禁止未知字段，严格校验

	err := decoder.Decode(&js)
	if err != nil {
		// 获取更详细的错误信息
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			// 计算错误位置附近的上下文
			start := max(0, syntaxErr.Offset-10)
			end := min(len(input), int(syntaxErr.Offset+10))
			context := input[start:end]

			fmt.Printf("JSON语法错误在位置 %d:\n", syntaxErr.Offset)
			fmt.Printf("错误附近内容: ...%s...\n", context)
			fmt.Printf("详细错误: %v\n", err)
		} else if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Printf("类型错误在字段 '%s': 期望类型 %s, 实际值: %v\n",
				unmarshalErr.Field, unmarshalErr.Type, unmarshalErr.Value)
		} else {
			fmt.Printf("JSON校验失败: %v\n", err)
		}
		return false
	}

	return true
}

func FormatJSON(input string) error {
	var formatted bytes.Buffer
	err := json.Indent(&formatted, []byte(input), "", "  ")
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println(formatted.String())
	return nil
}

func CompressJSON(jsonStr string) error {
	var data any
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return err
	}

	compressed, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println(string(compressed))

	return nil
}

var charMap = []string{`\"`, `\\`, `\/`}

func Unescape(s string) (string, bool) {

	needReplace := false
	for _, char := range charMap {
		if strings.Contains(s, char) {
			needReplace = true
		}
	}

	if needReplace {
		s = strings.ReplaceAll(s, `\"`, `"`) // 转义引号
		s = strings.ReplaceAll(s, `\\`, `\`) // 转义反斜杠
		s = strings.ReplaceAll(s, `\/`, `/`) // 转义斜杠
	}

	return s, needReplace
}
