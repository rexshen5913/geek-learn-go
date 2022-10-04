-- example HTTP POST script which demonstrates setting the
-- HTTP method, body, and adding a header

math.randomseed(os.time())

wrk.method = "POST"
--wrk.body   = "{\"email\": \"1230@qq.com\",\"password\":\"25d55ad283aa400af464c76d713c07ad\"}"
wrk.body   = "{\"email\": \"1230@qq.com\",\"password\":\"5416d7cd6ef195a0f7622a9c56b55e84\"}"
wrk.headers["Content-Type"] = "application/json"