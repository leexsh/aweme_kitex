include "user.thrift"
include "base.thrift"
namespace go feed

struct BaseResp {
    1: i64 status_code
    2: string status_msg
    3: i64 service_time
}

struct Video {
    1: string video_id
    2: user.User author
    3: string play_url
    4: string cover_url
    5: i64 favourite_count
    6: i64 comment_count
    7: bool is_favourite
    8: string title
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

service FeedService {
    FeedResponse Feed(1: FeedRequest req)
}