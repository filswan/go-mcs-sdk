package mcs

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"go-mcs-sdk/mcs/common"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type McsClient struct {
	Client
	UserWalletAddressForRegisterMcs string `json:"user_wallet_address_for_register_mcs"`
	UserWalletAddressPK             string `json:"user_wallet_address_pk"`
	ChainNameForRegisterOnMcs       string `json:"chain_name_for_register_on_mcs"`
}

func NewMcsClient() *McsClient {
	mcsClient := McsClient{}
	mcsClient = *mcsClient.GetConfig()
	return &mcsClient
}

func (client *McsClient) GetConfig() *McsClient {
	err := common.LoadEnv()
	if err != nil {
		log.Fatal(err)
		return client
	}
	walletAddress := os.Getenv("USER_WALLET_ADDRESS_FOR_REGISTER_MCS")
	if walletAddress == "" {
		err = fmt.Errorf("user wallet address is null in .env file")
		log.Fatal(err)
		return client
	}
	client.UserWalletAddressForRegisterMcs = walletAddress
	walletAddressPK := os.Getenv("USER_WALLET_ADDRESS_PK")
	if walletAddressPK == "" {
		err = fmt.Errorf("user wallet address private key is null in .env file")
		log.Fatal(err)
		return client
	}
	client.UserWalletAddressPK = walletAddressPK
	chainNetworkName := os.Getenv("CHAIN_NAME_FOR_REGISTER_ON_MCS")
	if chainNetworkName == "" {
		err = fmt.Errorf("chain network name is null in .env file")
		log.Fatal(err)
		return client
	}
	client.ChainNameForRegisterOnMcs = chainNetworkName
	mcsBackendBaseUrl := os.Getenv("MCS_BACKEND_BASE_URL")
	if mcsBackendBaseUrl == "" {
		err = fmt.Errorf("mcs backend base url is null in .env file")
		log.Fatal(err)
		return client
	}
	client.BaseURL = mcsBackendBaseUrl
	return client
}

func (client *McsClient) GetToken() error {
	user, err := client.NewUserRegisterService().SetWalletAddress(client.UserWalletAddressForRegisterMcs).Do(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	nonce := user.Data.Nonce
	privateKey, _ := crypto.HexToECDSA(client.UserWalletAddressPK)
	signature, _ := common.PersonalSign(nonce, privateKey)
	jwt, err := client.NewUserLoginService().SetNetwork(client.ChainNameForRegisterOnMcs).SetNonce(nonce).SetWalletAddress(client.UserWalletAddressForRegisterMcs).
		SetSignature(signature).Do(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	client.SetJwtToken(jwt.Data.JwtToken)
	return nil
}

func (client *McsClient) UserLogin(walletAddress, signature, nonce, network string) ([]byte, error) {
	httpRequestUrl := client.BaseURL + common.USER_LOGIN
	params := make(map[string]string)
	params["public_key_address"] = walletAddress
	params["nonce"] = nonce
	params["signature"] = signature
	params["network"] = network
	response, err := common.HttpPost(httpRequestUrl, client.JwtToken, params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(*(*string)(unsafe.Pointer(&response)))
	return response, nil
}

func (client *McsClient) UserRegister(walletAddress string) (*string, error) {
	httpRequestUrl := client.BaseURL + common.USER_REGISTER
	params := make(map[string]string)
	params["public_key_address"] = walletAddress
	response, err := common.HttpPost(httpRequestUrl, client.JwtToken, params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(*(*string)(unsafe.Pointer(&response)))
	var dict map[string]interface{}
	err = json.Unmarshal(response, &dict)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(dict)
	dataInReturn := dict["data"].(map[string]interface{})
	objectInData := dataInReturn["nonce"].(string)
	fmt.Println(objectInData)
	return &objectInData, nil
}

func (client *McsClient) GetJwtToken() error {
	nonce, err := client.UserRegister(client.UserWalletAddressForRegisterMcs)
	if err != nil {
		log.Println(err)
		return err
	}
	privateKey, _ := crypto.HexToECDSA(client.UserWalletAddressPK)
	signature, _ := common.PersonalSign(*nonce, privateKey)
	resp, err := client.UserLogin(client.UserWalletAddressForRegisterMcs, signature, *nonce, client.ChainNameForRegisterOnMcs)
	if err != nil {
		log.Println(err)
		return err
	}
	var dict map[string]interface{}
	err = json.Unmarshal(resp, &dict)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(dict)
	dataInReturn := dict["data"].(map[string]interface{})
	jwtToken := dataInReturn["jwt_token"].(string)
	client.SetJwtToken(jwtToken)
	return nil
}

func (client *McsClient) GetParams() ([]byte, error) {
	httpRequestUrl := client.BaseURL + common.MCS_PARAMS
	response, err := common.HttpGet(httpRequestUrl, client.JwtToken, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(*(*string)(unsafe.Pointer(&response)))
	return response, nil
}

func (client *McsClient) GetPriceRate() ([]byte, error) {
	httpRequestUrl := client.BaseURL + common.PRICE_RATE
	response, err := common.HttpGet(httpRequestUrl, client.JwtToken, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(*(*string)(unsafe.Pointer(&response)))
	return response, nil
}

func (client *McsClient) GetPaymentInfo(sourceFileUploadId int) ([]byte, error) {
	httpRequestUrl := client.BaseURL + common.PAYMENT_INFO + strconv.Itoa(sourceFileUploadId)
	params := make(map[string]int)
	params["source_file_upload_id"] = sourceFileUploadId
	response, err := common.HttpGet(httpRequestUrl, client.JwtToken, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(*(*string)(unsafe.Pointer(&response)))
	return response, nil
}

func (client *McsClient) GetUserTasksDeals(fileName, status string, pageNumber, pageSize int) ([]byte, error) {
	requestParam := "?file_name=" + fileName + "status=" + status + "page_number=" + strconv.Itoa(pageNumber) + "page_size=" + strconv.Itoa(pageSize)
	httpRequestUrl := client.BaseURL + common.TASKS_DEALS + requestParam
	response, err := common.HttpGet(httpRequestUrl, client.JwtToken, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(*(*string)(unsafe.Pointer(&response)))
	return response, nil
}
