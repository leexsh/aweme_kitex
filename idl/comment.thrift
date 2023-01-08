include "user.thrift"
include "feed.thrift"
namespace go comment

struct Comment {
    1: required string comment_id // 评论id
    2: user.User user
    3: string content // 评论内容
    4: string create_time // 创建日期
}

struct CommentActionRequest {
    1: required string token // token
    2: required string video_id // 评论视频id
    3: string action_type // 1-发布评论 2-删除评论
    4: optional string comment_content // 评论内容
    5: string comment_id // 评论id
}

struct CommentActionResponse{
    1: feed.BaseResp base_resp
    3: list<Comment> comment_list
}

struct CommentListRequest {
    1: required string token
    2: string video_id
}

struct CommentListResponse {
    1: feed.BaseResp base_resp
    2: list<Comment> comment_list
}
