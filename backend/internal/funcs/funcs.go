package funcs

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/howeyc/crc16"
	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	day  = 24 * time.Hour
	year = 365 * day
)

var printer = message.NewPrinter(language.English)

var TemplateFuncs = template.FuncMap{
	// Time functions
	"now":            time.Now,
	"timeSince":      time.Since,
	"timeUntil":      time.Until,
	"formatTime":     formatTime,
	"approxDuration": approxDuration,

	// String functions
	"uppercase": strings.ToUpper,
	"lowercase": strings.ToLower,
	"pluralize": pluralize,
	"slugify":   slugify,
	"safeHTML":  safeHTML,

	// Slice functions
	"join":           strings.Join,
	"containsString": slices.Contains[string],

	// Number functions
	"incr":        incr,
	"decr":        decr,
	"formatInt":   formatInt,
	"formatFloat": formatFloat,

	// Boolean functions
	"yesno": yesno,

	// URL functions
	"urlSetParam": urlSetParam,
	"urlDelParam": urlDelParam,
}

func formatTime(format string, t time.Time) string {
	return t.Format(format)
}

func approxDuration(d time.Duration) string {
	if d < time.Second {
		return "less than 1 second"
	}

	ds := int(math.Round(d.Seconds()))
	if ds == 1 {
		return "1 second"
	} else if ds < 60 {
		return fmt.Sprintf("%d seconds", ds)
	}

	dm := int(math.Round(d.Minutes()))
	if dm == 1 {
		return "1 minute"
	} else if dm < 60 {
		return fmt.Sprintf("%d minutes", dm)
	}

	dh := int(math.Round(d.Hours()))
	if dh == 1 {
		return "1 hour"
	} else if dh < 24 {
		return fmt.Sprintf("%d hours", dh)
	}

	dd := int(math.Round(float64(d / day)))
	if dd == 1 {
		return "1 day"
	} else if dd < 365 {
		return fmt.Sprintf("%d days", dd)
	}

	dy := int(math.Round(float64(d / year)))
	if dy == 1 {
		return "1 year"
	}

	return fmt.Sprintf("%d years", dy)
}

func pluralize(count any, singular string, plural string) (string, error) {
	n, err := toInt64(count)
	if err != nil {
		return "", err
	}

	if n == 1 {
		return singular, nil
	}

	return plural, nil
}

func slugify(s string) string {
	var buf bytes.Buffer

	for _, r := range s {
		switch {
		case r > unicode.MaxASCII:
			continue
		case unicode.IsLetter(r):
			buf.WriteRune(unicode.ToLower(r))
		case unicode.IsDigit(r), r == '_', r == '-':
			buf.WriteRune(r)
		case unicode.IsSpace(r):
			buf.WriteRune('-')
		}
	}

	return buf.String()
}

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func incr(i any) (int64, error) {
	n, err := toInt64(i)
	if err != nil {
		return 0, err
	}

	n++
	return n, nil
}

func decr(i any) (int64, error) {
	n, err := toInt64(i)
	if err != nil {
		return 0, err
	}

	n--
	return n, nil
}

func formatInt(i any) (string, error) {
	n, err := toInt64(i)
	if err != nil {
		return "", err
	}

	return printer.Sprintf("%d", n), nil
}

func formatFloat(f float64, dp int) string {
	format := "%." + strconv.Itoa(dp) + "f"
	return printer.Sprintf(format, f)
}

func yesno(b bool) string {
	if b {
		return "Yes"
	}

	return "No"
}

func urlSetParam(u *url.URL, key string, value any) *url.URL {
	nu := *u
	values := nu.Query()

	values.Set(key, fmt.Sprintf("%v", value))

	nu.RawQuery = values.Encode()
	return &nu
}

func urlDelParam(u *url.URL, key string) *url.URL {
	nu := *u
	values := nu.Query()

	values.Del(key)

	nu.RawQuery = values.Encode()
	return &nu
}

func toInt64(i any) (int64, error) {
	switch v := i.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	// Note: uint64 not supported due to risk of truncation.
	case string:
		return strconv.ParseInt(v, 10, 64)
	}

	return 0, fmt.Errorf("unable to convert type %T to int", i)
}




const (
	BounceableTag    = 0x11
	NonBounceableTag = 0x51
	TestFlag         = 0x80
)

type Address struct {
	wc int
	hashPart []byte
	isTestOnly bool
	isBounceable bool
	isUrlSafe bool
	IsUserFriendly bool
}


type ParseResult struct {
	workchain int
	hashPart []byte
	isTestOnly bool
	isBounceable bool
}




func parseFriendlyAddress(address string) (*ParseResult, error) {

	var result ParseResult

	if len(address) != 48 {
		return nil, errors.New("invalid address length 1")
	}


	data := stringToBytes(base64toString(address))
	if len(data) != 36 {
		return nil, errors.New("invalid address length 2")
	}

	// convert to golang code above
	addr := data[:34]
	// crc := data[34:36]

	// crc16 hashsum check


	tag := addr[0]

	if tag == TestFlag {
		result.isTestOnly = true
		tag = tag ^ TestFlag
	}
	if tag != BounceableTag && tag != NonBounceableTag {
		return nil, errors.New("unknown address tag")
	}
	result.isBounceable = tag == BounceableTag

	workchain := int(addr[1])
	if workchain == 0xff {
		result.workchain = -1
	}
	if workchain != 0 && workchain != -1 {
		return nil, errors.New("invalid address wc")
	}

	result.hashPart = addr[2:34]
	return &result, nil
}


// convert to golang code above
func NewAddress(anyForm string) (*Address, error) {
	var addr Address
	if strings.Contains(anyForm, "-") || strings.Contains(anyForm, "_") {
		addr.isUrlSafe = true
		anyForm = strings.Replace(anyForm, "-", "+", -1)
		anyForm = strings.Replace(anyForm, "_", "/", -1)
	} else {
		addr.isUrlSafe = false
	}
	if strings.Contains(anyForm, ":") {
		arr := strings.Split(anyForm, ":")
		if len(arr) != 2 {
			return nil, errors.New("Invalid address " + anyForm)
		}
		var wc int
		wc, err := strconv.Atoi(arr[0])
		if err != nil {
			return nil, errors.New("Invalid address wc " + anyForm)
		}

		if wc != 0 && wc != -1 {
			return nil, errors.New("Invalid address wc " + anyForm)
		}
		hex := arr[1]
		if len(hex) != 64 {
			return nil, errors.New("Invalid address hex " + anyForm)
		}
		addr.IsUserFriendly = false
		addr.wc = wc
		addr.hashPart = hexToBytes(hex)
		addr.isTestOnly = false
		addr.isBounceable = false


		return &addr, nil
	} else {
		addr.IsUserFriendly = true
		var parseResult *ParseResult
		parseResult, err := parseFriendlyAddress(anyForm)
		if err != nil {
			return nil, err
		}

		addr.wc = parseResult.workchain
		addr.hashPart = parseResult.hashPart
		addr.isTestOnly = parseResult.isTestOnly
		addr.isBounceable = parseResult.isBounceable

		return &addr, nil
	}
}



// convert code above to golang code
func hexToBytes(s string) []byte {
	var to_hex_array = []string{}
	// convert to golang code above
	var to_byte_map = map[string]byte{}

	for i := 0; i <= 0xff; i++ {
		var s = fmt.Sprintf("%02x", i)
		to_hex_array = append(to_hex_array, s)
		to_byte_map[s] = byte(i)
	}

	s = strings.ToLower(s)
	length2 := len(s)
	if length2 % 2 != 0 {
		panic("hex string must have length a multiple of 2")
	}
	length := length2 / 2
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		i2 := i * 2
		b := s[i2:i2+2]
		if !hasOwnProperty(to_byte_map, b) {
			panic("invalid hex character " + b)
		}
		result[i] = to_byte_map[b]
	}
	return result
}



func (a *Address) ToString() string {
	if a.IsUserFriendly {
		return fmt.Sprintf("%d:%s", a.wc, bytesToHex(a.hashPart))
	} else {
		var tag = a.isBounceable
		
		var addr = make([]byte, 34)

	
		addr[0] = byte(boolToByte(tag))
		addr[1] = byte(a.wc)
		copy(addr[2:], a.hashPart)
		var addressWithChecksum = make([]byte, 36)

		copy(addressWithChecksum, addr)
		copy(addressWithChecksum[34:], uint16ToByte(crc16.ChecksumCCITTFalse(addr)))
		


		var addressBase64 = stringToBase64(string(addressWithChecksum))
		if a.isUrlSafe {
			addressBase64 = strings.Replace(addressBase64, "+", "-", -1)
			addressBase64 = strings.Replace(addressBase64, "/", "_", -1)
		}
		return addressBase64
	}
}


func hasOwnProperty(m map[string]byte, key string) bool {
	_, ok := m[key]
	return ok
}


func stringToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// convert bool to byte
func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func uint16ToByte(uintvariable uint16) []byte {
	var bytes = make([]byte, 2)
	bytes[0] = byte(uintvariable >> 8)
	bytes[1] = byte(uintvariable & 0xFF)
	return bytes
}




func bytesToHex(buffer []byte) string {
	var to_hex_array = []string{}
	// convert to golang code above
	var to_byte_map = map[string]byte{}

	for i := 0; i <= 0xff; i++ {
		var s = fmt.Sprintf("%02x", i)
		to_hex_array = append(to_hex_array, s)
		to_byte_map[s] = byte(i)
	}

	var hex_array = make([]string, len(buffer))
	for i := 0; i < len(buffer); i++ {
		hex_array[i] = to_hex_array[buffer[i]]
	}
	return strings.Join(hex_array, "")
}


func stringToBytes(s string) []byte {
	b := make([]byte, len(s))
	for i := range b {
		b[i] = s[i]
	}
	return b
}
// base64 to string
func base64toString(s string) string {
	var b, _ = base64.StdEncoding.DecodeString(s)
	return string(b)
}