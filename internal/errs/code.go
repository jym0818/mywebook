package errs

//用户模块

const (
	//客户端错误  遇到这个错误不用管 prometheus不需要监控
	UserInputValid = 401001
	//这个错误需要监控，频繁出现可能有问题，可能有人在暴力破解
	UserInvalidOrPassword   = 401002
	UserInternalServerError = 501001
)

//文章模块

const (
	ArticleInvalidInput        = 402001
	ArticleInternalServerError = 502001
)
