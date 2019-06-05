http = require("http")
json = require("json")

module = {}

function module.superNodeList()
    res=http.post("http://152.136.11.202:8888/v1/chain/get_table_rows", [[{
        "code":"eosio",
        "scope":"eosio",
        "table":"producers",
        "json": true
    }]])
    tb = json.decode(res)
    resjson = {}
    for k,v in ipairs(tb.rows) 
    do
        resjson[k] = string.gsub(v.url,":8888","",1)
    end
    return json.encode(resjson)
end

return module