package log

import "testing"
import "os"
import "io/ioutil"
import "fmt"
import "time"

func TestLoggerStdout(t *testing.T) {
    var config Config
    config.Layout        = LY_DEFAULT
    config.LayoutStyle   = LS_DEFAULT
    config.Level         = INFO
    config.TimeFormat    = TF_LONG
    config.Utc           = false

    logger, err := New(os.Stdout, config)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    defer logger.Wait()

    logger.Debug("Test Debug(). This line will not output.")
    logger.Debugf("Test Debugf(). This line will not output.")
    logger.Info("Test Info().")
    logger.Infof("Test Infof(): %s, %d", "string", 123)
    logger.Warn("Test Warn(). ", "warning")
    logger.Warnf("Test Warnf(): %s", "warning")
    logger.Error("Test Error().")
    logger.Errorf("Test Errorf().")

    logger.Print(WARN, "Test Print().")
    logger.Printf(WARN, "Test Printf().")
    logger.Write([]byte("Test Write()."))
}


func TestLoggerStdout2(t *testing.T) {
    var config Config
    config.Rotate = R_HOURLY
    config.RotatePattern = RP_DEFAULT
    _, err := New(os.Stdout, config)
    if err == nil {
        t.Fail()
    }
    if err.Error() != "Stdout or stderr could not be rotated." {
        t.Error(err)
        t.Fail()
    }
}


func TestLoggerStdout3(t *testing.T) {
    var config Config
    config.Rotate = R_NONE
    config.RotatePattern = RP_DEFAULT
    _, err := New(os.Stdout, config)
    if err == nil {
        t.Fail()
    }
    if err.Error() != "Rotate pattern does not match rotate value." {
        t.Error(err)
        t.Fail()
    }
}


// Test to write log to a temporary file.
func TestLoggerFile(t *testing.T) {
    file, err := ioutil.TempFile("", "test_log_")
    if err != nil {
        t.Error(err)
        t.Fail()
    }

    filename := file.Name()
    t.Logf("Create temp file: %s", filename)

    var config Config
    config.Layout        = LY_DEFAULT
    config.LayoutStyle   = LS_DEFAULT
    config.Level         = INFO
    config.TimeFormat    = TF_TIMELONG
    config.Utc           = true

    logger, err := New(file, config)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    logger.Debug("Test Debug().")
    logger.Info("Test Info().")
    logger.Write([]byte("Test Write()."))

    t.Logf("Write log message to temp file: %s", filename)

    err = file.Close()
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    t.Logf("Close temp file: %s", filename)

    err = os.Remove(filename)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    t.Logf("Remove temp file: %s", filename)

    logger.Wait()
}


func TestSimpleLogger(t *testing.T) {
    Output("std logger")
    Outputf("std %s", "logger")
}


func TestLevelType(t *testing.T) {
    var ok bool
    _, ok = String2Level("DEBUG")
    if !ok {
        t.Fail()
    }
    _, ok = String2Level("notice")
    if !ok {
        t.Fail()
    }
    _, ok = String2Level("")
    if ok {
        t.Fail()
    }
    _, ok = String2Level("haha")
    if ok {
        t.Fail()
    }
}


func ExampleNew() {

    file, err := ioutil.TempFile("", "test_log_")
    if err != nil {
        fmt.Println(os.Stderr, err)
        os.Exit(1)
    }

    var config Config
    config.Level = INFO
    config.Rotate = R_HOURLY

    logger, err := New(file, config)
    if err != nil {
        fmt.Println(os.Stderr, err)
        os.Exit(1)
    }
    defer logger.Wait()

    logger.Info("something happens")
}


func ExampleNew_another() {

    var config Config
    config.Layout        = LY_LEVEL
    config.LayoutStyle   = "{level}: {msg}"
    config.Level         = DEBUG
    config.TimeFormat    = TF_DEFAULT
    config.Utc           = false

    logger, err := New(os.Stdout, config)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    logger.Debug("Test Debug().")
    logger.Notice("Test Notice().")
    logger.Noticef("Test Noticef(). %s", "xyz")
    logger.Info("Test Info(). ", "abc", 456)
    logger.Warnf("Test Warn(). %s", "string")
    logger.Wait()
    // Output: DEBUG: Test Debug().
    // NOTICE: Test Notice().
    // NOTICE: Test Noticef(). xyz
    // INFO: Test Info(). abc456
    // WARN: Test Warn(). string
}


func ExampleNew_handle() {

    // Define a function to do some custom work.
    var handle Handle
    handle.Func = func(msg Message) {
        time.Sleep(2 * time.Second)
        if msg.Level >= WARN {
            fmt.Println("haha:", msg.Msg)
        }
    }
    handle.Level = WARN

    var config Config
    config.Layout = LY_MSGONLY
    config.LayoutStyle = "{msg}"
    config.Level = INFO

    logger, err := New(os.Stdout, config, handle)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer logger.Wait()
    logger.Debug("debug")
    logger.Info("info")
    logger.Warn("warn")
    logger.Error("error")
    // Output: info
    // warn
    // error
    // haha: warn
    // haha: error
}


// Benchmark for writing log to stdout.
func BenchmarkLoggerStdout(b *testing.B) {

    b.StopTimer()

        var config Config
        config.Layout        = LY_TIME
        config.LayoutStyle   = LS_SIMPLE
        config.Level         = DEBUG
        config.TimeFormat    = TF_DEFAULT
        config.Utc           = true

        logger, err := New(os.Stdout, config)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

    b.StartTimer()

        for i := 0; i < b.N; i++ {
            logger.Debug("Test writing log to stdout. The more you write, the slower the speed is.")
        }

    b.StopTimer()

        logger.Wait()
}


// Benchmark for writing log to file.
func BenchmarkLoggerFile(b *testing.B) {

    b.StopTimer()

        file, err := ioutil.TempFile("", "test_log_")
        if err != nil {
            b.Error(err)
            b.Fail()
        }

        filename := file.Name()

        var config Config
        config.Layout        = LY_TIME
        config.LayoutStyle   = LS_SIMPLE
        config.Level         = DEBUG
        config.TimeFormat    = TF_DEFAULT
        config.Utc           = true

        logger, err := New(file, config)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

    b.StartTimer()

        for i := 0; i < b.N; i++ {
            logger.Debug("Test writing log to file. The more you write, the slower the speed is.")
        }

    b.StopTimer()

        err = os.Remove(filename)
        if err != nil {
            b.Error(err)
            b.Fail()
        }

        logger.Wait()
}


func BenchmarkMsg2bytes(b *testing.B) {
    b.StopTimer()

        var m Message
        m.Time = time.Now()
        m.Msg = "Test message. Test message. Test message. Test message. Test message. Test message."
        m.Level = INFO

        logger, err := New(os.Stdout, Config{})
        if err != nil {
            b.Error(err)
            b.Fail()
        }

    b.StartTimer()

        for i := 0; i < b.N; i++ {
            logger.msg2bytes(m)
        }
}
