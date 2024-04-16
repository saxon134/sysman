package toolbox

// GetYearDays 获取指定年的总天数
func GetYearDays(year int) int {
	//平年还是闰年
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return 366
	}
	return 365
}

// GetYearMonthDays 获取指定年、月份的总天数
func GetYearMonthDays(year int, month int) int {
	day31 := map[int]bool{
		1:  true,
		3:  true,
		5:  true,
		7:  true,
		8:  true,
		10: true,
		12: true,
	}
	if day31[month] == true {
		return 31
	}

	// 有30天的月份
	day30 := map[int]bool{
		4:  true,
		6:  true,
		9:  true,
		11: true,
	}
	if day30[month] == true {
		return 30
	}

	//计算平年还是闰年
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		// 得出二月天数
		return 29
	}

	// 得出平年二月天数
	return 28
}
