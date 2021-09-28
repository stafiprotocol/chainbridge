package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func VersionCompare(version1, version2 string) (ret int) {
	defer func() {
		switch {
		case ret < 0:
			ret = -1
		case ret > 0:
			ret = 1
		}
	}()

	version1Slice := strings.Split(version1, ".")
	version2Slice := strings.Split(version2, ".")
	if len(version1Slice) != len(version2Slice) || len(version1Slice) != 3 {
		panic(fmt.Sprintf("version format err: 1 %s 2 %s", version1, version2))
	}

	for i := range version1Slice {
		v1, err := strconv.Atoi(version1Slice[i])
		if err != nil {
			panic(err)
		}
		v2, err := strconv.Atoi(version2Slice[i])
		if err != nil {
			panic(err)
		}
		if v1 == v2 {
			continue
		} else {
			return v1 - v2
		}

	}
	return 0

}
