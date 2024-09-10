# pkg_login
第三方 OAuth2 授权登录，参考https://github.com/netnr/Netnr.Login，使用go语言实现
### 安装
```
go get github.com/juxiaoming/pkg_login
```

### 支持第三方登录
<table>
    <tr><th>三方</th><th>参考文档</th><th>应用申请（已登录）</th></tr>
    <tr>
        <td><img src="https://gs.zme.ink/static/login/dingtalk.svg" height="30" title="钉钉/DingTalk"></td>
        <td><a target="_blank" href="https://open.dingtalk.com/document/tutorial/scan-qr-code-to-log-on-to-third-party-websites">参考文档</a></td>
        <td><a target="_blank" href="https://open-dev.dingtalk.com/#/loginMan">应用申请</a></td>
    </tr>
    <tr>
        <td><img src="https://gs.zme.ink/static/login/feishu.svg" height="30" title="飞书/FeiShu"></td>
        <td><a target="_blank" href="https://open.feishu.cn/document/common-capabilities/sso/web-application-sso/web-app-overview">参考文档</a></td>
        <td><a target="_blank" href="https://open.feishu.cn/app">应用申请</a></td>
    </tr>
    <tr>
        <td><img src="https://gs.zme.ink/static/login/gitee.svg" height="30" title="码云/Gitee"></td>
        <td><a target="_blank" href="https://gitee.com/api/v5/oauth_doc">参考文档</a></td>
        <td><a target="_blank" href="https://gitee.com/oauth/applications">应用申请</a></td>
    </tr>
    <tr>
        <td><img src="https://gs.zme.ink/static/login/github.svg" height="30" title="GitHub"></td>
        <td><a target="_blank" href="https://docs.github.com/en/developers/apps/building-oauth-apps/authorizing-oauth-apps">参考文档</a></td>
        <td><a target="_blank" href="https://github.com/settings/developers">应用申请</a></td>
    </tr>
    <tr>
        <td><img src="https://gs.zme.ink/static/login/google.svg" height="30" title="谷歌/Google"></td>
        <td><a target="_blank" href="https://developers.google.com/identity/protocols/oauth2/web-server">参考文档</a></td>
        <td><a target="_blank" href="https://console.developers.google.com/apis/credentials">应用申请</a></td>
    </tr>
</table>

### 更多
由于账号原因微信、qq、微博、支付宝还没有封装，等我！

### 使用
```go
        pkg_login.Init(pkg_login.NewFeiShuConf("your_id", "your_secret", "redirect_url"))
	server, err := pkg_login.NewServer(pkg_login.ImplementFeiShu)
	if err != nil {
		fmt.Println("初始化错误")
		return
	}
	fmt.Println(server.RedirectUrl())

	fmt.Println(server.GetUserinfo("your_code"))
```