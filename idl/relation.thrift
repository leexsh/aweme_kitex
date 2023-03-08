include "user.thrift"
include "base.thrift"
namespace go relation

struct RelationActionRequest {
    1: string token
    2: string to_user_id
    3: string action_type
}

struct RelationActionResponse {
    1: base.BaseResp base_resp
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

// 查询两个用户关系的请求
struct QueryRelationRequest {
    1: string userId
    2: string to_userId
    3: bool is_follow
}

// 查询两个用户关系的相应
struct QueryRelationResponse {
    1: base.BaseResp base_resp
    2: bool is_follow
}

service RelationService {
    // relation操作
    RelationActionResponse RelationAction(1: RelationActionRequest req)
    // 获取关注列表
    FollowListResponse FollowList(1: FollowListRequest req)
    // 获取粉丝列表
    FollowerListResponse FollowerList(1: FollowerListRequest req)
    // 查询是否关注
    QueryRelationResponse QueryRelation(1: QueryRelationRequest req)
}