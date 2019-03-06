package utils

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	. "git.tor.ph/hiveon/pool/models"

)

func GetUserByEmail(email string) (OAuthUser, error) {
	url := fmt.Sprintf(GetConfig().GetString("pool.idp_api")+"/users/email/%s", email)
	var user OAuthUser
	res, err := resty.R().Get(url)
	if err != nil {
		return user, err
	}
	json.Unmarshal(res.Body(), &user)
	return user, nil
}
