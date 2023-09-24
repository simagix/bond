/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * utils.go
 */

package bond

import (
	"fmt"
	"regexp"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func Stringify(doc interface{}, opts ...string) string {
	var b []byte
	var err error
	if len(opts) == 2 {
		b, err = bson.MarshalExtJSONIndent(doc, false, false, opts[0], opts[1])
		if err != nil {
			return err.Error()
		}
	} else {
		b, err = bson.MarshalExtJSON(doc, false, false)
		if err != nil {
			return err.Error()
		}
	}
	str := string(b)
	re := regexp.MustCompile(`{\s*"\$oid":\s?("[a-fA-F0-9]{24}")\s*}`)
	str = re.ReplaceAllString(str, "ObjectId($1)")
	re = regexp.MustCompile(`{\s*"\$date":\s?("\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z")\s*}`)
	str = re.ReplaceAllString(str, "ISODate($1)")
	return str
}

func StringifyIndent(doc interface{}) string {
	return Stringify(doc, "", "  ")
}

// ToInt converts to int
func ToInt(num interface{}) int {
	f := fmt.Sprintf("%v", num)
	x, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return 0
	}
	return int(x)
}

// ToInt64 converts to int64
func ToInt64(num interface{}) int64 {
	f := fmt.Sprintf("%v", num)
	x, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return 0
	}
	return int64(x)
}
