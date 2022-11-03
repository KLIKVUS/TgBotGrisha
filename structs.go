package main

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
	Id      int     `json:"message_id"`
	User    User    `json:"from"`
	Date    int     `json:"date"`
	Chat    Chat    `json:"chat"`
	Text    string  `json:"text"`
	Sticker Sticker `json:"sticker"`
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
type Sticker struct {
	File_Id        string `json:"file_id"`
	File_Unique_Id string `json:"file_unique_id"`
}

type MainStruct struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}
type Result struct {
	Id         int    `json:"id"`
	Is_bot     bool   `json:"is_bot"`
	First_Name string `json:"first_name"`
	Username   string `json:"username"`
	// CanJoin_Groups              bool   `json:"can_join_groups"`
	// Can_Read_All_Group_Messages bool   `json:"can_read_all_group_messages"`
	// Supports_Inline_Queries     bool   `json:"supports_inline_queries"`
	Abilities []string `json:"abilities"`
}

type WeatherResponse struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
}
type Weather struct {
	Description string `json:"description"`
}
type Main struct {
	Temp_Min float32 `json:"temp_min"`
	Temp_Max float32 `json:"temp_max"`
}

type SendSticker struct {
	Chat_Id int    `json:"chat_id"`
	Sticker string `json:"sticker"`
}
