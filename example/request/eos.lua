-- 引用json解析模块
json = require("json")
-- 此脚本获取eos统计数据
-- 请求列表
requestList = {
    -- 获取矿工总数
    {
        url="http://152.136.11.202:8888/v1/chain/get_table_rows",
        method="POST",
        body=[[{
            "code":"hddpool12345",
            "scope":"hddpool12345",
            "table":"gusercount",
            "json":true
        }]]
    },
    -- 获取持币用户数
    {
        url="http://152.136.11.202:8888/v1/chain/get_table_rows",
        method="POST",
        body=[[{
            "code":"eosio",
            "scope":"eosio",
            "table":"gcount",
            "json":true
        }]]
    },  
    -- 测试调用本地接口
    {
        url="http://127.0.0.1:9002/api/v0/node/id",
        method="GET",
    }  
}
-- 设置执行间隔 单位为秒 60*60为1小时
setInterval(6)
-- 循环添加请求
for k,v in ipairs(requestList)
do
    addRequest(v)
end

-- 成功回调
function onSuccess (res)
    print("有一条结果返回了:"..res)
end
-- 错误回调
function onError (err)
    perror("error:"..err)
end
-- 所有都完成的回调
i = 0
function onAllDone(res)
    i = i + 1
    print("所有结果都返回了:"..i..res)
    -- -- 发送数据给指定服务器
    -- err = send("http://127.0.0.1:9003/",res)
    -- -- 如果错误不为空 打印错误
    -- if err ~= nil then
    --     print(err)
    -- end 
    

    -- json解码示例
    resj = json.decode(res)
    print(resj[1]["id"])
end