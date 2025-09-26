package dot

import (
	"regexp"
)

const (
	mobilePattern = `^1[34578][0-9]{9}$`
	mailPattern   = `^[a-z0-9A-Z]+([\-_\.][a-z0-9A-Z]+)*@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)*?\.)+[a-zA-Z]{2,4}$`
)

var (
	mobileRegexp = regexp.MustCompile(mobilePattern)
	mailRegexp   = regexp.MustCompile(mailPattern)
)

func IsMobile(s any) bool {
	if val, ok := s.([]byte); ok {
		return mobileRegexp.Match(val)
	}

	return mobileRegexp.MatchString(ToString(s))
}

func IsMail(s any) bool {
	if val, ok := s.([]byte); ok {
		return mailRegexp.Match(val)
	}

	return mailRegexp.MatchString(ToString(s))
}
