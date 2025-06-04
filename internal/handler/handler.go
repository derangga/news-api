package handler

type HandlerRegistry struct {
	TopicsHandler       TopicsHandler
	NewsArticlesHandler NewsHandler
}

func NewHandlerRegistry(
	topicsHandler TopicsHandler,
	newsArticlesHandler NewsHandler,
) HandlerRegistry {
	return HandlerRegistry{
		TopicsHandler:       topicsHandler,
		NewsArticlesHandler: newsArticlesHandler,
	}
}
