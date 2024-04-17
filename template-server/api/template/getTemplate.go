package template

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	
	"github.com/gin-gonic/gin"

	"github.com/chipong/template/common"
	"github.com/chipong/template/common/proto"
	"github.com/chipong/template/common/log"
	"github.com/chipong/template/common/util"
)

var (
	getTemplatePool = sync.Pool{
		New: func() interface{} {
			return new(GetTemplateAPI)
		},
	}
)

type GetTemplateAPI struct {
	wg          sync.WaitGroup
	uid         string
	
	req         *proto.GetTemplateReq

	templates	[]*proto.OZTemplate
}

func (r *GetTemplateAPI) reset() {
	r.templates = nil
}

func GetTemplate(c *gin.Context) {
	r := getTemplatePool.Get().(*GetTemplateAPI)
	r.reset()
	defer getTemplatePool.Put(r)

	if !r.Check(c) {
		return
	}

	if !r.Load(c) {
		return
	}

	if !r.Exec(c) {
		return
	}

	if !r.Save(c) {
		return
	}

	if !r.Logging(c) {
		return
	}

	if !r.Answer(c) {
		return
	}

	r.wg.Wait()
}

func (r *GetTemplateAPI) Check(c *gin.Context) bool {
	log.Println(c, c.Request.RequestURI)

	r.req = &proto.GetTemplateReq{}
	if err := util.Unmarshal(c.Request.Body, r.req); err != nil {
		log.Println(c, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": common.ErrCodeJSONParsing,
			"err_msg":  err.Error()})
		return false
	}

	jsonStr, _ := json.Marshal(r.req)
	log.Println(c, "request: ", r.req.Uid, string(jsonStr))
	return true
}

func (r *GetTemplateAPI) Load(c *gin.Context) bool {
	chs := make([](<-chan *common.ChErrCode), 0)
	ctx, cancel := context.WithCancel(c.Request.Context())

	loadTemplateCh := util.LoadTemplate(ctx, cancel,
		func(ch chan *common.ChErrCode) {
			defer close(ch)
		},
		r.req.Uid, &r.templates, false)
	chs = append(chs, loadTemplateCh)

	for _, ch := range chs {
		errCh := <-ch
		if errCh != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err_code": errCh.ErrCode,
				"err_msg":  errCh.Err.Error()})
			return false
		}
	}

	return true
}

func (r *GetTemplateAPI) Exec(c *gin.Context) bool {
	return true
}

func (r *GetTemplateAPI) Save(c *gin.Context) bool {
	return true
}

func (r *GetTemplateAPI) Logging(c *gin.Context) bool {
	return true
}

func (r *GetTemplateAPI) Answer(c *gin.Context) bool {
	ans := &proto.GetTemplateAns{
		ErrCode:		common.Success,
		ErrMsg:			"success",

		Templates:		r.templates,
	}
	c.JSON(http.StatusOK, ans)
	return true
}