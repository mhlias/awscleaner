package resources

import (
  "time"
)


func older_than(created int64, days int16) bool {
  
  seconds_now := time.Now().Unix()

  days_since_created := (seconds_now - created)/(3600*24)

  if days_since_created > int64(days) {
    return true
  }

  return false


}


func is_whitelisted(list []string , id string) bool {
    for _, b := range list {
        if b == id {
            return true
        }
    }
    return false
}