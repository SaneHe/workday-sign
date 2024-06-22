package sign

import (
	"fmt"
	"github.com/sanehe/workday-sign/internal/util"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// CmdSign is sign-in tool
var CmdSign = &cobra.Command{
	Use:   "sign",
	Short: "sign: 打卡签到",
	Long:  `sign: 打卡签到`,
	Run:   run,
}

var (
	dateFormat = "2006-01-02"
	apiURL     string
	key        string
	iv         string
	body       = `{"title": "今天是 %s %s %s", "body": "打卡时间到啦! 卡卡卡", "sound": "birdsong", "url": "wxwork://"}`
	weekdays   = [...]string{
		"周日", "周一", "周二", "周三", "周四", "周五", "周六",
	}
)

func init() {
	if apiURL = os.Getenv("API_URL"); apiURL == "" {
		apiURL = "https://api.day.app/"
	}
	
	CmdSign.Flags().StringVarP(&apiURL, "api-url", "r", apiURL, "bark server api")
	CmdSign.Flags().StringVarP(&body, "body", "b", body, "推送加密消息内容")
	CmdSign.Flags().StringVarP(&key, "key", "k", key, "推送加密 key")
	CmdSign.Flags().StringVarP(&iv, "iv", "i", iv, "推送加密初始向量")
}

func run(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		panic("参数不足: 请输入推送设备标识")
	}
	
	now, deviceKey := time.Now(), args[0]
	dayOfWeek, nowFormat := now.Weekday(), now.Format(dateFormat)
	
	dayData, _ := getDayData(now)
	weekName := weekdays[dayOfWeek]
	
	if dayData != nil && !dayData.IsOffDay {
		fmt.Printf("今天是: %s %s 为 %s, 签到结果: %v\n", nowFormat, weekName, dayData.Name, notify(deviceKey, nowFormat, weekName, dayData.Name))
		return
	}
	
	// 周末 或 节假日
	if (dayData != nil && dayData.IsOffDay) || dayOfWeek == time.Saturday || dayOfWeek == time.Sunday {
		fmt.Printf("今天是: %s %s, 休息啦\n", nowFormat, weekName)
		return
	}
	
	fmt.Printf("今天是: %s %s 为 %s, 签到结果: %v\n", nowFormat, weekName, "工作日", notify(deviceKey, nowFormat, weekName, "工作日"))
}

func notify(deviceKey, date, weekName, dayName string) error {
	body = fmt.Sprintf(body, date, weekName, dayName)
	// 加密推送参数为空
	if key != "" && iv != "" {
		return util.SendToBark(apiURL, deviceKey, body, key, iv)
	}
	
	return util.SendToBarkDirectly(apiURL, deviceKey, body)
}
