package handlers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// func sendCard
func msgFilter(msg string) string {
	//replace @到下一个非空的字段 为 ''
	regex := regexp.MustCompile(`@[^ ]*`)
	return regex.ReplaceAllString(msg, "")

}

// Parse rich text json to text
func parsePostContent(content string) string {
	/*
		{
		    "title":"我是一个标题",
		    "content":[
		        [
		            {
		                "tag":"text",
		                "text":"第一行 :",
		                "style": ["bold", "underline"]
		            },
		            {
		                "tag":"a",
		                "href":"http://www.feishu.cn",
		                "text":"超链接",
		                "style": ["bold", "italic"]
		            },
		            {
		                "tag":"at",
		                "user_id":"@_user_1",
		                "user_name":"",
		                "style": []
		            }
		        ],
		        [
		            {
		                "tag":"img",
		                "image_key":"img_47354fbc-a159-40ed-86ab-2ad0f1acb42g"
		            }
		        ],
		        [
		            {
		                "tag":"text",
		                "text":"第二行:",
		                "style": ["bold", "underline"]
		            },
		            {
		                "tag":"text",
		                "text":"文本测试",
		                "style": []
		            }
		        ],
		        [
		            {
		                "tag":"img",
		                "image_key":"img_47354fbc-a159-40ed-86ab-2ad0f1acb42g"
		            }
		        ],
		        [
		            {
		                "tag":"media",
		                "file_key": "file_v2_0dcdd7d9-fib0-4432-a519-41d25aca542j",
		                "image_key": "img_7ea74629-9191-4176-998c-2e603c9c5e8g"
		            }
		        ],
		        [
		            {
		                "tag": "emotion",
		                "emoji_type": "SMILE"
		            }
		        ]
		    ]
		}
	*/
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)

	if err != nil {
		fmt.Println(err)
	}

	if contentMap["content"] == nil {
		return ""
	}
	var text string
	// deal with title
	if contentMap["title"] != nil && contentMap["title"] != "" {
		text += contentMap["title"].(string) + "\n"
	}
	// deal with content
	contentList := contentMap["content"].([]interface{})
	for _, v := range contentList {
		for _, v1 := range v.([]interface{}) {
			if v1.(map[string]interface{})["tag"] == "text" {
				text += v1.(map[string]interface{})["text"].(string)
			}
		}
		// add new line
		text += "\n"
	}
	return msgFilter(text)
}

func parseContent(content, msgType string) string {
	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha
	if msgType == "post" {
		return parsePostContent(content)
	}

	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	if contentMap["text"] == nil {
		return ""
	}
	text := contentMap["text"].(string)
	return msgFilter(text)
}

func processMessage(msg interface{}) (string, error) {
	msg = strings.TrimSpace(msg.(string))
	msgB, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	msgStr := string(msgB)

	if len(msgStr) >= 2 {
		msgStr = msgStr[1 : len(msgStr)-1]
	}
	return msgStr, nil
}

func processNewLine(msg string) string {
	return strings.Replace(msg, "\\n", `
`, -1)
}

func processQuote(msg string) string {
	return strings.Replace(msg, "\\\"", "\"", -1)
}

// 将字符中 \u003c 替换为 <  等等
func processUnicode(msg string) string {
	regex := regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)
	return regex.ReplaceAllStringFunc(msg, func(s string) string {
		r, _ := regexp.Compile(`\\u`)
		s = r.ReplaceAllString(s, "")
		i, _ := strconv.ParseInt(s, 16, 32)
		return string(rune(i))
	})
}

func cleanTextBlock(msg string) string {
	msg = processNewLine(msg)
	msg = processUnicode(msg)
	msg = processQuote(msg)
	return msg
}

func parseFileKey(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if contentMap["file_key"] == nil {
		return ""
	}
	fileKey := contentMap["file_key"].(string)
	return fileKey
}

func parseImageKey(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if contentMap["image_key"] == nil {
		return ""
	}
	imageKey := contentMap["image_key"].(string)
	return imageKey
}
