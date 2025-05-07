package server

var requestsMap = make(map[string]int)

func GetByIp(ip string) int {
	if count, ok := requestsMap[ip]; ok {
		return count
	}

	return 0
}

func IncrementByIp(ip string) {
	if count, ok := requestsMap[ip]; ok {
		requestsMap[ip] = count + 1
	} else {
		requestsMap[ip] = 1
	}
}
