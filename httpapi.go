package iotalerter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func httpAPI(w http.ResponseWriter, r *http.Request) {
	if !Config.Global.Communication.HTTP.Enabled {
		log.Println("alert: HTTPAPI Not Enabled! Function Called in Error!")
		return
	}

	SystemCommandAPIDefined := false
	UserCommandAPIDefined := false

	funcs := map[string]interface{}{
		"showmenu":    iotalerterMenu,
		"showbanner":  iotRelayBoardBannerAddColor,
		"showversion": iotRelayBoardShowVersion,
		"showconfig":  iotRelayBoardDumpConfig,
		"listapi":     listAPI}

	if !r.URL.Query().Has("command") {
		log.Println("error: URL Param 'command' is missing example http API commands should be of the format http://a.b.c.d/?command=listapi")
		fmt.Fprintf(w, "error: API should be of the format http://a.b.c.d:"+Config.Global.Communication.HTTP.Settings.Listenport+"/?command=listapi\n")
		return
	}

	APICommand := r.URL.Query().Get("command")

	for _, apicommand := range Config.Global.Communication.HTTP.Commands.Command {
		if APICommand == "listapi" && apicommand.Enabled {
			fmt.Fprintf(w, "200 OK: API Command %v Control Available\n", apicommand.Request)
			SystemCommandAPIDefined = true
			UserCommandAPIDefined = false
		} else if apicommand.Request == APICommand {
			SystemCommandAPIDefined = true
			UserCommandAPIDefined = false
		}
	}

	for _, apicommand := range Config.Global.Payloads.Payload {
		if apicommand.Name == APICommand {
			SystemCommandAPIDefined = false
			UserCommandAPIDefined = true
			break
		}
	}

	if SystemCommandAPIDefined {
		for _, apicommand := range Config.Global.Communication.HTTP.Commands.Command {
			if apicommand.Request == APICommand {
				_, err := Call(funcs, apicommand.Request)
				if err != nil {
					log.Println("error: Wrong Parameters to Call Function")
				}
			}
		}
		return
	}

	if UserCommandAPIDefined {
		for _, payloads := range Config.Global.Payloads.Payload {
			for _, payload := range payloads.Action {
				if payloads.Name == APICommand {
					if payload.Type == "gpio" && payload.Enabled {
						log.Printf("HTTPAPI Command GPIO item=%v param=%v", payload.Item, payload.Param)
						GPIOOutPinByItem(payload.Item, payload.Param)
					}
					if payload.Type == "http" && payload.Enabled {
						if payload.Method == "get" {
							HTTPAPIRequestGet(payload.Token, payload.URL)
						}
						if payload.Method == "post" {
							HTTPAPIRequestPost(payload.Token, payload.URL)
						}
					}
					if payload.Type == "mqtt" && payload.Enabled {
						log.Printf("payload=%v", payload.Payload)
					}
				}
			}
		}

		if !SystemCommandAPIDefined && !UserCommandAPIDefined {
			fmt.Fprintf(w, "Error: 404 API Not Defined\n")
			return
		}
	}
}

func listAPI() {
	for _, apicommand := range Config.Global.Communication.HTTP.Commands.Command {
		log.Printf("info: API Command %v Control Available\n", apicommand.Request)
	}
}

func HTTPAPIRequestGet(token string, urltocall string) {
	if !Config.Global.Communication.HTTP.Settings.Httpapirequest.Enabled {
		log.Println("error: HTTPAPIRequest Not Enabled in Config")
		return
	}
	if len(urltocall) == 0 {
		log.Println("error: HTTPAPIRequest URL Empty")
		return
	}

	response, err := http.Get(urltocall)

	if err != nil {
		log.Println("error: HTTPAPIRequest Error ", err.Error())
		return
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(responseData))
}

func HTTPAPIRequestPost(token string, urltocall string) {
	type Post struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	params := url.Values{}
	params.Add("title", "foo")
	params.Add("body", "bar")
	params.Add("userId", "1")

	resp, err := http.PostForm(urltocall, params)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	bodyString := string(body)
	log.Print(bodyString)
	post := Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	log.Printf("Post added with ID %d", post.ID)
}
