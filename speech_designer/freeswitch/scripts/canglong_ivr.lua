--[[
-- 作者：石晓旭
-- 日期：20190819
-- 描述：测试机器人
--]]

local http = require("socket.http")
local ltn12 = require("ltn12");
local base64 = require("mime")
local json = require("libs/json/json")
local xml = require("libs/xml/xmlSimple").newParser()


IvrContext = {
	-- 请求上下文
	m_request_url = "",
	m_request_format = "",
	m_request_query_method = "Frist",
	m_request_trigger = "",
	m_request_s_id = 0,
	m_request_sn_id = 0,
	m_request_speech_number = "",
	-- 通用配置
	m_sounds_dir = "",
	-- 智能客服
	m_default_asr_engine = "",
}


function IvrContext:new(o,side)
 	o = o or {}
	setmetatable(o, self)
	self.__index = self
 	side = side or 0

	return o	
end

function IvrContext:query_node_by_http()
	local request_body = {}
	local response_body = {}
	local request_json = ""

	request_body.SpeechNumber 	= self.m_request_speech_number
	request_body.SpeechId 		= self.m_request_s_id
	request_body.NodeParent 	= self.m_request_sn_id
	request_body.NodeTrigger 	= self.m_request_trigger
	request_body.QueryMethod        = self.m_request_query_method

	request_json 			= json.encode(request_body)

        local res, code, response_headers = http.request{
                url = self.m_request_url,
                method = "POST",
                headers =
                {
                        ["Content-Type"] = "application/json";
                        ["Content-Length"] = #request_json;
                },
                source = ltn12.source.string(request_json),
                sink = ltn12.sink.table(response_body),
        }

        if code ~= nil and type(code) ~= "string" and code >= 200 and code < 300 then
                if type(response_body) == "table" then
                        session:consoleLog("DEBUG", "IvrContext:query response body : " .. table.concat(response_body) .."\n")
                        local rc = json.decode(response_body[1])
                        return rc.Node
                else
                        session:consoleLog("ERR", "IvrContext:query response body type : " .. type(response_body) .."\n")
			session:hangup()
                        return nil
                end
        else
                session:consoleLog("ERR", "IvrContext:query failed to request : " .. code .." ... \n")
                session:hangup()
                return nil
        end

end

function IvrContext:transfer(node)
	local t_argv = nil	
	local argvs = json.decode(node.SN_Argv)
	if node.SN_Argc == 1 then
		t_argv = string.format("%s %s %s", argvs.target, argvs.dialplan, argvs.context)
	else
		session:consoleLog("ERR", "IvrContext:transfer bad arguments!\n")
		session:hangup()
	end

	session:consoleLog("DEBUG", "IvrContext:transfer arguments : " .. t_argv .. " !\n")
	session:execute("transfer", t_argv)
end

function IvrContext:play_and_detect_speech_file(node)
	local full
	if node.SN_Argc == 1 and node.SN_Argv ~= '' then
		full = string.format("%s/%s %s", self.m_sounds_dir , node.SN_File, node.SN_Argv)
	else
		full = string.format("%s/%s %s", self.m_sounds_dir, node.SN_File, self.m_default_asr_engine)
	end

	session:execute("play_and_detect_speech", full)

	local detect_result = session:getVariable('detect_speech_result')
	if xml ~= nil then
		local xml_parser_result = xml:ParseXmlText(detect_result)
		self.m_request_s_id = node.S_ID
		self.m_request_sn_id = node.SN_ID
		self.m_request_query_method = "Next"
		self.m_request_trigger = xml_parser_result.result.interpretation.input:value()
	else
		session:consoleLog("INFO", "IvrContext:play_and_detect_speech_file No result!\n")
	end
end

function IvrContext:play_and_hangup_file(node)
	session:execute("playback", table.concat({self.m_sounds_dir, "/", node.SN_File}))
	session:hangup()
end

function IvrContext:play_and_get_digits_file(node)
	local argvs = json.decode(node.SN_Argv)
	local destnum = session:playAndGetDigits(
		argvs.min, 
		argvs.max, 
		argvs.tries, 
		argvs.timeout, 
		argvs.terminators, 
		table.concat({self.m_sounds_dir, "/", node.SN_File}),
		argvs.invalid_file,
		argvs.regexp)

	if argvs.rollback == "false" then
		self.m_request_trigger = tostring(destnum)
		self.m_request_s_id = node.S_ID
		self.m_request_sn_id = node.SN_ID
		self.m_request_query_method = "Next"

		session:consoleLog("DEBUG","IvrContext: query next\n")
	elseif argvs.rollback == "true" then
		if tostring(destnum) == argvs.parent then
			self.m_request_query_method = "Parent"
			self.m_request_s_id = node.S_ID
			self.m_request_sn_id = node.SN_Parent

			session:consoleLog("DEBUG","IvrContext: query parent\n")
		elseif tostring(destnum) == argvs.frist then
			self.m_request_query_method = "Frist"

			session:consoleLog("DEBUG","IvrContext: query frist\n")
		else
			session:consoleLog("ERR", "IvrContext:play_and_get_digits_file bad input argvs.frist : " .. argvs.frist .. " argvs.parent : " .. argvs.parent .. "\n")
			session:hangup()
		end
	else
		session:consoleLog("ERR", "IvrContext:play_and_get_digits_file bad params argvs.select : " .. argvs.select .. "\n")
		session:hangup()
	end
end

function IvrContext:query()
	-- query_frist(node)
	-- self.m_request_s_id = 0 : 没有使用
	-- self.m_request_trigger = ""  : 没有使用
	-- self.m_request_sn_id = 0     : 没有使用
	-- self.m_request_speech_number : self.pre_proc 初始化，不用更新
	-- self.m_request_query_method  : self.proc 里面更新
	--
	-- query_next(node)
	-- self.m_request_speech_number : 不用更新
	-- self.m_request_trigger 	         : self.proc 里面更新
	-- self.m_request_query_method  = "Next" : self.proc 里面更新
	-- self.m_request_s_id = node.S_ID  : self.proc 里面更新
	-- self.m_request_sn_id = node.SN_ID  	 : self.proc 里面更新
	--

	-- query_parent(node)
	-- self.m_request_speech_number           : 不用更新
	-- self.m_request_trigger 	          : 没有使用
	-- self.m_request_query_method = "Parent" : self.proc 里面更新
	-- self.m_request_s_id = node.S_ID   : self.proc 里面更新
	-- self.m_request_sn_id = node.SN_Parent  : self.proc 里面更新
	--
	return self:query_node_by_http()
end



function IvrContext:proc(node)
	if type(node) == "table" then
		if next(node) then
			session:consoleLog("DEBUG", "IvrContext:proc node.SN_Action :  ".. node.SN_Action .."\n")
			if node.SN_Action == "play_and_detect_speech.file" then
				self:play_and_detect_speech_file(node)
			elseif node.SN_Action == "play_and_hangup.file" then
				self:play_and_hangup_file(node)
			elseif node.SN_Action == "play_and_get_digits.file" then
				self:play_and_get_digits_file(node)
			elseif node.SN_Action == "transfer" then
				self:transfer(node)
			else
				session:consoleLog("ERR", "IvrContext:proc failed to proc node.SN_Action unknown\n")
				session:hangup()
			end
		else
			session:consoleLog("ERR", "IvrContext:proc failed to proc node is nil\n")
		end
	else
		session:consoleLog("ERR", "IvrContext:proc failed to parse node type\n")
	end

end

function IvrContext:pre_run(...)
	session:answer()
	session:sleep(1000)

	self.m_request_query_method = "Frist"
	self.m_request_speech_number = argv[1]
	self.m_request_url = "http://192.168.43.94:8000/query"
	self.m_request_format = "{\"SpeechNumber\":\"%s\",\"SpeechId\":%d,\"NodeParent\":%d,\"NodeTrigger\":\"%s\"}"
	self.m_default_asr_engine = "detect:unimrcp:canglong-robot-unimrcp-v2 {Start-Input-Timers=false,no-input-timeout=6000,sentivity-level=0.5}one"
	self.m_sounds_dir = freeswitch.getGlobalVariable("robot_sounds_dir")	

	session:consoleLog("DEBUG", "IvrContext.m_request_url : " .. self.m_request_url .. "\n")
	session:consoleLog("DEBUG", "IvrContext.m_request_speech_number : " .. self.m_request_speech_number .. "\n")
	session:consoleLog("DEBUG", "IvrContext.m_default_asr_engine : " .. self.m_default_asr_engine .. "\n")
	session:consoleLog("DEBUG", "IvrContext.m_sounds_dir : " .. self.m_sounds_dir .. "\n")
end

function IvrContext:run(...)
	local node = self:query()

	while (session:ready()) do
		self:proc(node)

		node = self:query()
	end
end


-- lua 主函数
function main(...)
	ic = IvrContext:new()
	ic:pre_run()
	ic:run()
end

-- lua 异常处理
local status, err = pcall(main)
if status then
	session:consoleLog("DEBUG", "IvrContext:main success\n")
else
	session:consoleLog("ERR", "IvrContext:main exception : ".. err .."\n")
end

