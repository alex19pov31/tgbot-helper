package tgbothelper

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Arman92/go-tdlib"
)

var client *tdlib.Client
var allChats []*tdlib.Chat
var haveFullChatList bool

var accountName string
var configPath string

// ProxySocks5 - настройка прокси socks5
type ProxySocks5 struct {
	Server   string `json:"server"`
	Port     int32  `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ProxyHTTP - настройка http прокси
type ProxyHTTP struct {
	Server   string `json:"server"`
	Port     int32  `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	UseHTTP  bool   `json:"use_http"`
}

// ProxyMtproto - настройка mtproto прокси
type ProxyMtproto struct {
	Server string `json:"server"`
	Port   int32  `json:"port"`
	Secret string `json:"secret"`
}

// Client - клиент telegram user api
type Client struct {
	client       *tdlib.Client
	booted       bool
	APIID        string
	APIHash      string
	accountName  string
	proxySocks5  *ProxySocks5
	proxyHTTP    *ProxyHTTP
	proxyMtproto *ProxyMtproto
	allChats     []*tdlib.Chat
}

// SetSocks5Proxy - применить socks5 прокси
func (c *Client) SetSocks5Proxy(server string, port int32, username, password string) {
	c.proxySocks5 = &ProxySocks5{
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// SetHTTPProxy - применить http прокси
func (c *Client) SetHTTPProxy(server string, port int32, username, password string, useHTTP bool) {
	c.proxyHTTP = &ProxyHTTP{
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
		UseHTTP:  useHTTP,
	}
}

// SetMtprotoProxy - применить mtproxy прокси
func (c *Client) SetMtprotoProxy(server string, port int32, secret string) {
	c.proxyMtproto = &ProxyMtproto{
		Server: server,
		Port:   port,
		Secret: secret,
	}
}

func (c *Client) getClient() *tdlib.Client {
	if c.client != nil {
		return c.client
	}

	config := tdlib.Config{
		APIID:               c.APIID,
		APIHash:             c.APIHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  false,
		UseFileDatabase:     false,
		UseChatInfoDatabase: false,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./" + c.accountName + "/tdlib-db",
		FileDirectory:       "./" + c.accountName + "/tdlib-files",
		IgnoreFileNames:     false,
	}

	curClient := tdlib.NewClient(config)
	if c.proxySocks5 != nil && c.proxySocks5.Server != "" && c.proxySocks5.Port > 0 {
		curClient.AddProxy(c.proxySocks5.Server, c.proxySocks5.Port, true, tdlib.NewProxyTypeSocks5(c.proxySocks5.Username, c.proxySocks5.Password))
	}
	if c.proxyHTTP != nil && c.proxyHTTP.Server != "" && c.proxyHTTP.Port > 0 {
		curClient.AddProxy(c.proxyHTTP.Server, c.proxyHTTP.Port, true, tdlib.NewProxyTypeHttp(c.proxyHTTP.Username, c.proxyHTTP.Password, c.proxyHTTP.UseHTTP))
	}
	if c.proxyMtproto != nil && c.proxyMtproto.Server != "" && c.proxyMtproto.Port > 0 {
		curClient.AddProxy(c.proxyMtproto.Server, c.proxyMtproto.Port, true, tdlib.NewProxyTypeMtproto(c.proxyMtproto.Secret))
	}

	c.client = curClient
	return c.client
}

// GetClient - telegram клиент
func (c *Client) GetClient() *tdlib.Client {
	currentState, _ := c.getClient().Authorize()
	for ; currentState.GetAuthorizationStateEnum() != tdlib.AuthorizationStateReadyType; currentState, _ = c.getClient().Authorize() {
		switch currentState.GetAuthorizationStateEnum() {
		case tdlib.AuthorizationStateWaitPhoneNumberType:
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := c.getClient().SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		case tdlib.AuthorizationStateWaitCodeType:
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := c.getClient().SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		case tdlib.AuthorizationStateWaitPasswordType:
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := c.getClient().SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		}
		time.Sleep(300 * time.Millisecond)
	}

	c.getChatList(1000)

	return c.client
}

func (c *Client) getChatList(limit int) error {
	if c.allChats == nil {
		c.allChats = []*tdlib.Chat{}
	}

	if !haveFullChatList && limit > len(c.allChats) {
		offsetOrder := int64(math.MaxInt64)
		offsetChatID := int64(0)
		var lastChat *tdlib.Chat

		if len(c.allChats) > 0 {
			lastChat = c.allChats[len(c.allChats)-1]
			offsetOrder = int64(lastChat.Order)
			offsetChatID = lastChat.ID
		}

		// get chats (ids) from tdlib
		chats, err := c.getClient().GetChats(tdlib.JSONInt64(offsetOrder),
			offsetChatID, int32(limit-len(c.allChats)))
		if err != nil {
			return err
		}
		if len(chats.ChatIDs) == 0 {
			haveFullChatList = true
			return nil
		}

		for _, chatID := range chats.ChatIDs {
			// get chat info from tdlib
			chat, err := c.getClient().GetChat(chatID)
			if err == nil {
				c.allChats = append(c.allChats, chat)
			} else {
				return err
			}
		}
		return c.getChatList(limit)
	}
	return nil
}

// GetChatByName возвращает по названию
func (c *Client) GetChatByName(name string) (*tdlib.Chat, error) {
	for _, chat := range c.allChats {
		if name == chat.Title {
			return chat, nil
		}
	}

	return &tdlib.Chat{}, errors.New("Chat not found")
}
