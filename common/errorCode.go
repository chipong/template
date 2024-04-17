package common

// Error Code ...
const (
	Success               = 0    // 성공 ...
	ErrCodeUnauthorized   = 1001 // 인증 실패 ...
	ErrCodeBadRequest     = 1002 // 잘못된 요청
	ErrCodeJSONParsing    = 1003 // json 파싱 실패
	ErrCodeInvalidToken   = 1004 // 무효한 토큰
	ErrCodeBanAccount     = 1005 // 벤 계정 요청
	ErrCodeSessionKey     = 1006 // 세션 키 만료 혹은 잘못된 세션키
	ErrCodeInternalServer = 1007 // 내부 서버간 오류
	ErrCodeNotFoundData   = 1008 // 데이터 획득 실패
	ErrCodeInvalidCoin    = 1009 // 잘못된 코인
	ErrCodeFirebase       = 1010 // 파이어베이스 인증 오류
	ErrCodeDynamoDB       = 1011 // 다이나모디비 서버 오류
	ErrCodeRedis          = 1012 // 레디스 서버 오류
	ErrCodeLimitOver      = 1013 // 제한 항목 초과
	ErrCodeNotEnoughAsset = 1014 // 제화가 부족
	ErrCodeAppleIAP       = 1015 // 애플 영수증 확인
	ErrCodeGoogleIAP      = 1016 // 구글 영수증 확인
	ErrCodeGRPC           = 1017 // GRPC 통신 오류
	ErrCodeMysqlDB        = 1018 // MysqlDB 서버 오류
	ErrCodeWebsock        = 1019 // Websocket 오류
	ErrCodeTransaction    = 1020 // Transaction 오류
	ErrCodeAlreadyProc    = 1021 // 이미 처리중인 작업이 있음
	ErrCodeRandomBox      = 1022 // 랜덤 상자 오류
	ErrCodeShardIndex     = 1023 // 사딩 인덱스 오류
	ErrCodeAlreadyPlay    = 1024 // 이미 플레이중인 데이터 있음

)

var errCodeText = map[int]string{
	Success:               "Success",
	ErrCodeUnauthorized:   "ErrCodeUnauthorized",
	ErrCodeBadRequest:     "ErrCodeBadRequest",
	ErrCodeJSONParsing:    "ErrCodeJSONParsing",
	ErrCodeInvalidToken:   "ErrCodeInvalidToken",
	ErrCodeBanAccount:     "ErrCodeBanAccount",
	ErrCodeSessionKey:     "ErrCodeSessionKey",
	ErrCodeInternalServer: "ErrCodeInternalServer",
	ErrCodeNotFoundData:   "ErrCodeNotFoundData",
	ErrCodeInvalidCoin:    "ErrCodeInvalidCoin",
}

// ErrorString ...
func ErrorString(code int) string {
	return errCodeText[code]
}
