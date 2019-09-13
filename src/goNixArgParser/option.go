package goNixArgParser

import (
	"bytes"
	"strings"
)

func (opt *Option) isDelimiter(r rune) bool {
	for _, delimiter := range opt.Delimiters {
		if r == delimiter {
			return true
		}
	}
	return false
}

func (opt *Option) splitValues(input string) []string {
	values := strings.FieldsFunc(input, opt.isDelimiter)
	values = opt.filterValues(values)

	return values
}

func (opt *Option) filterValues(values []string) []string {
	if opt.UniqueValues {
		uniqueValues := make([]string, 0, len(values))
		uniqueValues = appendUnique(uniqueValues, values...)
		values = uniqueValues
	}

	return values
}

func (opt *Option) GetHelp() []byte {
	buffer := &bytes.Buffer{}

	for i, flag := range opt.Flags {
		if i > 0 {
			buffer.WriteString("|")
		}
		buffer.WriteString(flag.Name)
	}

	if opt.AcceptValue {
		buffer.WriteString(" <value>")
		if opt.MultiValues {
			buffer.WriteString(" ...")
		}
	}

	buffer.WriteByte('\n')

	if len(opt.EnvVars) > 0 {
		buffer.WriteString("EnvVar: ")

		for i, envVar := range opt.EnvVars {
			if i > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(envVar)
		}

		buffer.WriteByte('\n')
	}

	if len(opt.DefaultValues) > 0 {
		buffer.WriteString("Default: ")

		for i, d := range opt.DefaultValues {
			if i > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(d)
		}

		buffer.WriteByte('\n')
	}

	if len(opt.Summary) > 0 {
		buffer.WriteString(opt.Summary)
		buffer.WriteByte('\n')
	}

	if len(opt.Description) > 0 {
		buffer.WriteString(opt.Description)
		buffer.WriteByte('\n')
	}

	return buffer.Bytes()
}
