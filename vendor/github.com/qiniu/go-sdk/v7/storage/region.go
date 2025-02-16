package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
	"golang.org/x/sync/singleflight"
)

// 存储所在的地区，例如华东，华南，华北
// 每个存储区域可能有多个机房信息，每个机房可能有多个上传入口
type Region struct {
	// 上传入口
	SrcUpHosts []string `json:"src_up,omitempty"`

	// 加速上传入口
	CdnUpHosts []string `json:"cdn_up,omitempty"`

	// 获取文件信息入口
	RsHost string `json:"rs,omitempty"`

	// bucket列举入口
	RsfHost string `json:"rsf,omitempty"`

	ApiHost string `json:"api,omitempty"`

	// 存储io 入口
	IovipHost string `json:"io,omitempty"`
}

type RegionID string

// GetDefaultReion 根据RegionID获取对应的Region信息
func GetRegionByID(regionID RegionID) (Region, bool) {
	if r, ok := regionMap[regionID]; ok {
		return r, ok
	}
	return Region{}, false
}

func (r *Region) String() string {
	str := ""
	str += fmt.Sprintf("SrcUpHosts: %v\n", r.SrcUpHosts)
	str += fmt.Sprintf("CdnUpHosts: %v\n", r.CdnUpHosts)
	str += fmt.Sprintf("IovipHost: %s\n", r.IovipHost)
	str += fmt.Sprintf("RsHost: %s\n", r.RsHost)
	str += fmt.Sprintf("RsfHost: %s\n", r.RsfHost)
	str += fmt.Sprintf("ApiHost: %s\n", r.ApiHost)
	return str
}

func endpoint(useHttps bool, host string) string {
	host = strings.TrimSpace(host)
	host = strings.TrimLeft(host, "http://")
	host = strings.TrimLeft(host, "https://")
	if host == "" {
		return ""
	}
	scheme := "http://"
	if useHttps {
		scheme = "https://"
	}
	return fmt.Sprintf("%s%s", scheme, host)
}

// 获取rsfHost
func (r *Region) GetRsfHost(useHttps bool) string {
	return endpoint(useHttps, r.RsfHost)
}

// 获取io host
func (r *Region) GetIoHost(useHttps bool) string {
	return endpoint(useHttps, r.IovipHost)
}

// 获取RsHost
func (r *Region) GetRsHost(useHttps bool) string {
	return endpoint(useHttps, r.RsHost)
}

// 获取api host
func (r *Region) GetApiHost(useHttps bool) string {
	return endpoint(useHttps, r.ApiHost)
}

var (
	// regionHuadong 表示华东机房
	regionHuadong = Region{
		SrcUpHosts: []string{
			"up.qiniup.com",
			"up-nb.qiniup.com",
			"up-xs.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload.qiniup.com",
			"upload-nb.qiniup.com",
			"upload-xs.qiniup.com",
		},
		RsHost:    "rs.qbox.me",
		RsfHost:   "rsf.qbox.me",
		ApiHost:   "api.qiniu.com",
		IovipHost: "iovip.qbox.me",
	}

	// regionHuadongZhejiang 表示华东-浙江2
	regionHuadongZhejiang = Region{
		SrcUpHosts: []string{
			"up-cn-east-2.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-cn-east-2.qiniup.com",
		},
		RsHost:    "rs-cn-east-2.qiniuapi.com",
		RsfHost:   "rsf-cn-east-2-qiniuapi.com",
		ApiHost:   "api-cn-east-2.qiniuapi.com",
		IovipHost: "iovip-cn-east-2.qiniuio.com",
	}

	// regionHuabei 表示华北机房
	regionHuabei = Region{
		SrcUpHosts: []string{
			"up-z1.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-z1.qiniup.com",
		},
		RsHost:    "rs-z1.qbox.me",
		RsfHost:   "rsf-z1.qbox.me",
		ApiHost:   "api-z1.qiniuapi.com",
		IovipHost: "iovip-z1.qbox.me",
	}
	// regionHuanan 表示华南机房
	regionHuanan = Region{
		SrcUpHosts: []string{
			"up-z2.qiniup.com",
			"up-gz.qiniup.com",
			"up-fs.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-z2.qiniup.com",
			"upload-gz.qiniup.com",
			"upload-fs.qiniup.com",
		},
		RsHost:    "rs-z2.qbox.me",
		RsfHost:   "rsf-z2.qbox.me",
		ApiHost:   "api-z2.qiniuapi.com",
		IovipHost: "iovip-z2.qbox.me",
	}

	// regionNorthAmerica 表示北美机房
	regionNorthAmerica = Region{
		SrcUpHosts: []string{
			"up-na0.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-na0.qiniup.com",
		},
		RsHost:    "rs-na0.qbox.me",
		RsfHost:   "rsf-na0.qbox.me",
		ApiHost:   "api-na0.qiniuapi.com",
		IovipHost: "iovip-na0.qbox.me",
	}
	// regionSingapore 表示新加坡机房
	regionSingapore = Region{
		SrcUpHosts: []string{
			"up-as0.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-as0.qiniup.com",
		},
		RsHost:    "rs-as0.qbox.me",
		RsfHost:   "rsf-as0.qbox.me",
		ApiHost:   "api-as0.qiniuapi.com",
		IovipHost: "iovip-as0.qbox.me",
	}
	// regionFogCnEast1 表示雾存储华东区
	regionFogCnEast1 = Region{
		SrcUpHosts: []string{
			"up-fog-cn-east-1.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-fog-cn-east-1.qiniup.com",
		},
		RsHost:    "rs-fog-cn-east-1.qbox.me",
		RsfHost:   "rsf-fog-cn-east-1.qbox.me",
		ApiHost:   "api-fog-cn-east-1.qiniuapi.com",
		IovipHost: "iovip-fog-cn-east-1.qbox.me",
	}
)

const (
	// region code
	RIDHuadong         = RegionID("z0")
	RIDHuadongZheJiang = RegionID("cn-east-2")
	RIDHuabei          = RegionID("z1")
	RIDHuanan          = RegionID("z2")
	RIDNorthAmerica    = RegionID("na0")
	RIDSingapore       = RegionID("as0")
	RIDFogCnEast1      = RegionID("fog-cn-east-1")
)

// regionMap 是RegionID到具体的Region的映射
var regionMap = map[RegionID]Region{
	RIDHuadong:         regionHuadong,
	RIDHuadongZheJiang: regionHuadongZhejiang,
	RIDHuanan:          regionHuanan,
	RIDHuabei:          regionHuabei,
	RIDSingapore:       regionSingapore,
	RIDNorthAmerica:    regionNorthAmerica,
	RIDFogCnEast1:      regionFogCnEast1,
}

/// UcHost 为查询空间相关域名的API服务地址
/// 设置 UcHost 时，如果不指定 scheme 默认会使用 https
/// UcHost 已废弃，建议使用 SetUcHost
//Deprecated
var UcHost = "https://uc.qbox.me"

var ucHost = ""

func SetUcHost(host string, useHttps bool) {
	ucHost = endpoint(useHttps, host)
}

func getUcHostByDefaultProtocol() string {
	return getUcHost(true)
}

func getUcHost(useHttps bool) string {
	if ucHost != "" {
		return ucHost
	}

	if strings.Contains(UcHost, "://") {
		return UcHost
	} else {
		return endpoint(useHttps, UcHost)
	}
}

// UcQueryRet 为查询请求的回复
type UcQueryRet struct {
	TTL    int `json:"ttl"`
	Io     map[string]map[string][]string
	IoInfo map[string]UcQueryIo `json:"io"`
	Up     map[string]UcQueryUp `json:"up"`
}

func (uc *UcQueryRet) UnmarshalJSON(data []byte) error {
	t := struct {
		TTL    int                  `json:"ttl"`
		IoInfo map[string]UcQueryIo `json:"io"`
		Up     map[string]UcQueryUp `json:"up"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	uc.TTL = t.TTL
	uc.IoInfo = t.IoInfo
	uc.Up = t.Up
	uc.setup()
	return nil
}

func (uc *UcQueryRet) setup() {
	if uc.Io != nil || uc.IoInfo == nil {
		return
	}

	uc.Io = make(map[string]map[string][]string)
	ioSrc := uc.IoInfo["src"].toMapWithoutInfo()
	if ioSrc != nil && len(ioSrc) > 0 {
		uc.Io["src"] = ioSrc
	}

	ioOldSrc := uc.IoInfo["old_src"].toMapWithoutInfo()
	if ioOldSrc != nil && len(ioOldSrc) > 0 {
		uc.Io["old_src"] = ioOldSrc
	}
}

// UcQueryUp 为查询请求回复中的上传域名信息
type UcQueryUp struct {
	Main   []string `json:"main,omitempty"`
	Backup []string `json:"backup,omitempty"`
	Info   string   `json:"info,omitempty"`
}

// UcQueryIo 为查询请求回复中的上传域名信息
type UcQueryIo struct {
	Main   []string `json:"main,omitempty"`
	Backup []string `json:"backup,omitempty"`
	Info   string   `json:"info,omitempty"`
}

func (io UcQueryIo) toMapWithoutInfo() map[string][]string {

	ret := make(map[string][]string)
	if io.Main != nil && len(io.Main) > 0 {
		ret["main"] = io.Main
	}

	if io.Backup != nil && len(io.Backup) > 0 {
		ret["backup"] = io.Backup
	}

	return ret
}

type regionCacheValue struct {
	Region   *Region   `json:"region"`
	Deadline time.Time `json:"deadline"`
}

type regionCacheMap map[string]regionCacheValue

var (
	regionCachePath     = filepath.Join(os.TempDir(), "qiniu-golang-sdk", "query.cache.json")
	regionCache         sync.Map
	regionCacheLock     sync.RWMutex
	regionCacheSyncLock sync.Mutex
	regionCacheGroup    singleflight.Group
	regionCacheLoaded   bool = false
)

func SetRegionCachePath(newPath string) {
	regionCacheLock.Lock()
	defer regionCacheLock.Unlock()

	regionCachePath = newPath
	regionCacheLoaded = false
}

func loadRegionCache() {
	cacheFile, err := os.Open(regionCachePath)
	if err != nil {
		return
	}
	defer cacheFile.Close()

	var cacheMap regionCacheMap
	if err = json.NewDecoder(cacheFile).Decode(&cacheMap); err != nil {
		return
	}
	for cacheKey, cacheValue := range cacheMap {
		regionCache.Store(cacheKey, cacheValue)
	}
}

func storeRegionCache() {
	err := os.MkdirAll(filepath.Dir(regionCachePath), 0700)
	if err != nil {
		return
	}

	cacheFile, err := os.OpenFile(regionCachePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer cacheFile.Close()

	cacheMap := make(regionCacheMap)
	regionCache.Range(func(cacheKey, cacheValue interface{}) bool {
		cacheMap[cacheKey.(string)] = cacheValue.(regionCacheValue)
		return true
	})
	if err = json.NewEncoder(cacheFile).Encode(cacheMap); err != nil {
		return
	}
}

// GetRegion 用来根据ak和bucket来获取空间相关的机房信息
func GetRegion(ak, bucket string) (*Region, error) {
	regionCacheLock.RLock()
	if regionCacheLoaded {
		regionCacheLock.RUnlock()
	} else {
		regionCacheLock.RUnlock()
		func() {
			regionCacheLock.Lock()
			defer regionCacheLock.Unlock()

			if !regionCacheLoaded {
				loadRegionCache()
				regionCacheLoaded = true
			}
		}()
	}

	regionID := fmt.Sprintf("%s:%s", ak, bucket)
	//check from cache
	if v, ok := regionCache.Load(regionID); ok && time.Now().Before(v.(regionCacheValue).Deadline) {
		return v.(regionCacheValue).Region, nil
	}

	newRegion, err, _ := regionCacheGroup.Do(regionID, func() (interface{}, error) {
		reqURL := fmt.Sprintf("%s/v2/query?ak=%s&bucket=%s", getUcHostByDefaultProtocol(), ak, bucket)

		var ret UcQueryRet
		err := client.DefaultClient.CallWithForm(context.Background(), &ret, "GET", reqURL, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("query region error, %s", err.Error())
		}

		if len(ret.Io["src"]["main"]) <= 0 {
			return nil, fmt.Errorf("empty io host list")
		}

		ioHost := ret.Io["src"]["main"][0]
		srcUpHosts := ret.Up["src"].Main
		if ret.Up["src"].Backup != nil {
			srcUpHosts = append(srcUpHosts, ret.Up["src"].Backup...)
		}
		cdnUpHosts := ret.Up["acc"].Main
		if ret.Up["acc"].Backup != nil {
			cdnUpHosts = append(cdnUpHosts, ret.Up["acc"].Backup...)
		}

		region := &Region{
			SrcUpHosts: srcUpHosts,
			CdnUpHosts: cdnUpHosts,
			IovipHost:  ioHost,
			RsHost:     DefaultRsHost,
			RsfHost:    DefaultRsfHost,
			ApiHost:    DefaultAPIHost,
		}

		//set specific hosts if possible
		setSpecificHosts(ioHost, region)
		regionCache.Store(regionID, regionCacheValue{
			Region:   region,
			Deadline: time.Now().Add(time.Duration(ret.TTL) * time.Second),
		})

		regionCacheSyncLock.Lock()
		defer regionCacheSyncLock.Unlock()

		storeRegionCache()
		return region, nil
	})

	if err != nil {
		return nil, err
	} else {
		return newRegion.(*Region), nil
	}
}

type ucRegionsRet struct {
	Regions []ucRegionRet `json:"regions"`
}

type ucRegionRet struct {
	Id  string       `json:"id"`
	Up  ucDomainsRet `json:"up"`
	Io  ucDomainsRet `json:"io"`
	Uc  ucDomainsRet `json:"uc"`
	Rs  ucDomainsRet `json:"rs"`
	Rsf ucDomainsRet `json:"rsf"`
	Api ucDomainsRet `json:"api"`
}

type ucDomainsRet struct {
	Main   []string `json:"domains"`
	Backup []string `json:"old,omitempty"`
}

var (
	regionIdCache      = make(map[string]*Region)
	regionIdCacheGroup singleflight.Group
)

// getRegionByRegionId 用来根据 ak 和 sk 来获取空间相关的机房信息，由于返回的 Region 结构体与先前不兼容，所以本接口暂不开放
func getRegionByRegionId(regionId string, credentials *auth.Credentials) (region *Region, err error) {
	var cacheValue interface{}

	cacheValue, err, _ = regionIdCacheGroup.Do("query", func() (interface{}, error) {
		if v, ok := regionIdCache[regionId]; ok {
			return v, nil
		}
		reqURL := fmt.Sprintf("%s/regions", getUcHostByDefaultProtocol())
		var ret ucRegionsRet
		ctx := context.TODO()
		qErr := client.DefaultClient.CredentialedCallWithForm(ctx, credentials, auth.TokenQiniu, &ret, "GET", reqURL, nil, nil)
		if qErr != nil {
			return nil, fmt.Errorf("query region error, %s", qErr.Error())
		}
		for _, r := range ret.Regions {
			upHosts := r.Up.Main
			if len(r.Up.Backup) > 0 {
				upHosts = append(upHosts, r.Up.Backup...)
			}
			if len(r.Io.Main) == 0 {
				return nil, fmt.Errorf("empty io host list")
			}
			if len(r.Rs.Main) == 0 {
				return nil, fmt.Errorf("empty rs host list")
			}
			if len(r.Rsf.Main) == 0 {
				return nil, fmt.Errorf("empty rsf host list")
			}
			if len(r.Api.Main) == 0 {
				return nil, fmt.Errorf("empty api host list")
			}

			region := &Region{
				SrcUpHosts: upHosts,
				CdnUpHosts: upHosts,
				IovipHost:  r.Io.Main[0],
				RsHost:     r.Rs.Main[0],
				RsfHost:    r.Rsf.Main[0],
				ApiHost:    r.Api.Main[0],
			}
			regionIdCache[r.Id] = region
		}
		if v, ok := regionIdCache[regionId]; ok {
			return v, nil
		} else {
			return nil, fmt.Errorf("region id is not found")
		}
	})
	if err != nil {
		return nil, err
	}
	region = cacheValue.(*Region)
	return
}

func regionFromHost(ioHost string) (Region, bool) {
	if strings.Contains(ioHost, "-z1") {
		return GetRegionByID(RIDHuabei)
	}
	if strings.Contains(ioHost, "-z2") {
		return GetRegionByID(RIDHuanan)
	}

	if strings.Contains(ioHost, "-na0") {
		return GetRegionByID(RIDNorthAmerica)
	}
	if strings.Contains(ioHost, "-as0") {
		return GetRegionByID(RIDSingapore)
	}
	return Region{}, false
}

func setSpecificHosts(ioHost string, region *Region) {
	r, ok := regionFromHost(ioHost)
	if ok {
		region.RsHost = r.RsHost
		region.RsfHost = r.RsfHost
		region.ApiHost = r.ApiHost
	}
}

type RegionInfo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func GetRegionsInfo(mac *auth.Credentials) ([]RegionInfo, error) {
	var regions struct {
		Regions []RegionInfo `json:"regions"`
	}
	qErr := client.DefaultClient.CredentialedCallWithForm(context.Background(), mac, auth.TokenQiniu, &regions, "GET", getUcHostByDefaultProtocol()+"/regions", nil, nil)
	if qErr != nil {
		return nil, fmt.Errorf("query region error, %s", qErr.Error())
	} else {
		return regions.Regions, nil
	}
}
