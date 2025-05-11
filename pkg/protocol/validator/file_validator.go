package validator

import (
	"errors"
	"net"
	"os"
)

func Validate(ip string, filePath string) error {
	if err := validateFilePath(filePath); err != nil {
		return errors.New("file path is invalid")
	}

	if err := validateIp(ip); err != nil {
		return errors.New("ip is invalid")
	}

	return nil
}

func validateIp(ip string) error {
	if ip == "" {
		return errors.New("ip is empty")
	}

	if net.ParseIP(ip) == nil {
		return errors.New("invalid ip address")
	}

	return nil
}

func validateFilePath(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	return nil
}
