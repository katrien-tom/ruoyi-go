package post

type PostResponse struct {
	PostID   int64   `json:"postId"`
	PostCode string  `json:"postCode"`
	PostName string  `json:"postName"`
	PostSort int     `json:"postSort"`
	Status   string  `json:"status"`
	Remark   *string `json:"remark"`
}

type PostListResponse struct {
	Rows  []PostResponse `json:"rows"`
	Total int64          `json:"total"`
}
