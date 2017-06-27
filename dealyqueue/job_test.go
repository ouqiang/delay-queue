package delayqueue

import "testing"

func TestJob (t *testing.T) {
    key := "putjob"
    value := "putjob-value"
    err := PutJob(key, value)
    if err != nil {
        t.Fail()
    }
    newValue, err := GetJob(key)
    if err != nil {
        t.Fail()
    }
    if newValue != value {
        t.Fail()
    }

    err = RemoveJob(key)
    if err != nil {
        t.Fail()
    }
}