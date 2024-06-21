package sign

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

var YearNotFound = fmt.Errorf("所请求的年份数据不存在")

const yearDataPath = "/%d.json"

type (
	YearData struct {
		Year   int       `json:"year"`
		Papers []string  `json:"papers"`
		Days   []DayData `json:"days"`
	}

	DayData struct {
		Name     string `json:"name"`
		Date     string `json:"date"`
		IsOffDay bool   `json:"isOffDay"`
	}
)

func getDayData(now time.Time) (*DayData, error) {
	// 获取指定年份的数据
	data, err := getJsonData(now.Year())
	if err != nil {
		return nil, err
	}

	// 解析数据
	var yearData YearData
	if err = json.Unmarshal(data, &yearData); err != nil {
		return nil, err
	}

	if len(yearData.Days) == 0 {
		return nil, YearNotFound
	}

	// 判断是否是节假日
	notFormat := now.Format(dateFormat)
	for _, date := range yearData.Days {
		if date.Date == notFormat {
			return &date, nil
		}
	}

	return nil, nil
}

func getJsonData(year int) ([]byte, error) {
	pwd, _ := os.Getwd()
	filePath := path.Join(pwd, fmt.Sprintf(yearDataPath, year))

	// 判断数据文件是否存在
	if exists, _ := FileExists(filePath); exists {
		return os.ReadFile(filePath)
	}

	data, err := downloadJsonFile(year)
	if err != nil {
		return nil, err
	}
	_ = os.WriteFile(filePath, data, 0644)

	return data, nil
}

// downloadJsonFile 下载指定年份的数据
func downloadJsonFile(year int) ([]byte, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/NateScarlet/holiday-cn/master/%d.json", year)
	resp, err := http.Get(url)
	if err != nil {
		return nil, YearNotFound
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
