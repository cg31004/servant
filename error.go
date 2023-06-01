package cronjob

import (
	"fmt"
)

type ErrorCode string

const (
	ErrJobAlreadyRemove ErrorCode = "Job is already remove"
	ErrRegistered       ErrorCode = "Registered"
	ErrScheduleIsNil    ErrorCode = "Schedule is nil"
)

func newCronError(errorCode ErrorCode, msgs ...ErrorMsg) *CronError {
	e := &CronError{
		errorCode: errorCode,
	}

	for _, msg := range msgs {
		msg(e)
	}

	return e
}

type CronError struct {
	errorCode ErrorCode
	msg       string
}

func (err *CronError) Error() string {
	if err.msg == "" {
		return string(err.errorCode)
	}
	return fmt.Sprintf("%s: %s", err.errorCode, err.msg)
}

func (err *CronError) ErrorCode() ErrorCode {
	return err.errorCode
}

func (err *CronError) SetMsg(msg string) *CronError {
	err.msg = msg
	return err
}

type ErrorMsg func(o *CronError)

func ErrorWithMessage(msg string) ErrorMsg {
	return func(o *CronError) {
		if o.msg == "" {
			o.msg = msg
		} else {
			o.msg = fmt.Sprintf("%s %s", o.msg, msg)
		}
	}
}

func ErrorWithErrorMessage(err error) ErrorMsg {
	return ErrorWithMessage(err.Error())
}

//IsCronError 判斷是不是CronError
func IsCronError(err error) bool {
	_, ok := err.(*CronError)
	return ok
}

//ConvertCronError 判斷是不是CronError
func ConvertCronError(err error) (*CronError, bool) {
	if IsCronError(err) {
		return err.(*CronError), true
	}

	return nil, false
}

//CompareErrorCode 比較err 是不是屬於欲比對的ErrorCode
func CompareErrorCode(err error, errorCode ErrorCode) bool {
	cronError, ok := ConvertCronError(err)
	if !ok {
		return false
	}

	return cronError.ErrorCode() == errorCode
}
