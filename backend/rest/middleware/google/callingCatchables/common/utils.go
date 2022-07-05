package callingCatchables

import (
	"encoding/json"
	"net/http"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Route404Error(st *status.Status, rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(rw)
	er := encoder.Encode(struct {
		RpcError  codes.Code `json:"rpc_error"`
		HtmlError int        `json:"html_error"`
		Details   string     `json:"code"`
	}{HtmlError: rpcToHtmlError(st.Code()), RpcError: st.Code(),
		Details: st.Message()})

	if er != nil {
		panic(er)
	}
}

func rpcToHtmlError(code codes.Code) int {
	switch code {
	case 3:
		return 400
	default:
		return 404
	}
}

func Str2Int32(val string) (int32, error) {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return int32(i), nil
}

func Str2Bool(val string) (bool, error) {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}

	return b, nil
}
