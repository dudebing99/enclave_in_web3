package utils

// CheckError 检查错误，如果有错误会 panic
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
