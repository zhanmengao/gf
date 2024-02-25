package tracetyp

type RpcOperator int32

const (
	RpcOpInit = iota
	RpcOpHandler
	RpcOpSend
	RpcOpRequest
	RpcMqProducer
	RpcMqConsumer
)

func GetRpcOpName(op RpcOperator) string {
	switch op {
	case RpcOpInit:
		return "init"
	case RpcOpHandler:
		return "handler"
	case RpcOpSend:
		return "send"
	case RpcOpRequest:
		return "request"
	case RpcMqConsumer:
		return "consumer"
	case RpcMqProducer:
		return "producer"
	}
	return ""
}

const (
	TraceErrorCode       = "err_code"
	TraceErrorMsg        = "err_msg"
	TraceReadTablePrefix = "table_%s"
	TraceReadTable       = "%s_%%s"
	TraceTableDefine     = "table_define:%s"
	TraceTableAB         = "table_ab:%s "
	TraceTableKey        = "table_key:%v "
	TraceTableSubKey     = "table_subkey:%v "
	TraceTableResult     = "table_res[%d]:%s "
	TraceRedisKey        = "redis_key"
	TraceRedisKeyFmt     = "redis_keyfmt"
	TraceRedisField      = "redis_field"
	TraceRedisFieldFmt   = "redis_fieldfmt"
	TraceRedisCache      = "redis_cache"
	TraceRedisExist      = "redis_exist"
	TraceRedisTTL        = "redis_ttl"
	TraceReq             = "request"
	TraceResult          = "result"
	TraceRemoteAddr      = "remote_addr"
	TraceRemoteSrv       = "remote_srv"
	TraceRpcIsLocal      = "islocal"
	TraceBodySize        = "body_size"
	TraceUID             = "uid"
	TraceHttpHeader      = "http_header"
	TraceMysqlTable      = "mysql_table"
	TraceMysqlSql        = "sql"
	TraceMysqlRows       = "rows"
	TracePlatform        = "platform"
	TraceChannel         = "channel"
	TraceSeqID           = "seq_id"
	TraceReqID           = "req_id"
	TraceMysqlExist      = "mysql_exist"
)
