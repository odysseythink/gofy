package utils

import (
	"encoding/hex"
	"math"
	"math/rand"
	v2rand "math/rand/v2"
	"net"
	"runtime"
	"strings"

	uuid "github.com/satori/go.uuid"
	"mlib.com/mlog"
)

func RoundWithPrecision(x float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Round(x*p) / p
}

func GetIP() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		mlog.Errorf("net.Interfaces() failed:%v", err)
		return "0.0.0.0"
	}
	for _, intf := range netInterfaces {
		if (intf.Flags&net.FlagUp) != 0 && (intf.Flags&net.FlagLoopback) == 0 {
			addrs, _ := intf.Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.To4().String()
					}
				}
			}
		}
	}
	// addrs, err := net.InterfaceAddrs()
	// if err != nil {
	// 	mlog.Errorf("InterfaceAddrs failed:%v", err)
	// 	return ""
	// }
	// fmt.Println("")
	// for _, address := range addrs {
	// 	// 检查ip地址判断是否回环地址
	// 	if ipnet, ok := address.(*net.IPNet); ok {
	// 		fmt.Println("----", ipnet, "--IsLoopback=", ipnet.IP.IsLoopback(), "--IsMulticast=", ipnet.IP.IsMulticast(), "------IsInterfaceLocalMulticast=", ipnet.IP.IsInterfaceLocalMulticast(), "------IsLinkLocalMulticast=", ipnet.IP.IsLinkLocalMulticast(), "----IsLinkLocalUnicast=", ipnet.IP.IsLinkLocalUnicast(), "---IsGlobalUnicast", ipnet.IP.IsGlobalUnicast())
	// 	}
	// }
	// for _, address := range addrs {
	// 	// 检查ip地址判断是否回环地址
	// 	if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	// 		if ipnet.IP.To4() != nil {
	// 			// fmt.Println(ipnet.IP.String())
	// 			return ipnet.IP.To4().String()
	// 		}
	// 	}
	// }
	return "0.0.0.0"
}

// 把sql中的 '\' 和 '%' 替换为转义后的字符串，方便模糊查询
func ReplaceSqlSpecialCharacter(sql_str string) string {
	sql_str = strings.Replace(sql_str, "\\", "\\\\", -1)
	sql_str = strings.Replace(sql_str, "%", "\\%", -1)
	return sql_str
}

var (
	defaultStackSize = 40960
)

func GetCurrentGoroutineStack() string {
	buf := make([]byte, defaultStackSize)
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

const (
	whitespace      = " \t\n\r\v\f"
	ascii_lowercase = "abcdefghijklmnopqrstuvwxyz"
	ascii_uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ascii_letters   = ascii_lowercase + ascii_uppercase
	digits          = "0123456789"
	hexdigits       = digits + "abcdef" + "ABCDEF"
	octdigits       = "01234567"
)

func GenerateString(n int) string {
	const lettersDigits = ascii_letters + digits
	result := strings.Builder{}

	for i := 0; i < n; i++ {
		result.WriteByte(lettersDigits[rand.Intn(len(lettersDigits))])
	}

	return result.String()
}
func GenerateRandomHex(n int) string {
	randomBytes := make([]byte, n)
	_, err := (v2rand.NewChaCha8([32]byte([]byte(uuid.NewV4().String())))).Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(randomBytes)
}
