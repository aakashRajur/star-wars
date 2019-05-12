package util

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func MapStringFromInterfaces(args []interface{}) ([]string, error) {
	compiled := make([]string, 0)
	for _, each := range args {
		switch t := each.(type) {
		case float64:
			{
				compiled = append(compiled, strconv.FormatFloat(t, 'f', 64, -1))
				break
			}
		case float32:
			{
				compiled = append(compiled, strconv.FormatFloat(float64(t), 'f', 64, -1))
				break
			}
		case int64:
			{
				compiled = append(compiled, strconv.FormatInt(t, 10))
				break
			}
		case int32:
		case int16:
		case int8:
		case int:
			{
				compiled = append(compiled, strconv.FormatInt(int64(t), 10))
				break
			}
		case bool:
			{
				compiled = append(compiled, strconv.FormatBool(t))
				break
			}
		case string:
			{
				compiled = append(compiled, t)
				break
			}
		case []byte:
			{
				compiled = append(compiled, string(t[:]))
			}
		default:
			return nil, errors.New(`type not supported`)
		}
	}
	return compiled, nil
}

func GetType(value interface{}) string {
	return reflect.TypeOf(value).Kind().String()
}

func LogRequest(logger types.Logger, request *http.Request) {
	params := request.Context().Value(http.PARAMS).(map[string]interface{})
	query := request.URL.Query()

	logger.InfoFields(
		map[string]interface{}{
			"METHOD":     request.Method,
			"URI":        request.RequestURI,
			"URL":        request.URL.Path,
			"PARAMS":     params,
			"QUERY":      query,
			"REMOTE":     request.RemoteAddr,
			"REFERER":    request.Referer(),
			"USER-AGENT": request.UserAgent(),
		},
	)
}

func DurationToString(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	var parts []string

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}
