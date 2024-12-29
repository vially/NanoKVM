package config

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	gpioPowerFileEnvVar    = os.Getenv("GPIO_POWER_FILE")
	gpioPowerLEDFileEnvVar = os.Getenv("GPIO_POWER_LED_FILE")
	gpioResetFileEnvVar    = os.Getenv("GPIO_RESET_FILE")
	gpioHDDLedFileEnvVar   = os.Getenv("GPIO_HDD_LED_FILE")
)

type HWVersion int

const (
	HWVersionAlpha HWVersion = iota
	HWVersionBeta
	HWVersionPcie

	HWVersionFile = "/etc/kvm/hw"
)

var HWAlpha = Hardware{
	Version:      HWVersionAlpha,
	GPIOReset:    "/sys/class/gpio/gpio507/value",
	GPIOPower:    "/sys/class/gpio/gpio503/value",
	GPIOPowerLED: "/sys/class/gpio/gpio504/value",
	GPIOHDDLed:   "/sys/class/gpio/gpio505/value",
}

var HWBeta = Hardware{
	Version:      HWVersionBeta,
	GPIOReset:    "/sys/class/gpio/gpio505/value",
	GPIOPower:    "/sys/class/gpio/gpio503/value",
	GPIOPowerLED: "/sys/class/gpio/gpio504/value",
	GPIOHDDLed:   "",
}

var HWPcie = Hardware{
	Version:      HWVersionPcie,
	GPIOReset:    "/sys/class/gpio/gpio505/value",
	GPIOPower:    "/sys/class/gpio/gpio503/value",
	GPIOPowerLED: "/sys/class/gpio/gpio504/value",
	GPIOHDDLed:   "",
}

func (h HWVersion) String() string {
	switch h {
	case HWVersionAlpha:
		return "Alpha"
	case HWVersionBeta:
		return "Beta"
	case HWVersionPcie:
		return "PCIE"
	default:
		return "Unknown"
	}
}

func getHwVersion() HWVersion {
	content, err := os.ReadFile(HWVersionFile)
	if err != nil {
		return HWVersionAlpha
	}

	version := strings.ReplaceAll(string(content), "\n", "")
	if version == "beta" {
		return HWVersionBeta
	} else if version == "pcie" {
		return HWVersionPcie
	}

	return HWVersionAlpha
}

func getHardware() (h Hardware) {
	version := getHwVersion()

	switch version {
	case HWVersionAlpha:
		h = HWAlpha

	case HWVersionBeta:
		h = HWBeta

	case HWVersionPcie:
		h = HWPcie

	default:
		h = HWAlpha
		log.Error("Unsupported hardware version: %s", version)
	}

	if gpioPowerFileEnvVar != "" {
		h.GPIOPower = gpioPowerFileEnvVar
	}

	if gpioPowerLEDFileEnvVar != "" {
		h.GPIOPowerLED = gpioPowerLEDFileEnvVar
	}

	if gpioResetFileEnvVar != "" {
		h.GPIOReset = gpioResetFileEnvVar
	}

	if gpioHDDLedFileEnvVar != "" {
		h.GPIOHDDLed = gpioHDDLedFileEnvVar
	}

	return
}
