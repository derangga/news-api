package exception

var (
	ErrFailedInsertNews      = CustomError{Code: 20001, Message: "failed insert news"}
	ErrNewsNotFound          = CustomError{Code: 20002, Message: "news not found"}
	ErrFailedGetNews         = CustomError{Code: 20003, Message: "failed get news"}
	ErrFailedUpdateNews      = CustomError{Code: 20004, Message: "failed update news"}
	ErrFailedUpdateTopicNews = CustomError{Code: 20005, Message: "failed update topic news"}
	ErrFailedDeleteNews      = CustomError{Code: 20006, Message: "failed delete news"}
	ErrFailedDeleteTopicNews = CustomError{Code: 20007, Message: "failed delete topic news"}
)
