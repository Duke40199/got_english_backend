package utils

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var fullnameRegex = regexp.MustCompile("([a-zA-Z',.-]+( [a-zA-Z',.-]+)*){2,30}")
var phoneNumberRegex = regexp.MustCompile(`\+?(84|03|05|07|08|09|01[2|6|8|9])+([0-9]{8,10})\b`)

// isEmailValid checks if the email provided passes the required structure and length.
func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 55 {
		return false
	}
	return emailRegex.MatchString(e)
}

// isEmailValid checks if the email provided passes the required structure and length.
func IsPhoneNumberValid(e string) (bool, error) {
	if e != "" {
		if len(e) < 3 || len(e) > 15 {
			return false, errors.New("invalid phone number length (3 to 15 characters)")
		}
		if !phoneNumberRegex.MatchString(e) {
			return false, errors.New("incorrect phone number format")
		}
		return true, nil
	} else {
		return true, nil
	}
}

// isEmailValid checks if the email provided passes the required structure and length.
func IsFullnameValid(e string) (bool, error) {
	if e != "" {
		if len(e) < 3 || len(e) > 55 {
			return false, errors.New("invalid fullname length (3 to 55 characters)")
		}
		if !fullnameRegex.MatchString(e) {
			return false, errors.New("incorrect fullname format")
		}
		return true, nil
	} else {
		return true, nil
	}
}

// isEmailValid checks if the email provided passes the required structure and length.
func IsAddressValid(e string) (bool, error) {
	if e != "" {
		if len(e) < 3 || len(e) > 100 {
			return false, errors.New("invalid address length (3 to 100 characters)")
		}
		return true, nil
	} else {
		return true, nil
	}
}

// isEmailValid checks if the email provided passes the required structure and length.
func IsUsernameValid(e string) (bool, error) {
	if e != "" {
		if len(e) < 3 || len(e) > 55 {
			return false, errors.New("invalid username length (3 to 55 characters)")
		}
		return true, nil
	} else {
		return false, errors.New(`username contains ""`)
	}
}

// isEmailValid checks if the email provided passes the required structure and length.
func IsBirthdayValid(e string) (bool, error) {
	if e != "" {
		t, err := time.Parse("2006-01-02", e)
		if err != nil {
			return false, err
		}
		timeNow := time.Now()
		if t.Year() > timeNow.Year() {
			return false, errors.New("birthday cannot set in the future")
		} else if t.Year() == timeNow.Year() {
			if t.Month() > timeNow.Month() {
				return false, errors.New("birthday cannot set in the future")
			} else if t.Day() > timeNow.Day() {
				return false, errors.New("birthday cannot set in the future")
			}
		}
		return true, nil
	} else {
		return false, errors.New(`birthday contains ""`)
	}
}

func CalculateExpertEarningBySession(exchangeRateValue float32, coinValueInVND uint, coinCount uint) float32 {
	result := (1 - exchangeRateValue) * float32(coinValueInVND*coinCount)
	return result
}

// GetCurrentTime function is used to get the current time in milliseconds.
func GetCurrentEpochTimeInMiliseconds() int64 {
	var now = time.Now()
	ts := now.UnixNano() / 1000000
	return ts
}

func GetTimesByPeriod(period string) (time.Time, time.Time, error) {
	var startDate time.Time
	var endDate time.Time
	var timeNow = time.Now()
	switch period {
	case "weekly":
		{
			endDate = time.Now()
			year, week := endDate.ISOWeek()
			//Get week start date
			startDate = GetWeekStart(year, week)
			break
		}
	case "daily":
		{
			startDate = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.Local)
			endDate = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 23, 59, 59, 0, time.Local)
			break
		}
	case "monthly":
		{
			fmt.Println("================= monthly filter")
			startDate = time.Date(timeNow.Year(), timeNow.Month(), 1, 0, 0, 0, 0, time.Local)
			endDate = time.Date(timeNow.Year(), timeNow.Month()+1, 0, 23, 59, 59, 0, time.Local)
			break
		}
	case "yearly":
		{
			fmt.Println("========= Yearly filter")
			startDate = time.Date(timeNow.Year(), 1, 1, 00, 00, 00, 00, time.Local)
			endDate = time.Date(timeNow.Year(), 12, 31, 23, 59, 59, 00, time.Local)
			break
		}
	default:
		{
			fmt.Println("============ INVALID")
			return startDate, endDate, errors.New("invalid time period")
		}
	}
	return startDate, endDate, nil
}

func GetWeekStart(year, week int) time.Time {

	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if weekday := t.Weekday(); weekday == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(weekday)+1)
	}
	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

func CalculateRemainingCoin() int64 {
	var now = time.Now()
	ts := now.UnixNano() / 1000000
	return ts
}
