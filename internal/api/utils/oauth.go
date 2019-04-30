package utils

import (
	"encoding/json"
	"fmt"
	"git.tor.ph/hiveon/pool/models"
	"gopkg.in/resty.v1"
)

func GetUserByEmail(email string) (models.OAuthUser, error) {
	url := fmt.Sprintf(GetConfig().GetString("pool.idp_api")+"/users/email/%s", email)
	var user models.OAuthUser
	res, err := resty.R().Get(url)
	if err != nil {
		return user, err
	}
	json.Unmarshal(res.Body(), &user)
	return user, nil
}
