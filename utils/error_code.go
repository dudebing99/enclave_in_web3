package utils

type ErrorCode uint32

const (
	Success ErrorCode = iota
	InvalidParameter
	InternalError
	UpstreamError
)
