# 抽奖

## annual 年度抽奖


 wrk -t10 -c10 -d5 http://127.0.0.1:8080/api/v1/wechat/luck
 参数：
 -t : 模拟的线程数
 -c: 连接数
 --timeout 超时时间
 -d 测试时间
 
 结果：
 Latency：响应时间

 Req/Sec：每个线程每秒钟的完成的请求数
 Avg：平均
 Max：最大
 
Stdev：标准差
Requests/sec：QPS（每秒请求数）
Transfer/sec：每秒传输数量

QPS：Queries Per Second意思是“每秒查询率”，是一台服务器每秒能够相应的查询次数，是对一个特定的查询服务器在规定时间内所处理流量多少的衡量标准。

TPS：是TransactionsPerSecond的缩写，也就是事务数/秒。它是软件测试结果的测量单位。一个事务是指一个客户机向服务器发送请求然后服务器做出反应的过程。客户机在发送请求时开始计时，收到服务器响应后结束计时，以此来计算使用的时间和完成的事务个数，

wrk 测试post
```
# 1. 编写lua脚本，填写post的数据， 如  post.lua

wrk.method = "POST"

wrk.body  = '{"userId": "10001","coinType": "GT","type": "2","amount": "5.1"}'

wrk.headers["Content-Type"] = "application/json"

function request()

  return wrk.format('POST', nil, nil, body)

end

# 2. 执行wrk，开始压力测试:

wrk -t 16 -c 100 -d 30s --latency --timeout 5s -s post.lua http://localhost:8021/m/zh/order/new 

# wrk参数用法网上很多介绍，此处不再祥述

```


// atomic.AddInt() 原子操作