package fileWork

import (
	"fmt"
	"os"
	"strings"
)

// CheckInput проверяет ввод с консоли
func CheckInput(startDirectory string, sortType string) error {
	if _, err := os.Stat(startDirectory); os.IsNotExist(err) {
		return fmt.Errorf("directory by path : %s is not exist", startDirectory)
	}
	if strings.ToUpper(sortType) != ascType && strings.ToUpper(sortType) != deskType {
		return fmt.Errorf("sort type can be ASC or DESK! %s", sortType)
	}

	return nil
}
