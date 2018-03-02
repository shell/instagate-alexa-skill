package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ahmdrz/goinsta"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/yanatan16/golang-instagram/instagram"
)

// AlexaSkillEvent is a struct for the Alexa Skill Event
type AlexaSkillEvent struct {
	Session struct {
		New         bool   `json:"new"`
		SessionID   string `json:"sessionId"`
		Application struct {
			ApplicationID string `json:"applicationId"`
		} `json:"application"`
		Attributes interface{} `json:"attributes"`
		User       struct {
			Userid string `json:"userId"`
		} `json:"user"`
	} `json:"session"`
	Request struct {
		Type      string `json:"type"`
		RequestID string `json:"requestId"`
		Intent    struct {
			Name  string `json:"name"`
			Slots struct {
				Friend struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"Friend"`
			} `json:"slots"`
		} `json:"intent"`
		Locale    string `json:"locale"`
		Timestamp string `json:"timestamp"`
	} `json:"request"`
	Version string `json:"version"`
}

// AlexaResponse is a struct for the resonse to Alexa Skill Kit
type AlexaResponse struct {
	Version  string `json:"version,omitempty"`
	Response struct {
		OutputSpeech struct {
			Type string `json:"type"`
			Text string `json:"text"`
			SSML string `json:"ssml"`
		} `json:"outputSpeech"`
	} `json:"response"`
	ShouldEndSession  bool     `json:"shouldEndSession"`
	SessionAttributes struct{} `json:"sessionAttributes"`
}

// FriendsAliases are aliases for a friend slots
var FriendsAliases = map[string]string{
	"marina":  "Marina",
	"maya":    "Marina",
	"marisha": "Marina",
	"alexey":  "Alexey",
	"lyoha":   "Alexey",
	"dasha":   "Dasha",
	"daria":   "Dasha",
}

// FriendsMap is a map of my friends names and their usernames in instagram
var FriendsMap = map[string]string{
	"Marina": "maia_mois",
	"Alexey": "bigxam",
	"Dasha":  "sashchenko_",
}

// GetInstaHandler returns handler for a instagram fetching func
func GetInstaHandler(instaClient *goinsta.Instagram) func(request *AlexaSkillEvent) (*AlexaResponse, error) {
	return func(request *AlexaSkillEvent) (*AlexaResponse, error) {
		// stdout and stderr are sent to AWS CloudWatch Logs
		fmt.Println("Processing Lambda request for ", request.Request.Intent.Name)

		return processRequest(request, instaClient), nil
	}
}

func processRequest(event *AlexaSkillEvent, instaClient *goinsta.Instagram) *AlexaResponse {
	text := ""
	respType := "PlainText"

	switch event.Request.Intent.Name {
	case "MyFollowersCount":
		resp, err := instaClient.GetUserByUsername("penkinv")
		if err != nil {
			text = fmt.Sprintf("Can't reach instagram servers, error: %s", err)
		} else {
			text = fmt.Sprintf("You have %d followers", resp.User.FollowerCount)
		}
	case "FriendFollowersCount":
		friend := strings.ToLower(event.Request.Intent.Slots.Friend.Value)

		resp, err := instaClient.GetUserByUsername(FriendsMap[FriendsAliases[friend]])
		if err != nil {
			text = fmt.Sprintf("Can't reach instagram servers, error: %s", err)
		} else {
			text = fmt.Sprintf("%s has %d followers", friend, resp.User.FollowerCount)
		}
	case "MyLikes":
		// Special case: Using different instagram API to fetch my most recent media
		instaAPI := instagram.New(os.Getenv("InstaClientId"), os.Getenv("InstaClientSecret"), os.Getenv("InstaAccessToken"), false)
		if ok, err := instaAPI.VerifyCredentials(); !ok {
			text = fmt.Sprintf("Can't reach instagram servers, error: %s", err)
		}

		params := url.Values{}
		params.Set("count", "1")

		resp, err := instaAPI.GetUserRecentMedia("self", params)
		if err != nil {
			text = fmt.Sprintf("Can't reach instagram servers, error: %s", err)
		} else {
			text = fmt.Sprintf("You have %d likes on your most recent photo", resp.Medias[0].Likes.Count)
		}
	default:
		text = "I have no idea what you just said, could you speak more clearly?"
	}

	return generateAlexaResponse(text, respType)
}

func generateAlexaResponse(text string, respType string) *AlexaResponse {
	response := &AlexaResponse{
		Version: "1.0",
	}
	response.Response.OutputSpeech.Type = respType
	switch respType {
	case "PlainText":
		response.Response.OutputSpeech.Text = text
	case "SSML":
		response.Response.OutputSpeech.SSML = text
	}

	response.ShouldEndSession = true
	return response
}

func main() {
	instaClient := goinsta.New(os.Getenv("InstaLogin"), os.Getenv("InstaPassword"))
	if err := instaClient.Login(); err != nil {
		panic(err)
	}
	defer instaClient.Logout()

	fmt.Println("Starting lambdas")
	lambda.Start(GetInstaHandler(instaClient))
}
