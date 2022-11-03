package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattn/go-sqlite3"
)

const tgApiUrl string = "https://api.telegram.org/bot5700737410:AAEk9gWwDNhR4N_Z7JebIG8dc4fK6CgDdz8"

func main() {
	fmt.Println("> The bot is running!")

	sql.Register(
		"sqlite3_with_extensions",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"sqlite3_mod_regexp",
			},
		},
	)

	go UpdateLoop()

	router := mux.NewRouter()
	router.HandleFunc("/api", Handler)
	router.HandleFunc("/botName", NameHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe("localhost:8080", router)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var res MainStruct

	Ping()

	resp, err := http.Get(tgApiUrl + "/getMe")
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

func NameHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(GetName()))
}

func GetName() string {
	db, err := sql.Open("sqlite3", "file:database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var gotname string
	var resp sql.NullString
	err = db.QueryRow("SELECT name FROM bot_status").Scan(&resp)

	if err != nil {
		fmt.Println(err)
	}

	if resp.Valid {
		gotname = resp.String
	}

	return gotname
}

func AuthCheck(w http.ResponseWriter, _ *http.Request) {

}

func Login(w http.ResponseWriter, _ *http.Request) {

}

func Register(w http.ResponseWriter, _ *http.Request) {

}

func UpdateLoop() {
	db, err := sql.Open("sqlite3", "file:database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	lastId := 0
	for {
		newId := Update(lastId)
		if lastId != newId {
			lastId = newId

			str := fmt.Sprintf("UPDATE bot_status SET lastid = %d;", lastId)
			db.Exec(str)
		}
		time.Sleep(1 * time.Second)
	}
}

func Update(lastId int) int {
	raw, err := http.Get(tgApiUrl + "/getUpdates?offset=" + strconv.Itoa(lastId))
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
		var myFirstName = "Жопер" //GetName()

		ev := v.Result[len(v.Result)-1]
		// fmt.Println(ev.Message)
		txt := ev.Message.Text
		if txt == "/privet" {
			return SendMsg(lastId, ev, "Штоуж")
		}

		splitedTxt := strings.Split(txt, ", ")
		if splitedTxt[0] == myFirstName || splitedTxt[0] == fmt.Sprintf("@%s", GetMyMainName()) {
			switch strings.Split(splitedTxt[1], ": ")[0] {
			case "танцуй!":
				return SendMsg(lastId, ev, "Ты эбобо??")
			case "шути!":
				return SendMsg(lastId, ev, "Танцуют два негра и один упал.")
			case "дай погоду!":
				res, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=55.751244&lon=37.618423&lang=ru&appid=1759cf4bcb07551210be50cfe44c5c06")
				if err != nil {
					panic(err)
				}

				body, _ := io.ReadAll(res.Body)

				var v WeatherResponse
				err = json.Unmarshal(body, &v)
				if err != nil {
					panic(err)
				}

				weather := v.Weather[0]
				weatherMessage := fmt.Sprintf("Погода: \n%s", weather.Description)

				temp := v.Main
				tempMessage := fmt.Sprintf("Температура: \nmax: %f, min: %f", (temp.Temp_Max-32)*5/9, (temp.Temp_Min-32)*5/9)

				return SendMsg(lastId, ev, weatherMessage+"\n\n"+tempMessage)
			case "позови Олежу!":
				return SendMsg(lastId, ev, "Олежа, расскажи анекдот")
			case "дай конфиг препода!":
				return SendStick(
					SendMsg(lastId, ev, "1. Клава от бляди \n2. Крутой комп"),
					ev,
					"CAACAgIAAxkBAAMtY10MAAFrJdgnJMdebJv-51bnCvJkAAKgAAOrV8QLacLhq_lMpXsqBA",
				)
			case "придумай число до":
				num, err := strconv.Atoi(strings.Split(splitedTxt[1], ": ")[1])
				if err != nil {
					panic(err)
				}
				randNum := strconv.Itoa(rand.Intn(num))
				return SendMsg(lastId, ev, randNum)
			case "теперь ты":
				newName := strings.Split(splitedTxt[1], ": ")[1]

				if newName != "" {
					db, err := sql.Open("sqlite3", "file:database.db")
					if err != nil {
						panic(err)
					}
					defer db.Close()

					str := fmt.Sprintf("UPDATE bot_status SET name = '%s';", newName)
					db.Exec(str)

					return SendMsg(lastId, ev, fmt.Sprintln("Ну получается я", newName))
				}
			}
		}

		// return SendMsg(lastId, ev, "Не понимаю((")
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

	_, err := http.Post(tgApiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
	if err != nil {
		fmt.Println(err)
		return lastId
	}
	return ev.Id + 1
}

func SendStick(lastId int, ev UpdateStruct, sticker string) int {
	txtmsg := SendSticker{
		Chat_Id: ev.Message.Chat.Id,
		Sticker: sticker,
	}

	bytemsg, _ := json.Marshal(txtmsg)

	_, err := http.Post(tgApiUrl+"/sendSticker", "application/json", bytes.NewReader(bytemsg))
	if err != nil {
		fmt.Println(err)
		return lastId
	}
	return ev.Id + 1
}

func Ping() {
	txtmsg := SendMessage{
		Chat_Id: 911850117,
		Text:    "Страницу посетили.",
	}

	bytemsg, _ := json.Marshal(txtmsg)
	_, err := http.Post(tgApiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
	if err != nil {
		fmt.Println(err)
	}
}

func GetMyMainName() string {
	resp, err := http.Get(tgApiUrl + "/getMe")
	if err != nil {
		panic(err)
	}

	respBody, _ := io.ReadAll(resp.Body)

	var res MainStruct
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		panic(err)
	}

	return res.Result.Username
}
