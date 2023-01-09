include "user.thrift"
include "base.thrift"
namespace go relation

struct RelationActionRequest {
    1: string token
    2: string to_user_id
    3: string action_type
}

struct RelationActionResponse {
    1: i32 status_code
    2: string status_msg
}

struct FollowListRequest {
    1: string token
}

struct FollowListResponse {
    1: base.BaseResp base_resp
    2: list<user.User> user_list
}

struct FollowerListRequest {
    1: string token
}

struct FollowerListResponse {
   1: base.BaseResp base_resp
    2: list<user.User> user_list
}

