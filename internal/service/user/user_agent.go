package user_svc

import (
	"errors"
	"strings"

	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	user_repo "github.com/sfshf/gonoweb/internal/repo/user"
	jwt_util "github.com/sfshf/gonoweb/internal/util/jwt"
)

// 新建/更新用户的代理信息
// 入参 userInfo -- 用户的登录信息；0->user_xid，1->token
func UpsertUserAgent(ip, ua, tid string, userInfo ...string) error {
	var err error
	// 如果len(userInfo)>0，则是在登录，否则只是访客行为
	var record *TUserAgent
	if len(userInfo) > 0 {
		record, err = user_repo.UserAgent_FirstDeletedByIP(ip)
	} else {
		record, err = user_repo.UserAgent_FirstByIP(ip)
	}
	if err != nil {
		return err
	}
	var uxid string
	if len(userInfo) > 0 {
		uxid = userInfo[0]
	}
	var token string
	if len(userInfo) > 1 {
		token = userInfo[1]
	}
	if record == nil {
		// 新增记录
		if err := repo.Create(&TUserAgent{
			IP:      ip,
			Ua:      ua,
			TraceID: tid,
			UserXid: uxid,
			Token:   token,
		}); err != nil {
			return err
		}
	} else {
		// 如果len(userInfo)>0，则是在登录，否则只是访客行为
		if len(userInfo) > 0 {
			// 登录，则复用之前的记录
			if err := user_repo.UserAgent_ReliveByIP(ip, &TUserAgent{
				Ua:      ua,
				TraceID: tid,
				UserXid: uxid,
				Token:   token,
			}); err != nil {
				return err
			}
		} else {
			// 访客，则更新访客信息
			if err := user_repo.UserAgent_UpdateByIP(ip, &TUserAgent{
				Ua:      ua,
				TraceID: tid,
				UserXid: uxid,
				Token:   token,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckTokenWithIP(token, ip string) error {
	err := errors.New("非法的登录token")
	token = strings.TrimPrefix(token, jwt_util.BearerPrefix)
	record, _ := user_repo.UserAgent_FirstByToken(token)
	if record == nil {
		return err
	}
	if record.IP != ip {
		return err
	}
	return nil
}
