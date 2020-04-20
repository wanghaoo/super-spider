package main

import (
  _ "github.com/go-sql-driver/mysql"
  "net/http"
  "io/ioutil"
  "errors"
  "fmt"
  "encoding/json"
  "strconv"
  "time"
  sq "github.com/Masterminds/squirrel"
  "database/sql"
  "math/rand"
)

type ListResponse struct {
  Attr    Attr   `json:"attr"`
  Code    int    `json:"code"`
  Msg     string `json:"msg"`
  Success bool   `json:"success"`
}

type Attr struct {
  Datas []Data `json:"datas"`
}

type Data struct {
  SequenceId string `json:"sequenceId"`
}

type Detail struct {
  Attr DetailAttr `json:"attr"`
}

type DetailAttr struct {
  EventResultMap EventResultMap `json:"eventResultMap"`
  ActivityMap    ActivityMap    `json:"activityMap"`
}

type EventResultMap struct {
  RiskStatus string `json:"riskStatus"`
  RiskScore  string `json:"riskScore"`
}

type ActivityMap struct {
  AccountMobile     string `json:"accountMobile"`
  AccountName       string `json:"accountName"`
  MobileAddressCity string `json:"mobileAddressCity"`
  IdNumber          string `json:"idNumber"`
}

var (
  db         *sql.DB
  DBProxy    sq.DBProxyBeginner
)

const DBURL = "friends:Klfq_fds1@tcp(122.112.157.77:3306)/friends?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
func OpenDB() sq.DBProxyBeginner {
  var err error
  db, err = sql.Open("mysql", DBURL)
  if err != nil {
    panic(err)
  }

  db.SetMaxOpenConns(300)
  db.SetMaxIdleConns(300)
  db.SetConnMaxLifetime(10 * time.Second)


  err = db.Ping()
  if err != nil {
    panic(err)
  }

  DBProxy = sq.NewStmtCacheProxy(db)

  return DBProxy
}

func RandRange(min int, max int) int {
  return rand.Intn(max - min) + min
}

func main() {
  OpenDB()
  var i = 1
  for {
    url := fmt.Sprintf("https://oceanus.tongdun.cn/ruleengine/activity/history.json?operationType=doSearch&eventType=&policySetName=All&riskType=&riskStatus=&downRiskScore=&upRiskScore=&searchField=accountLogin&searchValue=&startDate=1587225600000&endDate=1587277474766&pageSize=50&curPage=%d&totalCount=32912&tdTraceId=1586865563503", i)
    i ++
    response, err := doHttp(url, "GET", nil, nil)
    if err != nil {
      fmt.Println(err)
    }
    data := ListResponse{}
    json.Unmarshal([]byte(response), &data)
    for _, d := range data.Attr.Datas {
      second := RandRange(5 , 10)
      time.Sleep(time.Duration(second) * time.Second)
      detailUrl := fmt.Sprintf("https://oceanus.tongdun.cn/ruleengine/activity/history.json?operationType=showDetail&evnetType=&sequenceId=%s&tdTraceId=%d", d.SequenceId, time.Now().UnixNano() / 1e6)
      detailJson, err := requestDetail(detailUrl, "GET")
      if err != nil {
        fmt.Println(err)
      }
      if len(detailJson) <= 0 {
        fmt.Println("======== it is end =========")
        return
      }
      detail := Detail{}
      json.Unmarshal([]byte(detailJson), &detail)
      if len(detail.Attr.EventResultMap.RiskScore) > 0 {
        score, err := strconv.Atoi(detail.Attr.EventResultMap.RiskScore)
        if err != nil {
          continue
        }
        if score <= 20 {
          _, err = sq.Insert("tongdun_data").Columns("accountName, accountMobile, mobileAddressCity, idNumber, riskStatus, riskScore").
            Values(detail.Attr.ActivityMap.AccountName, detail.Attr.ActivityMap.AccountMobile, detail.Attr.ActivityMap.MobileAddressCity, detail.Attr.ActivityMap.IdNumber, detail.Attr.EventResultMap.RiskStatus, detail.Attr.EventResultMap.RiskScore).
            RunWith(db).Exec()
          if err != nil {
            fmt.Println("insert into mysql error: ", err.Error())
          }
        }
      }
    }
  }
}

func requestDetail(url string, method string) (string, error) {
  client := &http.Client{}
  request, err := http.NewRequest(method, url, nil)
  if request == nil {
    return "", errors.New("build http request error")
  }
  request.Header.Set("Host", "oceanus.tongdun.cn")
  cookie := new(http.Cookie)
  cookie.Name = "TSESSIONID"
  cookie.Value = "DCYJID6V-P90F8VJL0XUET9CP75UR2-8CBAP69K-N2"
  request.AddCookie(cookie)
  response, err := client.Do(request)
  if err != nil {
    return "", err
  }
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return "", err
  }
  return string(body), nil
}

func doHttp(url string, method string, param map[string]interface{}, auth map[string]string) (string, error) {
  client := &http.Client{}
  request, err := http.NewRequest(method, url, nil)
  if request == nil {
    return "", errors.New("build http request error")
  }
  request.Header.Set("Host", "oceanus.tongdun.cn")
  cookie := new(http.Cookie)
  cookie.Name = "TSESSIONID"
  cookie.Value = "DCYJID6V-P90F8VJL0XUET9CP75UR2-8CBAP69K-N2"
  request.AddCookie(cookie)
  response, err := client.Do(request)
  if err != nil {
    return "", err
  }
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return "", err
  }
  return string(body), nil
}
