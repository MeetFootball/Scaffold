package util

import "errors"

var (
	ErrSystem          = errors.New("SystemError")
	ErrToken           = errors.New("TokenError")
	ErrScope           = errors.New("ScopeError")
	ErrMethodNotFound  = errors.New("MethodNotFound")
	ErrServiceNotFound = errors.New("ServiceNotFound")

	ErrLanguage                 = errors.New("LanguageError")            // 语言包错误
	ErrImageFormat              = errors.New("ImageFormat")              // 图片格式错误
	ErrFileOrDirectoryNotExists = errors.New("FileOrDirectoryNotExists") // 文件或路径不存在
)
