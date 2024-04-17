package template

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	
	"github.com/gin-gonic/gin"

	"github.com/chipong/template/common"
	"github.com/chipong/template/common/proto"
	"github.com/chipong/template/common/dynamodb"
	"github.com/chipong/template/common/log"
	"github.com/chipong/template/common/redisCache"
	"github.com/chipong/template/common/util"
	"github.com/chipong/template/template-server/router"
)

var (
	setTemplatePool = sync.Pool{
		New: func() interface{} {
			return new(SetTemplateAPI)
		},
	}
)

type SetTemplateAPI struct {
	wg				sync.WaitGroup
	// uid         	string

	req				*proto.SetTemplateReq

	templates		[]*proto.OZTemplate
	targetTemplate	*proto.OZTemplate
	isNew			bool
}

func (r *SetTemplateAPI) reset() {
	r.templates = nil
	r.targetTemplate = nil

	r.isNew = false
}

func SetTemplate(c *gin.Context) {
	r := setTemplatePool.Get().(*SetTemplateAPI)
	r.reset()
	defer setTemplatePool.Put(r)

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

func (r *SetTemplateAPI) Check(c *gin.Context) bool {
	log.Println(c, c.Request.RequestURI)
	// var err error
	// r.req.Uid, err = util.GetHeaderUid(c)
	// if err != nil {
	// 	log.Println(c, err)
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"err_code": common.ErrCodeShardIndex,
	// 		"err_msg":  err.Error()})
	// 	return false
	// }

	r.req = &proto.SetTemplateReq{}
	if err := util.Unmarshal(c.Request.Body, r.req); err != nil {
		log.Println(c, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": common.ErrCodeJSONParsing,
			"err_msg":  err.Error()})
		return false
	}

	jsonStr, _ := json.Marshal(r.req)
	log.Println(c, "request: ", r.req.Uid, string(jsonStr))

	_, err := router.TemplateTable.FindTemplate(r.req.Id)
	if err != nil {
		log.Println(c, err, r.req.Id)
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": common.ErrCodeNotFoundData,
			"err_msg":  err.Error()})
		return false
	}

	return true
}

func (r *SetTemplateAPI) Load(c *gin.Context) bool {
	chs := make([](<-chan *common.ChErrCode), 0)
	ctx, cancel := context.WithCancel(c.Request.Context())

	loadTemplateCh := util.LoadTemplate(ctx, cancel,
		func(ch chan *common.ChErrCode) {
			defer close(ch)

			for _, v := range r.templates {
				if v.Id == r.req.Id {
					r.targetTemplate = v
					r.isNew = false
				}
			}
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

func (r *SetTemplateAPI) Exec(c *gin.Context) bool {
	if r.targetTemplate == nil {
		r.targetTemplate = &proto.OZTemplate{
			Id:		r.req.Id,
			Count:	r.req.Count,
		}
		r.isNew = true

		r.templates = append(r.templates, r.targetTemplate)
	} else {
		r.targetTemplate.Count = r.req.Count
	}

	return true
}

func (r *SetTemplateAPI) Save(c *gin.Context) bool {
	defer func() {
		if cover := recover(); cover != nil {
			log.Println(c, "tx rollback err: ", cover)
			c.JSON(http.StatusBadRequest, gin.H{
				"err_code": common.ErrCodeDynamoDB,
				"err_msg":  cover})
		} else {
			r.wg.Add(1)
			go func() {
				defer r.wg.Done()
				redisCache.SetTemplate(r.req.Uid, r.templates)
			}()
		}
	}()

	tx := dynamodb.BeginPQLTransaction()

	if r.isNew {
		dynamodb.AddPQLTransaction(tx,
			dynamodb.PutTemplateTx(r.req.Uid, r.targetTemplate))
	} else {
		dynamodb.AddPQLTransaction(tx,
			dynamodb.UpdateTemplateCountTx(r.req.Uid, "TEMPLATE#"+r.targetTemplate.Id, r.targetTemplate.Count))
	}

	err := dynamodb.EndPQLTransaction(c.Request.Context(), tx)
	if err != nil {
		log.Println(c, err)
		panic(err.Error())
	}

	return true
}

func (r *SetTemplateAPI) Logging(c *gin.Context) bool {
	return true
}

func (r *SetTemplateAPI) Answer(c *gin.Context) bool {
	ans := &proto.SetTemplateAns{
		ErrCode:		common.Success,
		ErrMsg:			"success",

		Template:		r.targetTemplate,
	}
	c.JSON(http.StatusOK, ans)
	return true
}