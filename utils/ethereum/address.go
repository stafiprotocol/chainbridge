package ethereum

import "regexp"

var re = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

func IsAddressValid(addr string) bool {
	return re.MatchString(addr)
}
