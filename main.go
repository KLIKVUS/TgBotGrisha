package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var myName string = "Гриша"

const apiUrl string = "https://api.telegram.org/" + "bot5794246977:AAE2ab_pZ97tyY1AoNHhjFqQIS0bg3ffYuA"

func main() {
	go UpdateLoop()

	router := mux.NewRouter()
	router.HandleFunc("/", Handler)

	http.ListenAndServe("localhost:8080", router)
}

type UpdateResponse struct {
	Ok     bool           `json:"ok"`
	Result []UpdateStruct `json:"result"`
}
type UpdateStruct struct {
	Id                  int     `json:"update_id"`
	Message             Message `json:"message"`
	Edited_Message      Message `json:"edited_message"`
	Channel_Post        Message `json:"channel_post"`
	Edited_Channel_Post Message `json:"edited_channel_post"`
}
type Message struct {
	Id   int    `json:"id"`
	User User   `json:"from"`
	Date int    `json:"date"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}
type SendMessage struct {
	Chat_Id             int    `json:"chat_id"`
	Text                string `json:"text"`
	Reply_To_Message_Id int    `json:"reply_to_message_id"`
}
type User struct {
	Id         int    `json:"id"`
	Is_Bot     bool   `json:"is_bot"`
	First_Name string `json:"first_name"`
	Username   string `json:"username"`
	Is_Prem    bool   `json:"is_premium"`
}
type Chat struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}

type MainStruct struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}
type Result struct {
	Id         int    `json:"id"`
	Is_bot     bool   `json:"is_bot"`
	First_Name string `json:"first_name"`
	// Username                    string `json:"username"`
	// CanJoin_Groups              bool   `json:"can_join_groups"`
	// Can_Read_All_Group_Messages bool   `json:"can_read_all_group_messages"`
	// Supports_Inline_Queries     bool   `json:"supports_inline_queries"`
	Abilities []string `json:"abilities"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var res MainStruct

	resp, err := http.Get(apiUrl + "/getMe")
	if err != nil {
		panic(err)
	}

	respBody, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &res)
	if err != nil {
		fmt.Println(err)
	}

	res.Result.Abilities = append(res.Result.Abilities, "Replying to messages /privet")

	respReady, err := json.Marshal(res.Result)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", "inline")
	w.Write(respReady)
}

func UpdateLoop() {
	lastId := 0
	for {
		lastId = Update(lastId)
		time.Sleep(5 * time.Second)
	}
}

func Update(lastId int) int {
	raw, err := http.Get(apiUrl + "/getUpdates?offset=" + strconv.Itoa(lastId))
	if err != nil {
		panic(err)
	}

	body, _ := io.ReadAll(raw.Body)

	var v UpdateResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		panic(err)
	}

	if len(v.Result) > 0 {
		ev := v.Result[len(v.Result)-1]
		txt := ev.Message.Text
		if txt == "/privet" {
			return SendMsg(lastId, ev, "Штоуж")
		}

		splitedTxt := strings.Split(txt, ",")
		splitedTxt1 := strings.Split(txt, " ")
		maybeNewName := splitedTxt1[len(splitedTxt1)-1]
		// Когдя ты "слишком умный" для regexp-a
		if splitedTxt[0] == myName {
			switch splitedTxt[1] {
			case " танцуй!":
				return SendMsg(lastId, ev, "Ты эбобо??")
			case " шути!":
				return SendMsg(lastId, ev, "Танцуют два негра и один упал.")

			// Гений 1000lvl
			case " теперь ты " + maybeNewName:
				myName = maybeNewName
				return SendMsg(lastId, ev, "Ну получаетс я "+maybeNewName)
			default:
				return SendMsg(lastId, ev, "Не понимаю((")
			}
		}
	}

	return lastId
}

func SendMsg(lastId int, ev UpdateStruct, text string) int {
	txtmsg := SendMessage{
		Chat_Id:             ev.Message.Chat.Id,
		Text:                text,
		Reply_To_Message_Id: ev.Message.Id,
	}

	bytemsg, _ := json.Marshal(txtmsg)

	_, err := http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
	if err != nil {
		fmt.Println(err)
		return lastId
	}
	return ev.Id + 1
}
