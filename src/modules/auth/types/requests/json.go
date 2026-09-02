package requests

import (
	"bytes"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/parsers/zjson"
)

type Schema interface {
	Parse(data any, destPtr any, options ...z.ExecOption) z.ZogIssueList
}

func ParseJSON(body []byte, schema Schema, dest any) z.ZogIssueList {
	return schema.Parse(zjson.Decode(bytes.NewReader(body)), dest)
}

func ParseMap(data map[string]any, schema Schema, dest any) z.ZogIssueList {
	return schema.Parse(data, dest)
}

func Flatten(issues z.ZogIssueList) map[string][]string {
	return z.Issues.Flatten(issues)
}
