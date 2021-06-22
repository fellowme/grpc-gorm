# grpc-gin-gorm
grpc-gin-gorm


# 1初步
现在只有grpc 需要添加gin 作为中台   
采用前端 ---->>>gin中台---->>>grpc后端操作数据
采用追踪链 日志记录慢接口的数据  进行优化
# 新增grpc deadline

# 新增拦截器 进行链式操作
# 新增etcd 负载均衡  如果采用k8s 部署 不需要etcd 负载均衡
https://juejin.cn/post/6844903678269210632 etcd 用户验证
http://thesecretlivesofdata.com/raft/   etcd raft算法演示

