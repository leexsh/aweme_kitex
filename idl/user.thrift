include "base.thrift"
namespace go user
struct User {
    1: string user_id
    2: string name;
    3: i64 follow_count
    4: i64 follower_count
    5: bool is_follow
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
    1: string user_name
    2: string password
}

struct UserInfoResponse {
    1: base.BaseResp base_resp
    2: list<User> user
 }

 service UserService {
     UserRegisterResponse Register(1: UserRegisterRequest req)
     UserLoginResponse Login(1: UserLoginRequest req)
     UserInfoResponse UserInfo(1: UserInfoRequest req)
 }