// Code generated by protoc-gen-errors. DO NOT EDIT.
// versions:
// protoc-gen-errors v1.2.0

package gerror

import (
	"github.com/zhanmengao/gf/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

//系统错误
func ErrNone() *errors.Error {
	return errors.NewErrorSkip(0, "Err None", "Err None", 1)
}

//SERVER DECODE
func ErrServerDecode() *errors.Error {
	return errors.NewErrorSkip(10, "Err Server Decode", "Err Server Decode", 1)
}

//SERVER ENCODE
func ErrServerEncode() *errors.Error {
	return errors.NewErrorSkip(20, "Err Server Encode", "Err Server Encode", 1)
}

//SERVICE NOT_EXISTS
func ErrServerServiceNotExists() *errors.Error {
	return errors.NewErrorSkip(30, "Err Server Service Not Exists", "Err Server Service Not Exists", 1)
}

//SERVER NETWORK
func ErrServerNetwork() *errors.Error {
	return errors.NewErrorSkip(40, "Err Server Network", "Err Server Network", 1)
}

//SERVER DATABASE
func ErrServerDatabase() *errors.Error {
	return errors.NewErrorSkip(41, "Err Server Database", "Err Server Database", 1)
}

//SERVER_AUTH
func ErrServerAuth() *errors.Error {
	return errors.NewErrorSkip(50, "Err Server Auth", "Err Server Auth", 1)
}

//SERVER_BAD_PARAM
func ErrServerBadParam() *errors.Error {
	return errors.NewErrorSkip(60, "Err Server Bad Param", "Err Server Bad Param", 1)
}

//SERVER_BUSY
func ErrServerBusy() *errors.Error {
	return errors.NewErrorSkip(61, "Err Server Busy", "Err Server Busy", 1)
}

//CMD_NOT_REGISTER
func ErrServerCmdNotRegister() *errors.Error {
	return errors.NewErrorSkip(70, "Err Server Cmd Not Register", "Err Server Cmd Not Register", 1)
}

//INTERNAL_UNKNOWN
func ErrServerInternalUnknown() *errors.Error {
	return errors.NewErrorSkip(99, "Err Server Internal Unknown", "Err Server Internal Unknown", 1)
}

//DATA_NOT_FOUND
func ErrDataNotFound() *errors.Error {
	return errors.NewErrorSkip(404, "Err Data Not Found", "Err Data Not Found", 1)
}

//CONFIG_NOT_FOUND
func ErrConfigNotFound() *errors.Error {
	return errors.NewErrorSkip(101, "Err Config Not Found", "Err Config Not Found", 1)
}

//CONFIG_INCORRECT
func ErrConfigIncorrect() *errors.Error {
	return errors.NewErrorSkip(102, "Err Config Incorrect", "Err Config Incorrect", 1)
}

//ErrFreqLimit 频率限制错误
func ErrFreqLimit() *errors.Error {
	return errors.NewErrorSkip(1001, "Err Freq Limit", "Err Freq Limit", 1)
}

//ErrConcurrency 并发错误
func ErrConcurrency() *errors.Error {
	return errors.NewErrorSkip(1002, "Err Concurrency", "Err Concurrency", 1)
}

//ErrSrvStopped 停服
func ErrSrvStopped() *errors.Error {
	return errors.NewErrorSkip(1003, "Err Srv Stopped", "Err Srv Stopped", 1)
}

//ErrUserBanned 被封号
func ErrUserBanned() *errors.Error {
	return errors.NewErrorSkip(1004, "Err User Banned", "Err User Banned", 1)
}
