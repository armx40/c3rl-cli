package main

const (
	DATA_LOG_IS_ENCRYPTED                 = 1 << 0
	DATA_LOG_TIME_ENCRYPTED               = 1 << 1
	DATA_LOG_CODE_ENCRYPTED               = 1 << 2
	DATA_LOG_TAG_ENCRYPTED                = 1 << 3
	DATA_LOG_DATA_ENCRYPTED               = 1 << 4
	DATA_LOG_SIGNED                       = 1 << 5
	DATA_LOG_TIME_SIGNED                  = 1 << 6
	DATA_LOG_CODE_SIGNED                  = 1 << 7
	DATA_LOG_TAG_SIGNED                   = 1 << 8
	DATA_LOG_DATA_SIGNED                  = 1 << 9
	DATA_LOG_HMACED                       = 1 << 10
	DATA_LOG_TIME_HMACED                  = 1 << 11
	DATA_LOG_CODE_HMACED                  = 1 << 12
	DATA_LOG_TAG_HMACED                   = 1 << 13
	DATA_LOG_DATA_HMACED                  = 1 << 14
	DATA_LOG_CHECKSUM_INCLUDED            = 1 << 15
	DATA_LOG_SIZE_INCLUDES_SIZE_BYTES     = 1 << 16
	DATA_LOG_SIZE_INCLUDES_CHECKSUM_BYTES = 1 << 17
	DATA_LOG_TIME_IS_32_BITS              = 1 << 18
)

var API_HOST = ""

// const API_HOST = "http://127.0.0.1:9000/api/c3rl-cli/"
// const API_HOST = "https://es.c3rl.com/api/c3rl-cli/"
