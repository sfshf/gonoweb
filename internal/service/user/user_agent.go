package user

import (
	"errors"
	"strings"

	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"github.com/sfshf/gonoweb/internal/repo/user"
	. "github.com/sfshf/gonoweb/internal/service"
	"github.com/sfshf/gonoweb/internal/util/jwt"
)

// 新建/更新用户的代理信息
// 入参 userInfo -- 用户的登录信息；0->user_xid，1->token
func UpsertUserAgent(ip, ua, tid string, userInfo ...string) error {
	var err error
	var uxid string
	if len(userInfo) > 0 {
		uxid = userInfo[0]
	}
	var token string
	if len(userInfo) > 1 {
		token = userInfo[1]
	}
	// 如果uxid非空，则是在登录，否则只是访客行为
	var record *TUserAgent
	if uxid != "" {
		record, err = user.UserAgent_FirstUnscopedByXidAndIP(uxid, ip)
	} else {
		record, err = user.UserAgent_FirstByIP(ip)
	}
	if err != nil {
		return err
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
		// 如果uxid非空，则是在登录，否则只是访客行为
		if uxid != "" {
			// 登录，则激活之前的记录
			if err := user.UserAgent_ReliveByXidAndIP(uxid, ip, &TUserAgent{
				Ua:      ua,
				TraceID: tid,
				Token:   token,
			}); err != nil {
				return err
			}
		} else {
			// 访客，则更新访客信息
			if err := user.UserAgent_UpdateByIP(ip, &TUserAgent{
				Ua:      ua,
				TraceID: tid,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckTokenWithIP(token, ip string) error {
	err := errors.New("非法的登录token")
	token = strings.TrimPrefix(token, jwt.BearerPrefix)
	record, _ := user.UserAgent_FirstByToken(token)
	if record == nil {
		return err
	}
	if record.IP != ip {
		return err
	}
	return nil
}

func ListUserAgent(page, pageSize int) ([]TUserAgent, int64, *SvcErr) {
	db := repo.GormDB.Table(TableNameTUserAgent)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	var list []TUserAgent
	if err := db.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Where(`deleted_at=0`).
		Find(&list).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	return list, total, nil
}
