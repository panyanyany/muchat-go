package api_base

import (
    "encoding/json"
    "errors"
    "fmt"
)

var (
    ErrNoBalance = errors.New("ErrNoBalance")
    ErrBaned     = errors.New("ErrBaned")
    ErrUnknown   = errors.New("ErrUnknown")
    ErrApiKey    = errors.New("ErrApiKey")
)

func ExtractError(body []byte) (err error) {
    errResp := ErrorRespBody{}
    err = json.Unmarshal(body, &errResp)
    if err != nil {
        err = fmt.Errorf("json.Unmarshal(body, &errResp) failed, err=%w", err)
        return
    }
    if errResp.Error.Type == ErrorTypeInsufficientQuota {
        err = ErrNoBalance
        return
    }
    if errResp.Error.Type == ErrorTypeInvalidRequest {
        if errResp.Error.Param == "model" || errResp.Error.Code == ErrorTypeAccountDeactivated {
            err = ErrBaned
            return
        }
        if errResp.Error.Code == ErrorTypeInvalidApiKey {
            err = ErrApiKey
            return
        }
    }
    if errResp.Error.Type == ErrorTypeBillingNotActive {
        err = ErrBaned
        return
    }
    if errResp.Error.Type == ErrorTypeAccessTerminated {
        err = ErrBaned
        return
    }
    err = ErrUnknown
    //err = errors.New(fmt.Sprintf("status=%d, %v", resp.StatusCode, string(body)))
    return
}
