package leaderboard

import (
	"fmt"
	"log"

	"github.com/jinzhu/copier"

	"github.com/chipong/template/common/proto"
	"github.com/chipong/template/common/util"
	"github.com/chipong/template/common/redisCache"
)

const (
	URI_UPDATE = "/oz/game/leaderboard/rank/update"
)
/*
// @deprecated
func UpdateRanker(c *gin.Context, addr string, leaderboardType oz.LeaderboardType_T, score int64) (*oz.OZRanker, error) {
	uid, err := util.GetHeaderUid(c)
	if err != nil {
		return nil, err
	}
	ctx, err := RankUpdateRequestFactory(c, uid, leaderboardType, score, int32(util.OzNowUnix()))
	if err != nil {
		return nil, err
	}

	res, err := InnerRouter(uid, addr, ctx, false)
	if err != nil {
		return nil, err
	}

	temp, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	ans := &oz.RankUpdateAns{}
	err = json.Unmarshal(temp, ans)

	log.Printf("res : %v\n", res)
	log.Println("res cast result : ", ans)

	return ans.MyRanker, nil
}

// @deprecated
func RankUpdateRequestFactory(c *gin.Context, uid string, leaderboardType oz.LeaderboardType_T, score int64, at int32) (*gin.Context, error) {
	reqBody := &oz.RankUpdateReq{
		LeaderboardType: 	leaderboardType,
		Score: 				score,
	}
	
	reqMarshal, err := json.Marshal(reqBody)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ctx := c.Copy()

	// Using Deep Copy - Copy().Request.Clone()
	ctx.Request = c.Copy().Request.Clone(context.Background())
	ctx.Request.RequestURI = URI_UPDATE
	ctx.Request.Body = ioutil.NopCloser(ioutil.NopCloser(bytes.NewReader(reqMarshal)))

	return ctx, nil
}

// @deprecated
func InnerRouter(uid string, addr string, c *gin.Context, withSend bool) (map[string](interface{}), error) {
	if addr == "" {
		return nil, fmt.Errorf("not exist addr")
	}

	uri := ""
	if strings.HasPrefix(addr, "http") {
		uri = addr + c.Request.RequestURI
	} else {
		uri = "http://" + addr + c.Request.RequestURI
	}

	uri = strings.Replace(uri, "/v1", "", 1)
	uri = strings.Replace(uri, "/game", "", 1)

	resp, err := util.RequestHttpWithContext(c, uid, uri, "POST", c.Request.Body)
	if err != nil && withSend  {
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": common.ErrCodeBadRequest,
			"err_msg":  err.Error()})
		return nil, err
	}

	ans := make(map[string](interface{}))
	err = json.Unmarshal(resp, &ans)
	if err != nil && withSend {
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": common.ErrCodeBadRequest,
			"err_msg":  err.Error()})
		return nil, err
	}

	if withSend {
		c.JSON(http.StatusOK, ans)
	}
	return ans, nil
}
*/

func UpdateScore(uid string, leaderboard_type oz.LeaderboardType_T, score int64) (*oz.OZRanker, oz.LeaderboardUpdateStatus_T, error) {
	status := oz.LeaderboardUpdateStatus_NONE

	leaderboardType := oz.LeaderboardType_T(oz.LeaderboardType_T_value[oz.LeaderboardType_T_name[int32(leaderboard_type)]])
	if leaderboardType == oz.LeaderboardType_NONE || leaderboardType == oz.LeaderboardType_MAX {
		log.Println("LeaderboardType incorrect(NONE or MAX)")
		return nil, status, fmt.Errorf("LeaderboardType incorrect(NONE or MAX)")
	}

	myRanker, _ := redisCache.GetTargetRanker(leaderboardType.String(), uid)
	
	status = oz.LeaderboardUpdateStatus_CHANGED
	if myRanker == nil {
		status = oz.LeaderboardUpdateStatus_NEW
	} else {
		prevMyRanker := &oz.OZRanker{}
		copier.CopyWithOption(prevMyRanker, myRanker, copier.Option{DeepCopy: true})

		if score == myRanker.Score {
			status = oz.LeaderboardUpdateStatus_UNCHANGED
			return nil, status, nil
		}
	}

	at := int32(util.OzNowUnix())
	// TODO leaderboard_type 변경 필요
	rank, err := redisCache.SetRanker(leaderboardType.String(), uid, score, at)
	if err != nil {
		log.Println(err, "failed set score", oz.LeaderboardType_STAGE.String(), uid)
		return nil, oz.LeaderboardUpdateStatus_NONE, err
	}

	if myRanker == nil {
		myRanker = &oz.OZRanker{
			Uid: uid,
			Rank: rank,
			Score: score,
			UpdateAt: at,
		}
	} else {
		myRanker.Rank = rank
		myRanker.Score = score
		myRanker.UpdateAt = at
	}
	return myRanker, status, nil
}