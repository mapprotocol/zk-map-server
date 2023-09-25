package resp

const (
	CodeSuccess             = 2000
	CodeParameterErr        = 3000
	CodeExternalServerError = 4000
	CodeInternalServerError = 5000
)

const (
	MsgSuccess             = "Success"
	MsgParameterErr        = "Invalid Parameter"
	MsgExternalServerError = "External Server Error"
	MsgInternalServerError = "Internal Server Error"
)

var code2msg = map[int64]string{
	CodeSuccess:             MsgSuccess,
	CodeParameterErr:        MsgParameterErr,
	CodeExternalServerError: MsgExternalServerError,
	CodeInternalServerError: MsgInternalServerError,
}
