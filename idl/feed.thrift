include "user.thrift"
include "base.thrift"
namespace go feed

struct BaseResp {
    1: i64 status_code
    2: string status_msg
    3: i64 service_time
}

struct Video {
    1: string video_id // video id
    2: user.User author // 用户信息
    3: string play_url // 播放地址
    4: string cover_url // 封面地址
    5: i64 favourite_count // 点赞数量
    6: i64 comment_count // 评论数量
    7: bool is_favourite // 是否已经点赞
    8: string title // 标题
}

struct FeedRequest {
    1: i64 latest_time
    2: string token
}

struct FeedResponse {
    1: base.BaseResp base_resp
    2: list<Video> video_list
    3: i64 next_time
}

struct ChangeCommentCountRequest {
    1: string video_id
    2: i64 action
}

struct ChangeCommentCountResponse {
    1: base.BaseResp base_resp
}

struct CheckVideoInvalidRequest {
    1: list<string> video_id
}

struct CheckVideoInvalidResponse {
    1: base.BaseResp base_resp
}

struct GetVideosResponse {
    1: base.BaseResp base_resp
    2: list<Video> videos
}

service FeedService {
    // feed流操作
    FeedResponse Feed(1: FeedRequest req)
    // 修改评论数目
    ChangeCommentCountResponse ChangeCommentCnt(1: ChangeCommentCountRequest req)
    // 查询vid 是否合理
    CheckVideoInvalidResponse CheckVideoInvalid(1: CheckVideoInvalidRequest req)
    // 获取视频
    GetVideosResponse GetVideosById(1: CheckVideoInvalidRequest req)
}