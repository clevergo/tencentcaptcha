// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// in the LICENSE file.

package tencentcaptcha

import (
	"errors"

	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
)

// Application contains captcha's infomation, such as application ID, secret key.
type Application struct {
	client      *captcha.Client
	id          uint64
	secretKey   string
	captchaType uint64
}

// New returns a captcha application with the given client, application ID and secret key.
func New(client *captcha.Client, id uint64, secretKey string) *Application {
	return &Application{
		client:      client,
		id:          id,
		secretKey:   secretKey,
		captchaType: 9,
	}
}

// ID returns the application ID.
func (app *Application) ID() uint64 {
	return app.id
}

// Verify verifies whether is captcha info is valid.
func (app *Application) Verify(ticket, randstr, ipAddr string) error {
	req := captcha.NewDescribeCaptchaResultRequest()
	req.Ticket = &ticket
	req.Randstr = &randstr
	req.AppSecretKey = &app.secretKey
	req.CaptchaAppId = &app.id
	req.CaptchaType = &app.captchaType
	req.UserIp = &ipAddr
	resp, err := app.client.DescribeCaptchaResult(req)
	if err != nil {
		return err
	}
	if *resp.Response.CaptchaCode != 1 {
		return errors.New(*resp.Response.CaptchaMsg)
	}

	return nil
}
