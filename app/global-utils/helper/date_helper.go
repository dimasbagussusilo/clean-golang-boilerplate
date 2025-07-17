package helper

import (
	"appsku-golang/app/global-utils/constants"
	"fmt"
	"time"
)

// ToUTCfromGMT7 ...
func ToUTCfromGMT7(strTime string) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now(), err
	}

	date, err := time.ParseInLocation(constants.DATE_TIME_FORMAT_COMON, strTime, location)
	if err != nil {
		fmt.Printf("\nerror when parse strTime [%s] -> err: %v\n", strTime, err)
		return time.Now(), err
	}

	return date.In(time.UTC), nil
}

// FromUTCLocationToGMT7 ...
func FromUTCLocationToGMT7(date time.Time) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now(), err
	}

	return date.In(location), nil
}

// FromGMT7LocationUTCMin7 ...
func FromGMT7LocationUTCMin7(date time.Time) (time.Time, error) {
	date = date.Add(time.Hour * -7)
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now(), err
	}

	date = date.In(location)

	return date.In(time.UTC), nil
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// GetTimeLocationWIB get WIB location
func GetTimeLocationWIB() *time.Location {
	wib, _ := time.LoadLocation("Asia/Jakarta")
	return wib
}

func SetTimeZoneToWIB() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	time.Local = loc
}

// ConvertWibTimeToUtcWithTimeStartOrEnd isEndOfDay while value is 'true', time will be 23:59:59, while 'false' time will be 00:00:00
func ConvertWibTimeToUtcWithTimeStartOrEnd(dateStr string, isEndOfDay bool) string {
	if dateStr != "" {
		date, err := time.Parse(constants.DATE_FORMAT_COMMON, dateStr)
		if err != nil {
			return dateStr
		}
		if !isEndOfDay {
			startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, GetTimeLocationWIB())
			return startDate.UTC().Format("2006-01-02T15:04:05.999999999Z")
		} else {
			endDate := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, GetTimeLocationWIB())
			return endDate.UTC().Format("2006-01-02T15:04:05.999999999Z")
		}
	}
	return ""
}

func ValidateDate(date string, format string) bool {
	// Parse date using the provided format
	parsedDate, err := time.Parse(format, date)
	if err != nil {
		// Return false if parsing fails
		return false
	}

	// Reformat the parsed date to match the provided format and compare
	return parsedDate.Format(format) == date
}

func ValidateDateBeforeToday(date string, format string) (bool, error) {
	// Parse tanggal input berdasarkan format
	inputDate, err := time.Parse(format, date)
	if err != nil {
		return false, fmt.Errorf("format tanggal tidak valid: %v", err)
	}

	// Ambil tanggal hari ini tanpa waktu (hanya Y-M-D)
	today := time.Now().Truncate(24 * time.Hour)

	// Periksa apakah tanggal input tidak melebihi hari ini
	if inputDate.After(today) {
		return false, nil
	}

	return true, nil
}

func ValidateDateBeforeOrAfter(date, dateCompate, tipe, format string) (bool, error) {
	// Parse tanggal input berdasarkan format
	date1, err := time.Parse(format, date)
	if err != nil {
		return false, fmt.Errorf("format tanggal tidak valid: %v", err)
	}
	date2, err := time.Parse(format, dateCompate)
	if err != nil {
		return false, fmt.Errorf("format tanggal tidak valid: %v", err)
	}

	// Periksa apakah tanggal input tidak melebihi hari ini
	if tipe == "after" {
		if date1.After(date2) {
			return false, nil
		}
	} else if tipe == "before" {
		if date1.Before(date2) {
			return false, nil
		}
	}

	return true, nil
}

func StringToDateTime(inputFormat, outputFormat, dateStr string) (time.Time, error) {

	// Parse string ke tipe time
	parsedDate, err := time.Parse(inputFormat, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{}, err
	}

	// Format output sesuai dengan format yang diinginkan
	formattedDateStr := parsedDate.Format(outputFormat)
	formattedDate, err := time.Parse(outputFormat, formattedDateStr)
	if err != nil {
		fmt.Println("Error re-parsing formatted date:", err)
		return time.Time{}, err
	}

	// Cetak hasil
	return formattedDate, err
}
