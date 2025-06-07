# Prompt

1. Generate a query to display array string for news detail

```
Based on this table

Table topics {
  id SERIAL [pk, increment]
  name VARCHAR(100) [unique, not null]
  description TEXT
  slug VARCHAR(100) [unique, not null]
  created_at TIMESTAMP [default: `CURRENT_TIMESTAMP`]
  updated_at TIMESTAMP [default: `CURRENT_TIMESTAMP`]
  deleted_at TIMESTAMP
}

Table news_articles {
  id SERIAL [pk, increment]
  title VARCHAR(255) [not null]
  content TEXT [not null]
  summary TEXT
  author_id INTEGER [not null]
  slug VARCHAR(255) [unique, not null]
  status article_status [not null, default: 'draft']
  published_at TIMESTAMP
  created_at TIMESTAMP [default: `CURRENT_TIMESTAMP`]
  updated_at TIMESTAMP [default: `CURRENT_TIMESTAMP`]
  deleted_at TIMESTAMP
}

Table news_topics {
  id SERIAL [pk, increment]
  news_article_id INTEGER [not null]
  topic_id INTEGER [not null]
  created_at TIMESTAMP [default: `CURRENT_TIMESTAMP`]
  deleted_at TIMESTAMP
}

Generate a query to display detail news by slug, the query must show
- id
- title
- content
- slug
- author_name
- published_at
- topics: is a name of topic that retrieve from topic table, this field types is an array of string
```

2. Make update topic support update field partially

```
Based on this function

func (r topicRepository) UpdateTopicFileds(ctx context.Context, topic *entity.Topic) error {
	query := "UPDATE topics SET name=$1, description=$2, slug=$3 WHERE id=$4"

	_, err := r.db.ExecContext(ctx, query, topic.name, topic.description, topic.slug, topic.id)
	return err
}
```

3. I use the GPT to generate several unit test function. Here is the prompt

```
You're given a unit test setup function for usecase layer like this

type TopicsAccessor struct {
	topicRepo *mock_repository.MockTopicsRepository
	topicUC   usecase.TopicsUsecase
}

func newTopicAccessor(ctrl *gomock.Controller) TopicsAccessor {
	repo := mock_repository.NewMockTopicsRepository(ctrl)
	topicUC := usecase.NewTopicsUsecase(repo)
	return TopicsAccessor{
		topicRepo: repo,
		topicUC:   topicUC,
	}
}

Generate a unit test for function [Mention the function]. Follow this example unit test structure


func Test_CreateUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newUserAccessor(ctrl)
	repo := accessor.userRepo
	uc := accessor.userUC
	ctx := context.Background()

	mockReq := request.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	tests := []struct {
		testname  string
		initMock  func()
		assertion func(err error)
	}{
		{
			testname: "create user and repository return err then uc return error",
			initMock: func() {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			err := uc.CreateUser(ctx, mockReq)
			tt.assertion(err)
		})
	}
}

```

4. Asking GPT to generate github workflow

```
Generate a github action workflow config to run the unit test for the golang project. The project using go version 1.24.2 and to run the unit test you can use Makefile command

make run_test

```

5. Asking GPT to generate `news-api.yaml`.

```
Generate a api docs swagger with name `news-api.yaml` with requirement

1. Get All Topics
Endpoint: GET /api/v1/topics
Response JSON:
{
    "data": [
        {
            "id": 1,
            "name": "Technology",
            "description": "Latest tech news and innovations",
            "slug": "technology",
            "updated_at": "2025-06-05T06:33:08.982231Z"
        },
    ],
    "http_status": 200
}

(I list down all the API one by one including body request and response example)
```

6. Asking GPT to create grafana dashboard

```
Create grafana dashboard JSON to visualization request per second, average request duration for each endpoint, percentile latency, and error rate 5xx
```
