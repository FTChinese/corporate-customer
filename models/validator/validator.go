package validator

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"log"
	"strings"
	"unicode/utf8"
)

// Required checks if a string is empty.
// Leading and trailing spaces are trimmed.
func Required(str string) bool {
	return strings.TrimSpace(str) != ""
}

// StringInLength checks if a string's length, including multi bytes string,
// is within a range, inclusive.
func StringInLength(str string, min, max int) bool {
	if min > max {
		min, max = max, min
	}
	strLength := utf8.RuneCountInString(str)
	return strLength >= min && strLength <= max
}

// MinStringLength checks if a string's length is longer than min
func MinStringLength(str string, min int) bool {
	strLength := utf8.RuneCountInString(str)
	return strLength >= min
}

// MaxStringLength checks if a string's length is under max
// Return true if the length of str is under or equal to max; false otherwise
func MaxStringLength(str string, max int) bool {
	strLength := utf8.RuneCountInString(str)
	return strLength <= max
}

type Validator struct {
	fieldName  string
	isRequired bool
	min        int
	max        int
	isEmail    bool
	isURL      bool
}

func New(name string) *Validator {
	return &Validator{
		fieldName: name,
	}
}

func (v *Validator) Required() *Validator {
	v.isRequired = true
	return v
}

func (v *Validator) Min(min int) *Validator {
	v.min = min
	return v
}

func (v *Validator) Max(max int) *Validator {
	v.max = max
	return v
}

func (v *Validator) Range(min, max int) *Validator {
	v.min = min
	v.max = max
	return v
}

func (v *Validator) Email() *Validator {
	v.isEmail = true
	return v
}

func (v *Validator) URL() *Validator {
	v.isURL = true
	return v
}

func (v *Validator) Validate(value string) string {
	if v.isEmail && v.isURL {
		log.Fatal("The validated value cannot be both an email and url")
	}

	if v.isRequired && !Required(value) {
		return v.fieldName + "不能为空"
	}

	if v.min > 0 && v.max > 0 && !StringInLength(value, v.min, v.max) {
		return fmt.Sprintf("%s只允许%d～%d个字符。已输入%d个字符。", v.fieldName, v.min, v.max, len(value))
	}

	if v.min > 0 && !MinStringLength(value, v.min) {
		return fmt.Sprintf("%s不得少于%d个字符", v.fieldName, v.min)
	}

	if v.max > 0 && !MaxStringLength(value, v.max) {
		return fmt.Sprintf("%s不得超过%d个字符", v.fieldName, v.max)
	}

	if v.isEmail && !govalidator.IsEmail(value) {
		return fmt.Sprintf("%s不是有效的邮箱地址", value)
	}

	if v.isURL && !govalidator.IsURL(value) {
		return "请输入有效的URL"
	}

	return ""
}
