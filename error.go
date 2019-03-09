package main

import (
	"encoding/json"
	"log"
)

type returnCode struct {
	sendCloudV1
	Result     bool        `json:"result"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Info       interface{} `json:"info"`
}

type sendCloudV1 struct {
	Message     string   `json:"message"`
	EmailIDList []string `json:"email_id_list,omitempty"`
	Errors      []string `json:"errors,omitempty"`
}

func newReturnCode(c int) *returnCode {
	var code = map[int]string{
		200:   "请求成功",
		40011: "email不能为空",
		40012: "email格式非法",
		40209: "subject不能为空",
		40210: "subject格式非法",
		40211: "html不能为空",
		40212: "html格式非法",
		40213: "text不能为空",
		40214: "text格式非法",
		40501: "name不能为空",
		40801: "发信人地址from不能为空",
		40802: "发信人地址from格式错误",
		40803: "发信人名称fromName不能为空",
		40804: "发信人名称fromName格式错误",
		40805: "收件人地址不能为空",
		40806: "收件人地址数组中, 存在非法地址",
		40807: "收件人地址的数目不能超过100",
		40808: "邮件主题subject不能为空",
		40809: "邮件主题subject格式错误",
		40810: "回复地址replyto不能为空",
		40811: "回复地址replyto格式错误",
		40818: "attachments不能为空",
		40819: "附件过大",
		40830: "plain内容不能为空",
		40831: "plain内容格式错误",
		40852: "cc地址不能为空",
		40853: "cc地址格式错误",
		40854: "CC地址的数目不能超过100",
		40855: "bcc地址不能为空",
		40856: "bcc地址格式错误",
		40857: "BCC地址的数目不能超过100",
		40858: "respEmailId不能为空",
		40859: "respEmailId格式错误",
		40860: "gzipCompress不能为空",
		40861: "gzipCompress格式错误",
		40862: "to中有格式错误的地址列表",
		40863: "to中有不存在的地址列表",
		40864: "地址列表的数目不能超过5",
		40865: "html解压失败",
		40866: "plain解压失败",
		40867: "处理附件发生异常",
		40868: "headers不能为空",
		40869: "headers格式错误",
		40870: "html和plain不能同时为空",
		40871: "html格式错误",
		40901: "邮件发送失败",
		40902: "邮件处理发生未知异常",
		40903: "邮件发送成功",
		41001: "name不能为空串",
		41002: "name的长度应该为1-250个字符",
		41003: "name不符合域名规则",
		41004: "newName不能为空串",
		41005: "newName的长度应该为1-250个字符",
		41006: "newName不符合域名规则",
		41007: "type不能为空串",
		41008: "type不符合规则",
		41009: "verify不能为空串",
		41010: "verify不符合规则",
		41011: "verify解析错误",
		41013: "name参数错误, 多个域名",
		41101: "emailType不能为空串",
		41102: "emailType不符合规则",
		41103: "cType不能为空串",
		41104: "cType不符合规则",
		41105: "domainName不能为空串",
		41106: "domainName不符合规则",
		41107: "domainName的长度应该为1-250个字符",
		41108: "domainName所属的域名不存在",
		41109: "用户信息不存在",
		41110: "name不能为空串",
		41111: "name不符合规则, name的长度为6-32的字符串, 只能含有(A-Z,a-z,0-9,_)",
		41112: "apiUser不能超过10个",
		41113: "open不能为空串",
		41114: "open不符合规则",
		41115: "click不能为空串",
		41116: "click不符合规则",
		41117: "unsubscribe不能为空串",
		41118: "unsubscribe不符合规则",
		41119: "apiUser创建失败",
		49901: "url格式错误",
		49902: "http请求执行异常",
		49903: "http请求执行失败",
		49904: "http请求执行成功",
		49905: "http返回结果解析错误",
		49906: "http其他错误",
		501:   "服务器异常",
		6001:  "你没有权限访问",
		99999: "未知错误",
	}
	rc := &returnCode{
		StatusCode: c,
	}
	if c == 200 {
		rc.Result = true
	}
	var isExist bool
	if rc.Message, isExist = code[c]; !isExist {
		rc.StatusCode = 99999
		rc.Message = code[99999]
	}
	if rc.StatusCode == 200 {
		rc.sendCloudV1.Message = "success"
		rc.sendCloudV1.EmailIDList = append(rc.sendCloudV1.EmailIDList, "")
	} else {
		rc.sendCloudV1.Message = "error"
		rc.sendCloudV1.Errors = append(rc.sendCloudV1.Errors, rc.Message)
	}
	return rc
}

func (rc *returnCode) Error() string {
	b, err := json.Marshal(&rc.sendCloudV1)
	if err != nil {
		log.Println("returnCode err: ", err)
		return ""
	}
	return string(b)
}

func (rc *returnCode) Bytes() []byte {
	return []byte(rc.Error()+"\n")
}

func (rc *returnCode) code() int {
	return rc.StatusCode
}
