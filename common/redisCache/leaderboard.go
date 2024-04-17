package redisCache

import (
	"fmt"
	"math"
	"strconv"

	"github.com/chipong/template/common/proto"
)

/*
	OZ의 Leaderboard는 score를 -로 저장하여
	높은 점수의 유저가 상위 랭크로 바로 등록되도록 구성
	
*/

func SetRanker(leaderboardType, uid string, score int64, at int32) (int64, error) {
	key := fmt.Sprintf("%s:leaderboard:%s", appName, leaderboardType)
	compactScore := compactScore(score, at)
	rank, err := ZAdd(key, DataTTL, uid, compactScore)
	if err != nil {
		return 0, err
	}
	fmt.Printf("setRank uid : %s, rank : %d\n", uid, rank)
	return rank, nil
}

func GetRank(leaderboardType, uid string) (int64, error) {
	key := fmt.Sprintf("%s:leaderboard:%s", appName, leaderboardType)
	result, err := ZRank(key, DataTTL, uid)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func GetTargetRanker(leaderboardType, uid string) (*proto.Ranker, error) {
	key := fmt.Sprintf("%s:leaderboard:%s", appName, leaderboardType)
	targetRank, err := ZRank(key, DataTTL, uid)
	if err != nil {
		return nil, err
	}

	result, err := ZRange(key, DataTTL, (targetRank - 1), (targetRank - 1))
	if err != nil {
		return nil, err
	}

	tempScore, err := strconv.Atoi(result.([]interface{})[1].(string))
	if err != nil {
		return nil, err
	}
	targetScore, _ := uncompactScore(int64(tempScore))

	target := &proto.Ranker{
		Uid: 	uid,
		Rank: 	targetRank,
		Score: 	int64(targetScore),
	}

	return target, nil
}

func GetRankerList(leaderboardType, uid string, start, end int64) ([]*proto.Ranker, error) {
	key := fmt.Sprintf("%s:leaderboard:%s", appName, leaderboardType)
	results, err := ZRange(key, DataTTL, (start - 1), (end - 1))
	if err != nil {
		return nil, err
	}

	ranks := make([]*proto.Ranker, 0)
	rankCount := start
	for i := 0; i < len(results.([]interface{})); i += 2 {
		score, err := strconv.Atoi(results.([]interface{})[i + 1].(string))
		if err != nil {
			return nil, err
		}

		rankPoint, _ := uncompactScore(int64(score))
		temp := &proto.Ranker{
			Uid: results.([]interface{})[i].(string),
			Rank: int64(rankCount),
			Score: int64(rankPoint),
		}
		ranks = append(ranks, temp)
		rankCount++
	}

	// // my rank 조회
	// myRank, err := GetRank(leaderboardType, uid)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// myResult, err := ZRange(key, DataTTL, (myRank - 1), (myRank - 1))
	// if err != nil {
	// 	return nil, nil, err
	// }
	// tempScore, err := strconv.Atoi(myResult.([]interface{})[1].(string))
	// if err != nil {
	// 	return nil, nil, err
	// }

	// myScore, _ := uncompactScore(int64(tempScore))
	// myInfo := &proto.Ranker{
	// 	Uid: uid,
	// 	Rank: myRank,
	// 	Score: int64(myScore),
	// }

	return ranks, nil
}

func DeleteRanker(leaderboardType, uid string) error {
	key := fmt.Sprintf("%s:leaderboard:%s", appName, leaderboardType)
	_, err := ZRem(key, uid)
	if err != nil {
		return err
	}

	return nil
}

func DeleteLeaderboard(leaderboardType string) error {
	key := fmt.Sprintf("%s:leaderboard:%s", appName, leaderboardType)
	err := DelKeys([]string{key})
	if err != nil {
		return err
	}

	return nil
}

/*
compact/uncompact를 사용해야하는 상황이라 ZincrBy 사용 보류
func AddScore(leaderboardType, uid string, score int64, at int32) (int64, error) {
	key := fmt.Sprintf("%s:leaderboard:%s", appName, leaderboardType)
	compactScore := compactScore(score, at)
	result, err := ZincrBy(key, DataTTL, uid, compactScore, at)
	if err != nil {
		return 0, err
	}

	return result, nil
}
*/

// score가 같으면 timestamp로 sort
func compactScore(score int64, at int32) int64 {
	var Score int64
	Score = int64(score) << 32
	Score = Score + int64(math.MaxUint32-uint32(at))
	
	return -Score
}

func uncompactScore(score int64) (int64, int32) {
	score *= -1
	At := math.MaxUint32 - uint32(score)
	Score := int64(score >> 32)
	return Score, int32(At)
}