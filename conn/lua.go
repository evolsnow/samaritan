package conn

import "github.com/garyburd/redigo/redis"

var createProductLua = `local product_id = redis.call("INCR", "ProductIncrId")
redis.call("HMSET", "product:"..product_id,
"Name", KEYS[1],
"CreateTime", KEYS[2],
"ReviewTime", KEYS[3],
"PublishTime", KEYS[4],
"InterestStart", KEYS[5],
"InterestEnd", KEYS[6],
"Period", KEYS[7],
"FundsAmount", KEYS[8],
"MinPurchase", KEYS[9],
"SalseAmount",KEYS[10],
"InterestYear", KEYS[11],
"RepaymentType", KEYS[12],
"Category", KEYS[13],
"Status", KEYS[14],
"IntegerRequired", KEYS[15],
"RequiredLv", KEYS[16],
"BorrowId", KEYS[17],
"CreateStuffId", KEYS[18],
"ReviewStuffId", KEYS[19])
return product_id
`
var CreateProductScript = redis.NewScript(19, createProductLua)
