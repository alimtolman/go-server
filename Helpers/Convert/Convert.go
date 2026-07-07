package Convert

import "strconv"

func StrToUInt(value string, def uint32) uint32 {
	if number, err := strconv.ParseUint(value, 10, 32); err == nil {
		return uint32(number)
	}

	return def
}
