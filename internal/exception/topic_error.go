package exception

var (
	ErrFailedInsertTopic = CustomError{Code: 20001, Message: "failed insert topic"}
	ErrTopicNotFound     = CustomError{Code: 20002, Message: "topic not found"}
	ErrFailedGetTopic    = CustomError{Code: 20003, Message: "failed get topic"}
	ErrFailedUpdateTopic = CustomError{Code: 20004, Message: "failed update topic"}
	ErrNoFieldUpdate     = CustomError{Code: 20005, Message: "no field update"}
)
