package xRedis

func FundName(code string) string {
	return "fund:" + code
}

func FundLock() string {
	return "fund:lock"
}
