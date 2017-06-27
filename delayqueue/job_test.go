package delayqueue

import "testing"

func TestJob (t *testing.T) {
    key := "putjob"
    value := "putjob-value"
    err := putJob(key, value)
    if err != nil {
        t.Fail()
    }
    newValue, err := getJob(key)
    if err != nil {
        t.Fail()
    }
    if newValue != value {
        t.Fail()
    }

    err = removeJob(key)
    if err != nil {
        t.Fail()
    }
}