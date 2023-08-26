wrk.method="GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "PostmanRuntime/7.32.3"
-- 记得修改这个，你在登录页面登录一下，然后复制一个过来这里
wrk.headers["Authorization"]="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTMwMTk2NzQsIlVzZXJJRCI6MSwiVXNlckFnZW50IjoiQXBhY2hlLUh0dHBDbGllbnQvNC41LjEzIChKYXZhLzE3LjAuNikifQ.f05z0K5Q9St9CBnDtqWY98WU0Xv3MZDoRAGlbcQaqpM"