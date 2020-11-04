package tools

import (
	"crypto/md5"
	"crypto/rc4"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/satori/go.uuid"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
	. "unsafe"
)

func CreateMD5(str []byte)string  {
	ret := md5.Sum(str)
	return hex.EncodeToString(ret[:])
}
func RandNumStrHex() string {

	rand.Seed(rand.Int63())
	rnsh := fmt.Sprintf("%08X", rand.Uint32())
	return rnsh
}
func StrMac2Bytes(mac string) [6]byte {

	ret := [6]byte{0}

	if strings.HasPrefix(mac, "0x") || strings.HasPrefix(mac, "0X") {
		if len(mac) != 14 { //0x112233445566
			return ret
		}
		mac = mac[2:]
	} else {
		if len(mac) != 12 {
			return ret
		}
	}
	for k := 0; k < 6; k++ {
		val, err := strconv.ParseInt(mac[k*2:(k+1)*2], 16, 16)
		if err != nil {
			fmt.Println(err)
			return ret
		}
		ret[k] = byte(val)
	}
	return ret
}
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//非数字字符串自增，前面必须相同 最后一个自增+1，如果255则字符串变长后面为255
func StringAdd(str string) string {
	if str == "" {
		return ""
	}
	start := []byte(str)
	end := start
	//最大值
	if end[len(end)-1] == 255 || (end[len(end)-1] >= '0' && end[len(end)-1] <= '9') {
		return str + string(255)
	} else {
		end[len(end)-1] ++
		return string(end)
	}
}

//yyyy-mm-dd 字符串转int64
func TimeStrToInt64(t string) int64 {

	if t == "" {
		return 0
	}
	//设置本地时区
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05"[:len(t)], t, time.Local)
	time.Unix(tm.Unix(), 0).Format("2006-01-02 15:04:05"[:len(t)])
	return tm.Unix()
}

//yyyy-mm-dd 字符串转int64
func Int64ToTimeStr(tm int64, format string) string {

	if tm == 0 {
		return ""
	}
	return time.Unix(tm, 0).Format(format)
}
func GetMagicString() (ret string) {
	//1.生成一个uuid，如果出错，则用随机数代替

	var magic []byte
	u1, _ := uuid.NewV4()
	if u1 != uuid.Nil {
		magic = u1.Bytes()
	} else {
		rand.Seed(time.Now().UnixNano())
		magic = append(magic, func(n uint64) []byte {
			return []byte{
				byte(n),
				byte(n >> 8),
				byte(n >> 16),
				byte(n >> 24),
				byte(n >> 32),
				byte(n >> 40),
				byte(n >> 48),
				byte(n >> 56),
			}
		}(rand.Uint64())...)
	}

	ret = hex.EncodeToString(magic)
	return
}
func CheckSum(data []byte, length int) uint16 {
	var sum uint32
	var index int
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16

	return func(v uint16) uint16 {
		var i = 0x1
		bs := (*[4]byte)(Pointer(&i))
		if bs[0] == 1 {
			//小端
			var b = make([]byte, 2)
			binary.LittleEndian.PutUint16(b, v)
			return binary.BigEndian.Uint16(b)
		} else {
			//大端
			return v
		}
	}(uint16(^sum))
}
func Rc4Encrypt(in []byte, out []byte, key string) {
	r, _ := rc4.NewCipher([]byte(key))
	r.XORKeyStream(in, out)
}

//base64解码
func Base64Decode(str string) []byte {
	var ret []byte
	enc := base64.StdEncoding
	data, err := enc.DecodeString(str)
	if err != nil {
		fmt.Println("base64 DecodeString failed", err)
		return ret
	}

	return data
}

//base64编码
func Base64Encode(data []byte) string {
	enc := base64.StdEncoding
	str := enc.EncodeToString(data)
	return str
}

//拷贝结构体内容
//dst必须为指针
func CopyObjProperties(dst, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(fmt.Sprintf("%v", e))
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dst type should be a struct pointer")
	}
	if dstValue.IsNil() {
		return fmt.Errorf("dst is nil")
	}

	// src必须为结构体或者结构体指针
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return fmt.Errorf("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}

//本地序转网络序u16
func Htons(v uint16) uint16 {
	var i int = 0x1
	bs := (*[4]byte)(Pointer(&i))
	if bs[0] == 1 {
		//小端
		var b []byte = make([]byte, 2)
		binary.LittleEndian.PutUint16(b, v)
		return binary.BigEndian.Uint16(b)
	} else {
		//大端
		return v
	}
}

//本地序转网络序u16
func Htonsu32(v uint32) uint32 {
	var i int = 0x1
	bs := (*[4]byte)(Pointer(&i))
	if bs[0] == 1 {
		//小端
		var b []byte = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, v)
		return binary.BigEndian.Uint32(b)
	} else {
		//大端
		return v
	}
}

//本地序转网络序 u16

func Ntohs(v uint16) uint16 {
	var i int = 0x1
	bs := (*[4]byte)(Pointer(&i))
	if bs[0] == 1 {
		//小端
		var b []byte = make([]byte, 2)
		binary.BigEndian.PutUint16(b, v)
		return binary.LittleEndian.Uint16(b)
	} else {
		//大端
		return v
	}
}
func Ntohsu32(v uint32) uint32 {
	var i int = 0x1
	bs := (*[4]byte)(Pointer(&i))
	if bs[0] == 1 {
		//小端
		var b []byte = make([]byte, 4)
		binary.BigEndian.PutUint32(b, v)
		return binary.LittleEndian.Uint32(b)
	} else {
		//大端
		return v
	}
}

//生成随机mac
func RandMac() (mac string, err error) {
	buf := make([]byte, 6)

	rand.Seed(time.Now().UnixNano())
	_, err = rand.Read(buf)
	if err != nil {
		mac = "00:15:5D:7C:1F:2F"
		fmt.Println("error:", err)
		return
	}
	// Set the local bit
	buf[0] &= 0xfe
	mac = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	return
}
