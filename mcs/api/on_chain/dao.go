package api

import (
	"github.com/filswan/go-mcs-sdk/mcs/api/common/constants"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/utils"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/web"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"
)

type Deal2PreSign struct {
	DealId              int64 `json:"deal_id"`
	SourceFileUploadCnt int   `json:"source_file_upload_cnt"`
	BatchCount          int   `json:"batch_count"`
}

func (onChainClient *OnChainClient) GetDeals2PreSign() ([]*Deal2PreSign, error) {
	apiUrl := utils.UrlJoin(onChainClient.BaseUrl, constants.API_URL_DAO_GET_DEALS_2_PRE_SIGN)

	var deals []*Deal2PreSign
	err := web.HttpGet(apiUrl, onChainClient.JwtToken, nil, &deals)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return deals, nil
}

type Deal2SignBatchInfo struct {
	BatchNo int      `json:"batch_no"`
	WCid    []string `json:"w_cid"`
}

type Deal2Sign struct {
	OfflineDealId int64                 `json:"offline_deal_id"`
	CarFileId     int64                 `json:"car_file_id"`
	DealId        int64                 `json:"deal_id"`
	BatchCount    int                   `json:"batch_count"`
	BatchSizeMax  int                   `json:"batch_size_max"`
	BatchInfo     []*Deal2SignBatchInfo `json:"batch_info"`
}

func (onChainClient *OnChainClient) GetDeals2Sign() ([]*Deal2Sign, error) {
	apiUrl := utils.UrlJoin(onChainClient.BaseUrl, constants.API_URL_DAO_GET_DEALS_2_SIGN)

	var deals []*Deal2Sign
	err := web.HttpGet(apiUrl, onChainClient.JwtToken, nil, &deals)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return deals, nil
}

func (onChainClient *OnChainClient) GetDeals2SignHash() ([]*Deal2Sign, error) {
	apiUrl := utils.UrlJoin(onChainClient.BaseUrl, constants.API_URL_DAO_GET_DEALS_2_SIGN_HASH)

	var deals []*Deal2Sign
	err := web.HttpGet(apiUrl, onChainClient.JwtToken, nil, &deals)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return deals, nil
}
