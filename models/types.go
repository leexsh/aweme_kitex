package models

// ----------types & response---------------
type Favourite struct {
	Id      string `json:"identity,omitempty"`
	UserId  string `json:"user_id,omitempty"`
	VideoId string `json:"video_id,omitempty"`
}

// video
type Video struct {
	Id             string `json:"id"`     // id
	Author         User   `json:"author"` // author
	PlayUrl        string `json:"play_url"`
	CoverUrl       string `json:"cover_url"`
	FavouriteCount int64  `json:"favourite_count"`
	CommentCount   int64  `json:"comment_count"`
	IsFavourite    bool   `json:"is_favourite"`
	Title          string `json:"title"`
}

type Comment struct {
	Id         string `json:"id,omitempty"`
	User       User   `json:"author"`
	VideoId    string `json:"video_id"`
	Content    string `json:"content"`
	CreateDate string `json:"createDate"`
}

type User struct {
	UserId        string `json:"identity"`
	Name          string `json:"name"`
	Password      string `json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// -------------response--------------
type Response struct {
	StatusCode int32  `json:"status_code,omitempty"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserLogRstResponse struct {
	Response
	UserId   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Token    string `json:"token,omitempty"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type UserListResponse struct {
	Response
	UserList []*User `json:"user_list"`
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

// comment
type CommentActionResponse struct {
	Response
	Comment `json:"comment,omitempty"`
}
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}
