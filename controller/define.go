package controller

import "aweme_kitex/models"

type (
	Video    models.Video
	User     models.User
	Response models.Response
	Comment  models.Comment
)

var (
	db = models.DB
)
