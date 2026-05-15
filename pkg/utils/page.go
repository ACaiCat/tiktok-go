package utils

func NormalizePage(pageSize int32, pageNum int32, defaultPageSize int32, maxPageSize int32) (int, int) {
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	if pageNum < 0 {
		pageNum = 0
	}

	return int(pageSize), int(pageNum)
}
