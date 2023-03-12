include "base.thrift"
include "feed.thrift"
namespace go favourite

struct Favourite {
    1: required string favourite_id // 收藏id
    2: string user_id
    3: string video_id
    4: string action_type // 请求类型
}

struct FavouriteActionRequest {
    1: required string token
    2: string video_id
    3: string action_type
}


struct FavouriteActionResponse {
    1: base.BaseResp base_resp
}

struct FavouriteListRequest {
    1: required string token
    2: string user_id
}

struct FavouriteListResponse {
    1: base.BaseResp base_resp
    2: map<string, feed.Video> video_list
}

struct QueryVideoIsFavouriteRequest {
    1: string user_id
    2: list<string> videos_id
}
struct QueryVideoIsFavouriteResponse {
    1: base.BaseResp base_resp
    2: map<string, bool> is_favourites
}

service FavouriteService {
     FavouriteActionResponse FavouriteAction(1: FavouriteActionRequest req)
     FavouriteListResponse FavouriteList(1: FavouriteListRequest req)
     QueryVideoIsFavouriteResponse QueryVideoIsFavourite(1: QueryVideoIsFavouriteRequest req)
}