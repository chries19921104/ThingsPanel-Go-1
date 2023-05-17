package models

type TpOta struct {
	Id                 string `json:"id" gorm:"primaryKey"`
	PackageName        string `json:"package_name,omitempty"`
	PackageVersion     string `json:"package_version,omitempty"`
	PackageModule      string `json:"package_module,omitempty"`
	ProductId          string `json:"product_id,omitempty"`
	SignatureAlgorithm string `json:"signature_algorithm,omitempty"` //签名算法
	PackageUrl         string `json:"package_url,omitempty"`
	FileSize           string `json:"file_size,omitempty"`
	Description        string `json:"description,omitempty"`
	AdditionalInfo     string `json:"additional_info,omitempty"`
	CreatedAt          int64  `json:"created_at,omitempty"`
	Sign               string `json:"sign,omitempty"`
}

func (TpOta) TableName() string {
	return "tp_ota"
}
