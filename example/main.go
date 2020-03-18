package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/tencentcaptcha"
	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

var (
	secretID            = ""         // 安全凭证 ID
	secretKey           = ""         // 安全凭证 Key
	appID        uint64 = 2004372521 // 验证码应用 ID
	appSecretKey        = ""         // 验证码应用 Secret Key
	app          *tencentcaptcha.Application
)

func main() {
	credential := common.NewCredential(secretID, secretKey)
	client, err := captcha.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		log.Fatal(err)
	}
	app = tencentcaptcha.New(client, appID, appSecretKey)

	router := clevergo.NewRouter()
	router.Get("/", index)
	router.Post("/verify", verify)
	log.Println(http.ListenAndServe(":8080", router))
}

var tmpl = template.Must(template.New("index").Parse(`
<html>
<head>
	<title>腾讯验证码 Tencent Captcha</title>
</head>
<body>
<button id="TencentCaptcha"
     data-appid="{{ .appID }}"
     data-cbfn="callback"
     type="button"
>验证</button>
<script src="https://code.jquery.com/jquery-3.4.1.min.js"></script>
<script src="https://ssl.captcha.qq.com/TCaptcha.js"></script>
<script>
window.callback = function(res){
	console.log(res)
	// res（用户主动关闭验证码）= {ret: 2, ticket: null}
	// res（验证成功） = {ret: 0, ticket: "String", randstr: "String"}
	if(res.ret === 0){
		$.post('/verify', {ticket: res.ticket, randstr: res.randstr}, function(resp) {
			alert(resp)
		})
	}
}
</script>
</body>
</html>
`))

func index(ctx *clevergo.Context) error {
	return tmpl.Execute(ctx.Response, map[string]interface{}{
		"appID": app.ID(),
	})
}

func verify(ctx *clevergo.Context) error {
	ticket := ctx.Request.PostFormValue("ticket")
	randstr := ctx.Request.PostFormValue("randstr")
	ipAddr := "127.0.0.1"
	if err := app.Verify(ticket, randstr, ipAddr); err != nil {
		ctx.WriteString(err.Error())
	} else {
		ctx.WriteString("success")
	}
	return nil
}
