-- 请求列表
requestList = {
    {
        url="http://127.0.0.1:9002/api/v0/node/id",
        method="GET"
    },
    {
        url="http://152.136.18.176:9002/api/v0/node/id",
        method="GET"
    }, 
}
-- 设置执行间隔
setInterval(3)
-- 循环添加请求
for k,v in ipairs(requestList) 
do
    addRequest(v)
end

-- 成功回调
function onSuccess (res)
    -- 发送数据给指定服务器
    err = send("http://127.0.0.1:9003/",res)
    -- 如果错误不为空 打印错误
    if err ~= nil then
        print(err)
    end
end
-- 错误回调
function onError (err)
    print(err)
end
-- 所有都完成的回调
i = 0
function onAllDone(res)
    i = i + 1
    print("脚本2 所有结果都返回了:"..i..res)
end