package constant

import "time"

var RespNil = []byte("$-1\r\n")
var RespOk = []byte("+OK\r\n")
var RespZero = []byte(":0\r\n")
var RespOne = []byte(":1\r\n")
var TtlKeyNotExist = []byte(":-2\r\n")
var TtlKeyExistNoExpire = []byte(":-1\r\n")
var ActiveExpireFrequency = 100 * time.Millisecond
var ActiveExpireSampleSize = 20
var ActiveExpireThreshold = 0.1
var DefaultBPlusTreeDegree = 4

const BfDefaultInitCapacity = 100
const BfDefaultErrRate = 0.01

const ServerStatusIdle int32 = 0
const ServerStatusShutdown int32 = 1
const ServerStatusRunning int32 = 2
