include "base.thrift"
include "feed.thrift"
namespace go publish

struct PublishActionRequest {
    1: string token
    2: binary data
    3: string title
    4: string user_id
}

struct PublishActionResponse {
    1: base.BaseResp base_resp
}

struct PublishListRequest {
    1: string token
    2: string user_id
}

struct PublishListResponse {
    1: base.BaseResp base_resp
    2: list<feed.Video> video_list
}

service PublishService {
    PublishActionResponse PublishAction(1: PublishActionRequest req)
    PublishListResponse PublishList(1: PublishListRequest req)
}