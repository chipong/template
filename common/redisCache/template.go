package redisCache

import (
	"context"
	"errors"
	"fmt"

	"github.com/chipong/template/common/proto"
)

func GetTemplate(uid string) ([]*proto.OZTemplate, error) {
	key := fmt.Sprintf("%s:%s:template", appName, uid)
	results, err := LRange(key, DataTTL, &proto.OZTemplate{})
	if err != nil {
		return nil, err
	}

	if len(results) == 0 || err != nil {
		return nil, errors.New("not found data")
	}

	return results, nil
}

func SetTemplate(uid string, templates []*proto.OZTemplate) error {
	if len(templates) == 0 {
		return nil
	}

	key := fmt.Sprintf("%s:%s:template", appName, uid)
	err := LDelAndPush(key, DataTTL, templates)
	if err != nil {
		return err
	}
	return nil
}

func DelTemplate(uid string) error {
	key := fmt.Sprintf("%s:%s:template", appName, uid)
	_, err := write.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return nil
}
