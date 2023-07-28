package models

type SubnetRequest struct {
	VpcCidr   string `json:"vpcCidr"`
	SubnetCnt int    `json:"subnetCnt"`
}
