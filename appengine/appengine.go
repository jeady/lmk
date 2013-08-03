// +build appengine

package appengine

import (
  "fmt"
  "net/http"
  "time"

  "appengine"
  "appengine/datastore"

  "github.com/jeady/go-logging"
  lmk "github.com/jeady/lmk/engine"
)

var engine *lmk.Engine

var lmk_config_file string = "lmk.conf"

// Used for storing configuration values in the datastore.
type Config struct {
  Key   string
  Value string
}

func init() {
  http.HandleFunc("/", status)
  http.HandleFunc("/test-all", test_all)
  http.HandleFunc("/list", list)
  http.HandleFunc("/test-notify", test_notify)
  http.HandleFunc("/poll", poll)
}

func begin_lmk_request(r *http.Request) (e *lmk.Engine, c appengine.Context) {
  c = appengine.NewContext(r)

  // NOT SAFE FOR CONCURRENT REQUESTS.
  // Fortunately, this application should not require the use of concurrent
  // requests.
  logging.SetBackend(logging.NewGAELogger(c))

  e = lmk.NewEngineFromFile(lmk_config_file)
  if e == nil {
    panic("Could not start lmk engine.")
  }

  from, _, _ := e.Config().SmtpConfig()
  e.SetDefaultNotifier(NewGaeMailNotifier(from, c))
  e.SetUrlFetcher(NewGaeUrlFetcher(c))

  return
}

func status(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain")
  fmt.Fprintln(w, "This is an lmk! server running on Google App Engine.")
  fmt.Fprintln(w, "Find out more at github.com/jeady/lmk.")
}

func test_all(w http.ResponseWriter, r *http.Request) {
  engine, _ := begin_lmk_request(r)
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
  engine, _ := begin_lmk_request(r)
  w.Header().Set("Content-Type", "text/plain")

  fmt.Fprintln(w, "Enabled rules:\n")
  for _, rule := range engine.Rules() {
    fmt.Fprintln(w, "  * "+rule.Name())
  }
}

func test_notify(w http.ResponseWriter, r *http.Request) {
  engine, _ := begin_lmk_request(r)
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
  engine, context := begin_lmk_request(r)
  w.Header().Set("Content-Type", "text/plain")

  var configs []Config
  q := datastore.NewQuery("Config").Filter("Key =", "last_update")
  q.GetAll(context, &configs)
  if len(configs) < 1 {
    context.Infof("No value for last_update, running all.")
    for _, rule := range engine.Rules() {
      fmt.Fprintln(w, "Running "+rule.Name())
      engine.Run(rule)
    }
  } else {
    var last_update time.Time
    last_update.UnmarshalJSON([]byte(configs[0].Value))

    fmt.Fprintln(w, "Last updated "+last_update.String())
    for _, rule := range engine.RulesToPoll(last_update) {
      context.Infof("Running " + rule.Name())
      fmt.Fprintln(w, "Running "+rule.Name())
      engine.Run(rule)
    }
  }

  bytes, _ := time.Now().MarshalJSON()
  c := Config{
    Key:   "last_update",
    Value: string(bytes),
  }
  bKey := datastore.NewKey(context, "Config", "Key", 0, nil)
  _, err := datastore.Put(context, bKey, &c)

  if err != nil {
    fmt.Fprintln(w, "Problem putting: "+err.Error())
  }
}
