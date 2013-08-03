package main

import (
  "fmt"
  "net/http"
  "net/url"
  "os"
  "time"

  "github.com/garyburd/redigo/redis"
  lmk "github.com/jeady/lmk/engine"
)

var engine *lmk.Engine
var lmk_config_file string = "lmk.conf"

func init() {
  http.HandleFunc("/", status)
  http.HandleFunc("/test-all", test_all)
  http.HandleFunc("/list", list)
  http.HandleFunc("/test-notify", test_notify)
  http.HandleFunc("/poll", poll)
}

func status(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain")
  fmt.Fprintln(w, "This is an lmk! server running on Google App Engine.")
  fmt.Fprintln(w, "Find out more at github.com/jeady/lmk.")
}

func test_all(w http.ResponseWriter, r *http.Request) {
  engine = lmk.NewEngineFromFile(lmk_config_file)
  w.Header().Set("Content-Type", "text/plain")

  for _, rule := range engine.Rules() {
    sane, triggered := rule.TestTriggered()
    if !sane {
      fmt.Fprintln(w, rule.Name()+" is not sane.")
      continue
    }

    fmt.Fprintln(w, rule.Name()+" is sane.")
    if !triggered {
      fmt.Fprintln(w, rule.Name()+" has not been triggered.")
    } else {
      fmt.Fprintln(w, rule.Name()+" has been triggered.")
    }
    fmt.Fprintln(w, "")
  }
}

func list(w http.ResponseWriter, r *http.Request) {
  engine = lmk.NewEngineFromFile(lmk_config_file)
  w.Header().Set("Content-Type", "text/plain")

  fmt.Fprintln(w, "Enabled rules:")
  for _, rule := range engine.Rules() {
    fmt.Fprintln(w, "  * "+rule.Name())
  }
}

func test_notify(w http.ResponseWriter, r *http.Request) {
  engine = lmk.NewEngineFromFile(lmk_config_file)
  w.Header().Set("Content-Type", "text/plain")

  r.ParseForm()
  if _, ok := r.Form["to"]; !ok {
    fmt.Fprintln(w, "Must specify a 'to' query parameter.")
    return
  }

  err := engine.DefaultNotifier().Notify(
    r.Form.Get("to"),
    "test notification",
    "lmk! is correctly configured to send notifications from GAE. Cheers!")

  if err != nil {
    fmt.Fprintln(w, "Problem sending notifications.")
    return
  }
  fmt.Fprintln(w, "Successfully sent notifications.")
}

func poll(w http.ResponseWriter, r *http.Request) {
  engine = lmk.NewEngineFromFile(lmk_config_file)
  w.Header().Set("Content-Type", "text/plain")

  u, _ := url.Parse(os.Getenv("REDISTOGO_URL"))
  c, err := redis.Dial("tcp", u.Host)
  if err != nil {
    fmt.Println("Could not connect to redis.")
    return
  }

  pass, _ := u.User.Password()
  c.Do("AUTH", pass)

  exists, _ := redis.Bool(c.Do("EXISTS", "last_update"))
  if !exists {
    fmt.Println(w, "No value for last_update, running all.")
    for _, rule := range engine.Rules() {
      fmt.Fprintln(w, "Running "+rule.Name())
      engine.Run(rule)
    }
  } else {
    var last_update time.Time
    last_update_str, _ := redis.String(c.Do("GET", "last_update"))
    last_update.UnmarshalJSON([]byte(last_update_str))

    fmt.Fprintln(w, "Last updated "+last_update.String())
    for _, rule := range engine.RulesToPoll(last_update) {
      fmt.Println("Running " + rule.Name())
      fmt.Fprintln(w, "Running "+rule.Name())
      engine.Run(rule)
    }
  }

  bytes, _ := time.Now().MarshalJSON()
  c.Send("SET", "last_update", string(bytes))
  c.Flush()
  c.Receive()
}

func main() {
  err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
  if err != nil {
    panic(err)
  }
}
