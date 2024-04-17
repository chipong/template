package util

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/chipong/template/common"
	"github.com/chipong/template/common/dynamodb"
	"github.com/chipong/template/common/proto"
	"github.com/chipong/template/common/redisCache"
)
func LoadTemplate(ctx context.Context, cancel context.CancelFunc, checkFunc func(chan *common.ChErrCode),
	uid string, templates *[]*oz.OZTemplate, cache bool) <-chan *common.ChErrCode {
	job := func(ch chan *common.ChErrCode) {
		var err error
		*templates, err = redisCache.GetTemplate(uid)
		if err != nil || *templates == nil {
			*templates, err = dynamodb.GetTemplates(ctx, uid)
			if err != nil {
				log.Println(err)
				ch <- &common.ChErrCode{
					ErrCode: common.ErrCodeDynamoDB,
					Err:     err,
				}
				close(ch)
				cancel()
				return
			}

			log.Println("template loaded")
			if cache {
				redisCache.SetTemplate(uid, *templates)
			}
		}

		if checkFunc != nil {
			checkFunc(ch)
			return
		}

		close(ch)
	}
	return GoRoutineJob(ctx, cancel, job)
}

func GoRoutineJob(ctx context.Context,
	cancel context.CancelFunc,
	job func(chan *common.ChErrCode)) <-chan *common.ChErrCode {

	ch := make(chan *common.ChErrCode)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				job(ch)
				return
			}
		}
	}()
	return ch
}

func GoRoutineWaitJob(ctx context.Context,
	cancel context.CancelFunc,
	wait <-chan *common.ChErrCode,
	job func(chan *common.ChErrCode)) <-chan *common.ChErrCode {

	ch := make(chan *common.ChErrCode)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-wait:
				job(ch)
				return
			}
		}
	}()
	return ch
}

func GoId() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
