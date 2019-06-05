-- eos统计脚本
-- load("./request/eos.lua")
superListMod = require("example/node_list/")
res = superListMod.superNodeList()

print(res)
