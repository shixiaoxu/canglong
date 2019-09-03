--[[
-- 作者：石晓旭
-- 日期：20190819
-- 描述：测试机器人
--]]

local http = require("socket.http")
local ltn12 = require("ltn12");
local base64 = require("mime")
local json = require("libs/json/json")

local ttsDefaultParams  = "say:unimrcp:canglong-robot-unimrcp-v2:"
local asrDefaultParamsOfPlayFile = "detect:unimrcp:canglong-robot-unimrcp-v2 {Start-Input-Timers=false,no-input-timeout=6000,sentivity-level=0.5}one"
local robotSoundsdir =  freeswitch.getGlobalVariable("robot_sounds_dir")
local recvText = ""
local parentNode = 0
local node = nil

function canglong_query_node(node)
	local request_body = nil
	local response_body = {}
	if node == nil then
		session:consoleLog("DEBUG", "[canglong.robot.query.node -> ] request body :  111111 \n")
		request_body = string.format("{\"SpeechNumber\":\"%s\",\"SpeechId\":%d,\"NodeParent\":%d,\"NodeBody\":\"%s\"}", "911", 0, 0, "")
	else
		session:consoleLog("DEBUG", "[canglong.robot.query.node -> ] request body :  222222 ".. type(recvText) .." \n")
		request_body = string.format("{\"SpeechNumber\":\"%s\",\"SpeechId\":%d,\"NodeParent\":%d,\"NodeBody\":\"%s\"}", "911", node.S_ID , node.SN_ID, recvText)
	end
	
	session:consoleLog("DEBUG", "[canglong.robot.query.node -> ] request body : " .. request_body .."\n")

	local res, code, response_headers = http.request{
		url = "http://192.168.43.94:8000/query",
    		method = "POST",
    		headers =
    		{
        		["Content-Type"] = "application/json";
        		["Content-Length"] = #request_body;
    		},
    		source = ltn12.source.string(request_body),
    		sink = ltn12.sink.table(response_body),
	}

	if code ~= nil and type(code) ~= "string" and code >= 200 and code < 300 then
		if type(response_body) == "table" then
			session:consoleLog("DEBUG", "[canglong.robot.query.node -> ] response body : " .. table.concat(response_body) .."\n")
			local rc = json.decode(response_body[1])
			node = rc.Node
			return node
		else
			session:consoleLog("ERR", "[canglong.robot.query.node -> ] response body type : " .. type(response_body) .."\n")
			return nil
		end
	else 
		session:consoleLog("ERR", "[canglong.robot.query.node -> ] Failed to request : " .. code .." ... \n")
		return nil
	end
end

function canglong_proc_node(node)
	if type(node) == "table" then
		if next(node) then
			session:consoleLog("DEBUG", "[canglong.robot.porc.node -> ] node.SN_Action :  ".. node.SN_Action .."\n")
			if node.SN_Action == "play_and_detect_speech.file" then
				local full
				if node.SN_Argc == 1 and node.SN_Argv ~= '' then
					full = string.format("%s/%s %s", robotSoundsdir , node.SN_File, node.SN_Argv)
				else
					full = string.format("%s/%s %s", robotSoundsdir, node.SN_File, asrDefaultParamsOfPlayFile)
				end

				session:execute("play_and_detect_speech", full)

				local xml = session:getVariable('detect_speech_result')
				if xml ~= nil then
					recvText = base64.b64(xml)
			        	session:consoleLog("INFO", "[canglong.robot.run -> ] ".. xml .."\n")
				else
			        	session:consoleLog("INFO", "[canglong.robot.run -> ] No result!\n")
				end
			elseif node.SN_Action == "play_and_hangup.file" then
				session:execute("playback", table.concat({robotSoundsdir, "/", node.SN_File}))
				session:hangup()
			else
				session:consoleLog("ERR", "[canglong.robot.porc.node -> ] Failed to proc node.SN_Action unknown\n")
			end
		else
			session:consoleLog("ERR", "[canglong.robot.porc.node -> ] Failed to proc node is nil\n")
		end
	else
		session:consoleLog("ERR", "[canglong.robot.porc.node -> ] Failed to parse node type\n")
	end
end

-- lua 机器人处理
function canglong_robot_run()
	while (session:ready()) do
		node = canglong_query_node(node)
		if node ~= nil then
			local rc = canglong_proc_node(node)
		else
			session:consoleLog("ERR", "[canglong.robot.run -> ] Failed to query node\n")
		end
	end
end

-- lua 主函数
function main(...)
	session:answer()
	session:sleep(1000)
	
	session:consoleLog("DEBUG", "[canglong.robot.main -> ] run\n")
	canglong_robot_run()
end

-- lua 异常处理
local status, err = pcall(main)
if status then
	session:consoleLog("DEBUG", "[canglong.robot.post.main -> ] success\n")
else
	session:consoleLog("ERR", "[canglong.robot.post.main -> ] exception : ".. err .."\n")
end

