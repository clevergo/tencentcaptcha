// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package tencentcaptcha

import (
	"errors"

	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
)

// Captcha contains captcha's infomation, such as application ID, secret key.
type Captcha struct {
	client      *captcha.Client
	appID       uint64
	secretKey   string
	captchaType uint64
}

// New returns a captcha with the given client, application ID and secret key.
func New(client *captcha.Client, appID uint64, secretKey string) *Captcha {
	return &Captcha{
		client:      client,
		appID:       appID,
		secretKey:   secretKey,
		captchaType: 9,
	}
}

// AppID returns the application ID.
func (c *Captcha) AppID() uint64 {
	return c.appID
}

// Verify verifies whether is captcha info is valid.
func (c *Captcha) Verify(ticket, randstr, ipAddr string) error {
	req := captcha.NewDescribeCaptchaResultRequest()
	req.Ticket = &ticket
	req.Randstr = &randstr
	req.AppSecretKey = &c.secretKey
	req.CaptchaAppId = &c.appID
	req.CaptchaType = &c.captchaType
	req.UserIp = &ipAddr
	resp, err := c.client.DescribeCaptchaResult(req)
	if err != nil {
		return err
	}
	if *resp.Response.CaptchaCode != 1 {
		return errors.New(*resp.Response.CaptchaMsg)
	}

	return nil
}
