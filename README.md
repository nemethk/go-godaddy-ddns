# go-godaddy-ddns
## USAGE
### Varriables
| Tag   | Type | Details | Default Value |
|:--------:|:--------:|:--------:|:--------:|
| `-godaddy-key` | required | Godaddy API key. |  |
| `-godaddy-secret` | required | Godaddy API secret. |  |
| `-godaddy-domain` | required | The registered domain name. |  |
| `-slack-token` | required | Slack Bearer token. |  |
| `-channel-id` | required | Slack channel ID. |  |
| `-ip-provider` | optional | IP provider API. | `https://v4.ident.me/` |
### Run
```
go run main.go \
    -godaddy-key $GKEY \
    -godaddy-secret $GSECRET \
    -godaddy-domain $GDOMAIN \
    -slack-token $STOKEN \
    -channel-id $SCHANNEL
```