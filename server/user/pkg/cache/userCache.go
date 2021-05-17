package cache

import (
	"OpenFaas-User/model"
	"OpenFaas-User/st"
	"encoding/json"

	"fmt"
	"strings"
)

type UserCacheService struct{}

func NewUserCacheService() *UserCacheService {
	return &UserCacheService{}
}

func (s *UserCacheService) GetUserCacheKey() string {
	keys := []string{
		CACHE_TAG,
		CACHE_USER,
	}
	return strings.Join(keys, SEP)
}

func (s *UserCacheService) GetUserByUsernameOrUserId(usernameOrUserId string) (user *model.User, err error) {
	var (
		data []byte
	)
	user = new(model.User)
	key := getCacheKey(s.GetUserCacheKey(), usernameOrUserId)
	if Exists(key) {
		data, err = Get(key)
		if err != nil {
			st.Debug("user cache error", err)
			return
		} else {
			_ = json.Unmarshal(data, user)
			return user, nil
		}
	}
	return user, fmt.Errorf("cache not find")
}

func (s *UserCacheService) SetUserByUsernameAndUserId(user *model.User) (err error) {
	key := getCacheKey(s.GetUserCacheKey(), user.UserName)
	err = Set(key, *user, 3600)
	if err != nil {
		return
	}
	key = getCacheKey(s.GetUserCacheKey(), user.UserID.String())
	err = Set(key, *user, 3600)
	return
}
