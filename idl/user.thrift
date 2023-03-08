include "base.thrift"
namespace go user

struct User {
    1: string user_id // 用户标识
    2: string name; // 用户名称
    3: i64 follow_count // 关注数量
    4: i64 follower_count // 粉丝数量
    5: bool is_follow // 是否关注
}

struct UserRegisterRequest {
    1: string user_name
    2: string password
}

struct UserRegisterResponse {
    1: base.BaseResp base_resp
    2: string user_id
    3: string token
}

struct UserLoginRequest {
    1: string user_name
    2: string password
}

struct UserLoginResponse {
    1: base.BaseResp base_resp
    2: string user_id
    3: string token
}

struct UserInfoRequest {
    1: string token
    2: string user_id
}

struct UserInfoResponse {
    1: base.BaseResp base_resp
    2: list<User> user
 }

struct SingleUserInfoRequest {
    1: list<string> user_ids
}

struct SingleUserInfoResponse {
    1: base.BaseResp base_resp
    2: map<string, User> users
}

struct ChangeFollowStatusRequest {
    1: string user_id
    2: string to_user_id
    3: bool isFollow
}

 service UserService {
     // 用户注册
     UserRegisterResponse Register(1: UserRegisterRequest req)
     // 用户登录
     UserLoginResponse Login(1: UserLoginRequest req)
     // 登录后获取用户信息
     UserInfoResponse UserInfo(1: UserInfoRequest req)
     // 获取单个用户信息
     SingleUserInfoResponse GetUserInfoByUserId(1: SingleUserInfoRequest req)
     // 修改关注/被关注数目
     void ChangeFollowStatus(1: ChangeFollowStatusRequest req)
 }