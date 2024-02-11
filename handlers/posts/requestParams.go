package posts

type createPostRequestParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type changePostRequestParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
