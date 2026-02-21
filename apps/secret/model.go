package secret

import (
	"cmdb/apps/resource"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/infraboard/mcube/v2/crypto/cbc"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/mcube/v2/types"
)

// 泛型构造方法
func NewSecretSet() *types.Set[*Secret] {
	return types.New[*Secret]()
}

// 生成新的secret，基于vender、address、apikey实现唯一性
func NewSecret(in *CreateSecretRequest) *Secret {
	//  hash版本的UUID
	// 	Vendor Address ApiKey
	uid := uuid.NewMD5(uuid.Nil, fmt.Appendf(nil, "%d.%s.%s", in.Vendor, in.Address, in.ApiKey)).String()
	return &Secret{
		Id:                  uid,
		UpdateAt:            time.Now().Unix(),
		CreateSecretRequest: *in,
	}
}

type Secret struct {
	Id                  string `json:"id" bson:"id"`
	UpdateAt            int64  `json:"update_at" bson:"update_at"`
	CreateSecretRequest `bson:"inline"`
}

func (s *Secret) SetDefault() *Secret {
	if s.SyncLimit == 0 {
		s.SyncLimit = 10
	}
	return s
}

func (s *Secret) String() string {
	return pretty.ToJSON(s)
}

func NewCreateSecretRequest() *CreateSecretRequest {
	return &CreateSecretRequest{
		Regions:   []string{},
		SyncLimit: 10,
	}
}

type CreateSecretRequest struct {
	// 名称
	Name string `json:"name"`
	//
	Vendor resource.Vendor `json:"vendor"`
	// Vmware
	Address string `json:"address"`
	// 需要被脱敏
	// Musk
	ApiKey string `json:"api_key"`
	//
	ApiSecret string `json:"api_secret" mask:",5,4"`
	//
	isEncrypted bool

	// 资源所在区域
	Regions []string `json:"regions"`
	// 通过分页大小
	SyncLimit int64 `json:"sync_limit"`
}

func (r *CreateSecretRequest) SetIsEncrypted(v bool) {
	r.isEncrypted = v
}

func (r *CreateSecretRequest) GetSyncLimit() int64 {
	if r.SyncLimit == 0 {
		return 10
	}
	return r.SyncLimit
}

// 用 SECRET_KEY（base64解码得到密钥）对 r.ApiSecret 做 AES-CBC 加密，密文再base64编码
func (r *CreateSecretRequest) EncryptedApiSecret() error {
	if r.isEncrypted {
		return nil
	}
	// Hash, 对称，非对称
	// 对称加密 AES(cbc)
	// @v1,xxxx@xxxxx

	key, err := base64.StdEncoding.DecodeString(SECRET_KEY)
	if err != nil {
		return err
	}

	// 指定秘钥加密ApiSecret
	cipherText, err := cbc.MustNewAESCBCCihper(key).Encrypt([]byte(r.ApiSecret))
	if err != nil {
		return err
	}
	r.ApiSecret = base64.StdEncoding.EncodeToString(cipherText)
	r.SetIsEncrypted(true)
	return nil

}

// r.ApiSecret（base64解码得到密文）用 SECRET_KEY（base64解码得到密钥）做 AES-CBC 解密
func (r *CreateSecretRequest) DecryptedApiSecret() error {
	if r.isEncrypted {
		// base64解码
		cipherdText, err := base64.StdEncoding.DecodeString(r.ApiSecret)
		if err != nil {
			return err
		}

		// base64解码得到加密秘钥key
		key, err := base64.StdEncoding.DecodeString(SECRET_KEY)
		if err != nil {
			return err
		}
		// 使用key解密cipherText得到初始明文
		plainText, err := cbc.MustNewAESCBCCihper(key).Decrypt([]byte(cipherdText))
		if err != nil {
			return err
		}
		r.ApiSecret = string(plainText)
		r.SetIsEncrypted(false)
	}
	return nil
}
