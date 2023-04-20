package client

import (
	"fmt"
	"github.com/cihub/seelog"
	"gorm.io/gorm"
	"math"
	"muchat-go/chatgpt/api_base"
	"muchat-go/config"
	"muchat-go/models"
	"sync"
	"time"
)

type Client struct {
    Db           *gorm.DB
    Accounts     []models.OpenAiAccount
    accIndex     int
    accIndexLock sync.Mutex
    Concurrency  int
}

func NewClient(db *gorm.DB) (r *Client) {
    r = new(Client)
    r.Db = db
    r.Accounts = make([]models.OpenAiAccount, 0, 0)
    r.accIndex = -1
    r.Concurrency = 1
    return
}

func (r *Client) LoadAccounts() {
    r.accIndexLock.Lock()
    defer r.accIndexLock.Unlock()

    r.Accounts = []models.OpenAiAccount{}
    // status = 0 不能传 OpenAiAccount，因为 status 字段的默认值就是 0
    err := r.Db.Debug().Where("status = 0").Find(&r.Accounts).Error
    if err != nil {
        panic(fmt.Errorf("load open ai account failed, err=%w", err))
    }
    // 一次只用一个号
    if len(r.Accounts) > r.Concurrency {
        r.Accounts = r.Accounts[:r.Concurrency]
    }
    seelog.Infof("load open ai accounts: %v", len(r.Accounts))
}

func (r *Client) LoadAccountIntoDb() {
    accounts := config.LoadAccounts()
    for _, acc := range accounts {
        md := models.OpenAiAccount{}
        err := r.Db.Unscoped().Where(models.OpenAiAccount{Email: acc.Email}).First(&md).Error
        if err == gorm.ErrRecordNotFound {
            md.Email = acc.Email
            md.ApiKey = acc.ApiKey
            md.Password = acc.Password
            md.FirstTime = time.Unix(0, 0)
            md.ExpiresAt = time.Unix(0, 0)
            err = r.Db.Save(&md).Error
            if err != nil {
                panic(fmt.Errorf("save open ai account failed, err=%w", err))
            }
        } else if err != nil {
            panic(fmt.Errorf("query open ai account failed, err=%w", err))
        }
    }
}

func (r *Client) PickAccount() *models.OpenAiAccount {
    r.accIndexLock.Lock()
    defer r.accIndexLock.Unlock()

    if len(r.Accounts) == 0 {
        return nil
    }

    if r.accIndex+1 >= len(r.Accounts) {
        r.accIndex = -1
    }
    r.accIndex++
    acc := r.Accounts[r.accIndex]
    return &acc
}

func (r *Client) RefreshCredit() {
    seelog.Infof("refreshing credit")
    accs := []models.OpenAiAccount{}
    for _, acc := range r.Accounts {
        jd, err := api_base.CreditGrants(acc.ApiKey)
        if err != nil {
            seelog.Errorf("get credit grants failed, err=%v, email=%v", err, acc.Email)
            // 查询失败不代表账号不能用，先加进去
            accs = append(accs, acc)
            continue
        }
        expiresAt := time.Unix(int64(jd.Grants.Data[0].ExpiresAt), 0)

        acc.CreditUsed = jd.TotalUsed
        acc.CreditAvailable = jd.TotalAvailable
        acc.ExpiresAt = expiresAt

        err = r.Db.Model(acc).Updates(models.OpenAiAccount{
            CreditUsed:      jd.TotalUsed,
            CreditAvailable: jd.TotalAvailable,
            ExpiresAt:       expiresAt,
        }).Error
        if err != nil {
            seelog.Errorf("update open ai account failed, err=%v, email=%v", err, acc.Email)
            // 失败不代表账号不能用，先加进去
            accs = append(accs, acc)
            continue
        }
    }

    if len(accs) == 0 {
        seelog.Errorf("!! no available open ai account !!")
        return
    }

    r.accIndexLock.Lock()
    r.Accounts = accs
    r.accIndexLock.Unlock()
}

func (r *Client) RefreshUsage() {
    seelog.Infof("refreshing usage")
    accs := []models.OpenAiAccount{}
    for _, acc := range r.Accounts {
        // CurrentUsageUsd 变成 0 了
        //jd, err := api_base.Usage(acc.ApiKey)
        //if err != nil {
        //    seelog.Errorf("get usage failed, err=%v, email=%v", err, acc.Email)
        //    // 查询失败不代表账号不能用，先加进去
        //    accs = append(accs, acc)
        //    continue
        //}
        //seelog.Infof("usage: acc=%v, object=%v, currentUsageUsd=%v", acc.Email, jd.Object, jd.CurrentUsageUsd)
        //acc.UsdSpent = jd.CurrentUsageUsd

        jd, err := api_base.BillingUsage(acc.ApiKey)
        if err != nil {
            seelog.Errorf("get usage failed, err=%v, email=%v", err, acc.Email)
            // 查询失败不代表账号不能用，先加进去
            accs = append(accs, acc)
            continue
        }
        seelog.Infof("usage: acc=%v, object=%v, currentUsageUsd=%.2f", acc.Email, jd.Object, jd.TotalUsage/100)
        acc.UsdSpent = jd.TotalUsage / 100

        usdSpentLimit := 120.0
        if acc.UsdSpentLimit > 0 {
            usdSpentLimit = acc.UsdSpentLimit
        }
        gap := math.Abs(usdSpentLimit - acc.UsdSpent)
        status := acc.Status
        if gap >= 0.5 {
            accs = append(accs, acc)
        } else {
            // 先标记起来，实际还在继续用
            status = models.OpenAiAccountStatusPaused
        }

        err = r.Db.Model(acc).Updates(models.OpenAiAccount{
            UsdSpent: acc.UsdSpent,
            Status:   status,
        }).Error
        if err != nil {
            seelog.Errorf("update UsdSpent failed, err=%v, email=%v", err, acc.Email)
            continue
        }
    }

    if len(accs) == 0 {
        seelog.Errorf("!! RefreshUsage - no available open ai account !!")
        return
    }

    r.accIndexLock.Lock()
    r.Accounts = accs
    r.accIndexLock.Unlock()
}
