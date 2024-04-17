package core

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	prefixFileName string
	rxURL          = regexp.MustCompile(`^/regexp\d*`)
	logFileName    string
	logFile        *os.File
	createTime     time.Time

	fileLimit = 1024 * 1024 * 10 // 10Mb
	//fileLimit = 10
	sublog      *zerolog.Logger
	defaultPath string

	// txt file log
	textLogFileName       string
	textLogFile           *os.File
	textLogFileCreateTime time.Time
)

// LoggerConfig ...
type LoggerConfig struct {
	Logger *zerolog.Logger
	// UTC a boolean stating whether to use UTC time zone or local.
	UTC            bool
	SkipPath       []string
	SkipPathRegexp *regexp.Regexp
}

// DefaultLogger ...
func DefaultLogger() gin.HandlerFunc {
	return setLogger(LoggerConfig{
		Logger:         sublog,
		UTC:            true,
		SkipPath:       []string{"/"},
		SkipPathRegexp: rxURL,
	})
}

func setLogger(config ...LoggerConfig) gin.HandlerFunc {
	var newConfig LoggerConfig
	if len(config) > 0 {
		newConfig = config[0]
	}
	var skip map[string]struct{}
	if length := len(newConfig.SkipPath); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range newConfig.SkipPath {
			skip[path] = struct{}{}
		}
	}

	if newConfig.Logger == nil {
		sublog = &(zlog.Logger)
	} else {
		sublog = newConfig.Logger
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()
		track := true

		if _, ok := skip[path]; ok {
			track = false
		}

		if track &&
			newConfig.SkipPathRegexp != nil &&
			newConfig.SkipPathRegexp.MatchString(path) {
			track = false
		}

		if track {
			end := time.Now()
			latency := end.Sub(start)
			if newConfig.UTC {
				end = end.UTC()
			}

			msg := "Request"
			if len(c.Errors) > 0 {
				msg = c.Errors.String()
			}

			dumplogger := sublog.With().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", path).
				Str("ip", c.ClientIP()).
				Dur("latency", latency).
				Str("user-agent", c.Request.UserAgent()).
				Logger()

			switch {
			case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
				{
					dumplogger.Warn().
						Msg(msg)
				}
			case c.Writer.Status() >= http.StatusInternalServerError:
				{
					dumplogger.Error().
						Msg(msg)
				}
			default:
				dumplogger.Info().
					Msg(msg)
			}
		}

	}
}

// InitLog ...
func InitLog(path, name string) {
	prefixFileName = name
	defaultPath = path
	setLogLevel()

	createLumberJackZLogFile(defaultPath + "json/" + prefixFileName + ".json")
	createLumberJackLogFile(defaultPath + "log/" + prefixFileName + ".txt")

	//setZLogFile(defaultPath + "log/" + prefixFileName + "_" + time.Now().Format("01-02") + ".json")
	//setLogFile(defaultPath + "log/" + prefixFileName + "_" + time.Now().Format("01-02") + ".txt")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)

	//go checkFileSize()
	//go checkTextFileSize()
}

// Corsm ...
func Corsm() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Method", "GET, DELETE, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}
}

func setLogFile(filename string) (*os.File, error) {
	fileLog, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		filePathNameArry := strings.Split(filename, "/")
		path := ""
		for i := 0; i < len(filePathNameArry)-1; {
			log.Println(filePathNameArry[i])
			path += filePathNameArry[i]
			i++
		}
		err := os.Mkdir(path, os.FileMode(0755))
		if err != nil {
			return nil, err
		}

		fileLog, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}
	mutiWriter := io.MultiWriter(fileLog, os.Stdout)
	log.SetOutput(mutiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)

	textLogFileCreateTime = time.Now()
	textLogFileName = filename
	textLogFile = fileLog
	return fileLog, nil
}

func setLogLevel() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func setZLogFile(filename string) error {

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		filePathNameArry := strings.Split(filename, "/")
		path := ""
		for i := 0; i < len(filePathNameArry)-1; {
			path += filePathNameArry[i]
			i++
		}
		err := os.Mkdir(path, os.FileMode(0755))
		if err != nil {
			return err
		}

		f, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
	}

	createTime = time.Now()
	logWriter := io.MultiWriter(f, os.Stdout)
	l := zerolog.New(logWriter).With().Timestamp().Logger()
	sublog = &l
	logFileName = filename
	logFile = f
	return nil
}

func checkFileSize() {
	c := time.Tick(time.Minute * 10)
	//c := time.Tick(time.Second * 10)
	for now := range c {
		//log.Println(now.UTC())
		f, err := os.Stat(logFileName)
		if err != nil {
			log.Println(err.Error())
			return
		}

		if int(f.Size()) > fileLimit {
			createNewLogFile()
		} else if createTime.Day() != now.Day() {
			createNewLogFile()
		}
	}
}

func createNewLogFile() {
	logFile.Close()
	// json dir mov
	{
		//temp := strings.Split(logFileName, ".")
		//logBackupFileName := temp[0] + time.Now().Format("01-02_15_04") + ".json"
		logBackupFileName := logFileName

		err := MoveFile(logFileName, strings.Replace(logBackupFileName, "log", "json", 1))
		if err != nil {
			log.Println(err)
		}
		os.Truncate(logFileName, 0)
	}
	setZLogFile(defaultPath + "log/" + prefixFileName + "_" + time.Now().Format("01-02") + ".json")
}

func checkTextFileSize() {
	c := time.Tick(time.Minute * 10)
	//c := time.Tick(time.Second * 10)
	for now := range c {
		//log.Println(now.UTC())
		f, err := os.Stat(textLogFileName)
		if err != nil {
			log.Println(err.Error())
			return
		}

		if int(f.Size()) > fileLimit {
			createNewTextLogFile()
		} else if textLogFileCreateTime.Day() != now.Day() {
			createNewTextLogFile()
		}
	}
}

func createNewTextLogFile() {
	textLogFile.Close()
	// json dir mov
	{
		//temp := strings.Split(logFileName, ".")
		//logBackupFileName := temp[0] + time.Now().Format("01-02_15_04") + ".json"
		logBackupFileName := textLogFileName

		err := MoveFile(textLogFileName, strings.Replace(logBackupFileName, "log", "json", 1))
		if err != nil {
			log.Println(err)
		}
		os.Truncate(textLogFileName, 0)
	}
	setLogFile(defaultPath + "log/" + prefixFileName + "_" + time.Now().Format("01-02") + ".txt")
}

func createLumberJackLogFile(filename string) error {
	l := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
	}
	mutiWriter := io.MultiWriter(l, os.Stdout)
	log.SetOutput(mutiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)

	return nil
}

func createLumberJackZLogFile(filename string) error {
	l := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
	}
	logWriter := io.MultiWriter(l)
	zl := zerolog.New(logWriter).With().Timestamp().Logger()
	sublog = &zl
	return nil
}
