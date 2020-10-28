package utils

import "strconv"

// StringToUintSlice converts a string slice to uint slice
func StringToUintSlice(stringSlice []string) ([]uint, error) {
	tagIDsSliceInts := make([]uint, 0)

	for _, s := range stringSlice {
		sInt, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}

		tagIDsSliceInts = append(tagIDsSliceInts, uint(sInt))
	}

	return tagIDsSliceInts, nil
}
