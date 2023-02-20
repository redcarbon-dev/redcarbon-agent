package utils

func GetUniqueDataForColumnInMap(data []map[string]string, column string) []string {
	var u []string

	for _, d := range data {
		value := d[column]
		if !exist(u, value) {
			u = append(u, value)
		}
	}

	return u
}

func exist(values []string, item string) bool {
	for _, v := range values {
		if item == v {
			return true
		}
	}

	return false
}

func GroupMapByColumn(data []map[string]string, column string) map[string][]map[string]string {
	groups := map[string][]map[string]string{}

	for _, d := range data {
		columnValue, _ := d[column]

		values, ok := groups[columnValue]

		if ok {
			groups[columnValue] = append(values, d)
		} else {
			groups[columnValue] = []map[string]string{d}
		}
	}

	return groups
}
