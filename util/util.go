package util

import (
	"Laurene/go-log"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"time"
)

func IntInSlice(slice []int64, aI interface{}) bool {
	var a int64

	switch aI.(type) {
	case int:
		a = int64(aI.(int))
	case int64:
		a = aI.(int64)
	default:
		return false
	}

	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}

func StringInSlice(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ErrCheckFatal(err error, strs ...string) {
	if err != nil {
		if len(strs) != 0 {
			for _, str := range strs {
				err = errors.Wrap(err, str)
			}
		}
		log.Fatal(err)
	}
}

func TextCut(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n]) + "..."
	}
	return s
}

func FloatCut(f float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.8f", f), "0"), ".")
}

func CreateKey(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteRune(chars[r.Intn(len(chars))])
	}
	return b.String()
}
