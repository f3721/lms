package utils

import (
	"fmt"
	"time"
)

func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// StrToTime 字符串时间转time类型
func StrToTime(str string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	if len(str) == 10 {
		layout = "2006-01-02"
	}
	return time.Parse(layout, str)
}

// startDateStr、endDateStr "202305"
func GetMonthsBetweenDates(startDateStr, endDateStr string) []string {
	var months []string
	startDate, _ := time.Parse("2006-01", startDateStr)
	endDate, _ := time.Parse("2006-01", endDateStr)
	current := startDate
	for current.Before(endDate) {
		months = append(months, current.Format("200601"))
		// Add one month to the current date
		current = current.AddDate(0, 1, 0)
	}
	months = append(months, endDate.Format("200601"))
	return months
}

// 获取两个时间的差值 t1 - t2 格式为 x年x月x天x小时
func GetDiffTextByTime(t1, t2 time.Time) (text string) {
	if t2.Before(t1) {
		duration := t1.Sub(t2)
		year := int(duration.Hours() / 24 / 365)
		if year > 0 {
			text = text + fmt.Sprintf("%d年", year)
		}
		month := int(duration.Hours()/24) % 365 / 30
		if month > 0 {
			text = text + fmt.Sprintf("%d月", month)
		}
		day := int(duration.Hours()/24) % 365 % 30
		if day > 0 {
			text = text + fmt.Sprintf("%d天", day)
		}
		hour := int(duration.Hours()) % 24
		if hour > 0 {
			text = text + fmt.Sprintf("%d小时", hour)
		}
		//minute := int(duration.Minutes())%60
		//if minute > 0 {
		//	text = text + fmt.Sprintf("%分钟", minute)
		//}
		//seconds := int(duration.Seconds())%60
		//if seconds > 0 {
		//	text = text + fmt.Sprintf("%秒", seconds)
		//}
	}

	return
}
