package post

type PostListRequest struct {
	PostCode string `json:"postCode" form:"postCode"`
	PostName string `json:"postName" form:"postName"`
	Status   string `json:"status" form:"status"`
}

type AddPostRequest struct {
	PostCode string  `json:"postCode" binding:"required"`
	PostName string  `json:"postName" binding:"required"`
	PostSort int     `json:"postSort"`
	Status   string  `json:"status"`
	Remark   *string `json:"remark"`
}

type EditPostRequest struct {
	PostID   int64   `json:"postId" binding:"required"`
	PostCode string  `json:"postCode" binding:"required"`
	PostName string  `json:"postName" binding:"required"`
	PostSort int     `json:"postSort"`
	Status   string  `json:"status"`
	Remark   *string `json:"remark"`
}
