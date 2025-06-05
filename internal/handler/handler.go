package handler

type HandlerRegistry struct {
	UsersHandler        UsersHandler
	TopicsHandler       TopicsHandler
	NewsArticlesHandler NewsHandler
}

func NewHandlerRegistry(
	usersHandler UsersHandler,
	topicsHandler TopicsHandler,
	newsArticlesHandler NewsHandler,
) HandlerRegistry {
	return HandlerRegistry{
		UsersHandler:        usersHandler,
		TopicsHandler:       topicsHandler,
		NewsArticlesHandler: newsArticlesHandler,
	}
}
