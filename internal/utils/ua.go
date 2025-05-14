package utils

import (
	"strings"

	"github.com/mssola/useragent"
)

type UADeviceInfo struct {
	Device string
	OS     string
}

func ParseUserAgent(userAgent string) UADeviceInfo {
	ua := useragent.New(userAgent)

	osName := ua.OSInfo().Name
	osVersion := ua.OSInfo().Version
	os := osName
	if osVersion != "" {
		os = osName + " " + osVersion
	}

	device := "Unknown Device"
	if ua.Mobile() {
		deviceModel := ua.Model()
		if deviceModel != "" {
			device = "Mobile " + deviceModel
		} else {
			device += " " + ua.Platform()
		}
	} else {
		deviceModel := ua.Model()
		if deviceModel != "" {
			device = "Desktop " + deviceModel
		} else {
			device += " " + ua.Platform()
		}
	}

	// Special case for iPhone
	if strings.Contains(strings.ToLower(userAgent), "iphone") {
		device = "iPhone"
		if ua.Model() != "" {
			device += " " + ua.Model()
		}
	}

	return UADeviceInfo{
		Device: device,
		OS:     os,
	}
}
