package config

import (
    "github.com/cihub/seelog"
    "github.com/minio/pkg/wildcard"
    "net/url"
)

type GuestConfig struct {
    Enabled  bool   `yaml:"enabled"`
    MaxUsage int    `yaml:"max_usage"`
    Domain   string `yaml:"domain"`
}

type GuestConfigs []GuestConfig

func (r *GuestConfigs) GetDomainConfig(domain string) *GuestConfig {
    pUrl, err := url.Parse(domain)
    if err != nil {
        seelog.Errorf("parse domain failed, domain=%s, err=%v", domain, err)
        return nil
    }
    for _, x := range *r {
        if !wildcard.Match(x.Domain, pUrl.Hostname()) {
            continue
        }
        if !x.Enabled {
            break
        }
        return &x
    }
    return nil
}
