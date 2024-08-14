package http

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	bson2 "github.com/z26100/go-utils/bson"
	errors2 "github.com/z26100/go-utils/errors"
	"github.com/z26100/log-go"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type HttpClient interface {
	Get(uri string, resp interface{}) error
	Post(uri string, requestBody interface{}, response interface{}) error
	Delete(uri string, response interface{}) error
}

func TlsInsecureVerify() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func IsCustomError(err error) bool {
	return reflect.TypeOf(err).Name() == reflect.TypeOf(errors2.CustomError{}).Name()
}

func AsCustomError(err error) errors2.CustomError {
	return err.(errors2.CustomError)
}

var (
	Pretty = true
)

func RespondWithJson(context *gin.Context, resp interface{}, err error) {
	if CheckError(context, err, http.StatusBadRequest) {
		return
	}
	var data []byte
	switch Pretty {
	case true:
		data, err = json.MarshalIndent(resp, "", "  ")
	default:
		data, err = json.Marshal(resp)
	}
	fmt.Println(string(data))
	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	_, err = context.Writer.Write(data)
	if err != nil {
		log.Error(err)
	}
}

func GetRequestJWT(auth string) (string, error) {
	if auth == "" {
		return auth, errors.New("no authorization found")
	}
	tokens := strings.Split(auth, " ")
	if len(tokens) != 2 || tokens[0] != "Bearer" {
		return "", errors.New("No Bearer Token found")
	}
	return tokens[1], nil
}

func SendResults(ctx *gin.Context, data interface{}, count ...int64) {
	if data == nil {
		data = []bson.M{}
	}

	var ct int64
	if len(count) == 1 {
		ct = count[0]
	}
	var jsonData []byte
	var err error

	switch ctx.Query("pretty") == "true" {
	case true:
		jsonData, err = bson.MarshalExtJSONIndent(bson.M{"header": createHeader(ct), "body": data}, bson2.Canonical, bson2.EscapeHtml, "", " ")
	default:
		jsonData, err = bson.MarshalExtJSON(bson.M{"header": createHeader(ct), "body": data}, bson2.Canonical, bson2.EscapeHtml)
	}
	if CheckError(ctx, err, http.StatusBadRequest) {
		return
	}
	_, err = ctx.Writer.Write(jsonData)
}

type Header struct {
	Id        string `json:"_id"`
	Timestamp string `json:"timestamp"`
	Count     int64  `json:"count,omitempty"`
}

func createHeader(count int64) Header {
	uuid, _ := uuid.GenerateUUID()
	return Header{
		uuid,
		time.Now().UTC().Format(time.RFC3339),
		count,
	}
}

func CheckError(ctx *gin.Context, err error, code ...int) bool {
	if err != nil {
		log.Error(err)
		if len(code) > 0 {
			ctx.AbortWithStatus(code[0])
		} else {
			ctx.Abort()
		}
	}
	return err != nil
}
