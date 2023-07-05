package sendmqtt

import (
	"errors"
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/yaml.v3"
)

// type ShopContent struct {
// 	Name    string `json:"Name"`
// 	Address string `json:"Address"`
// 	Number  string `json:"Number"`
// }

// type PotType struct {
// 	Name         string `json:"Name"`
// 	SoupStandard int    `json:"SoupStandard"`
// 	PotTypeId    string `json:"PotTypeId"`
// }

// type Taste struct {
// 	Name           string   `json:"Name"`
// 	TasteId        string   `json:"TasteId"`
// 	PotTypeId      string   `json:"PotTypeId"`
// 	MaterialIdList []string `json:"MaterialIdList"`
// 	//RecipeId       string   `json:"RecipeId"`
// 	BottomPotId string `json:"BottomPotId"`
// }

// type Materials struct {
// 	Id        string `json:"Id"`
// 	Name      string `json:"Name"`
// 	Dosage    int    `json:"Dosage"`
// 	Unit      string `json:"Unit"`
// 	WaterLine int    `json:"WaterLine"`
// 	Station   string `json:"Station"`
// 	Resource  string `json:"Resource"`
// }

// type SendConfig struct {
// 	Shop      ShopContent
// 	PotType   []*PotType
// 	Taste     []*Taste
// 	Materials []*Materials
// 	Recipe    []*Recipe
// }

// type Recipe struct {
// 	BottomPotId      string `json:"BottomPotId"`
// 	BottomPot        string `json:"BottomPot"`
// 	PotTypeId        string `json:"PotTypeId"`
// 	BottomProperties string `json:"BottomProperties"`
// 	//SoupStandard     int64     `json:"SoupStandard"`
// 	MaterialIdList []string `json:"MaterialIdList"`
// }

// type HdlOrder struct {
// 	StoreID                    string `json:"storeId"`
// 	OrderID                    string `json:"orderId"`
// 	PotID                      string `json:"potId"`
// 	TableNumber                string `json:"tableNumber"`
// 	OrderTime                  string `json:"orderTime"`
// 	SoupAddingStartTime        string `json:"soupAddingStartTime"`
// 	SoupAddingFinishTime       string `json:"soupAddingFinishTime"`
// 	IngredientAddingStartTime  string `json:"ingredientAddingStartTime"`
// 	IngredientAddingFinishTime string `json:"ingredientAddingFinishTime"`
// 	PotSwitchingFinishTime     string `json:"potSwitchingFinishTime"`
// 	CreationTime               string `json:"creationTime"`
// }
// type PotTypeConfig struct {
// 	SendId  string
// 	PotType []*PotType
// }
// type TasteConfig struct {
// 	SendId string
// 	Taste  []*Taste
// }
// type MaterialsConfig struct {
// 	SendId    string
// 	Materials []*Materials
// }
// type RecipeConfig struct {
// 	SendId string
// 	Recipe []*Recipe
// }

/*店铺配置信息json案例,potType锅型,taste口味,materials配料,recipe锅底配方,shop店铺
{
	"shop": {
		"name":"海底捞A店",
		"address":"北京市海淀区中关村大街1号",
		"number":"A001"
	}
	"potType": [{
		"id": "1",
		"name": "锅底1"
	}],
	"taste": [{
		"id": "1",
		"name": "口味1", // 口味名称
		"tasteId":"", // POS口味Id
		"dosage": 1, // 用量
		"unit": "g", // 单位
		"station": "鲜料工位", // 工位
		"materialsName": "物料1", // 物料名称
		"waterLine": 1 // 加汤水位标准
	}],
	"materials": [{
		"id": "1",
		"name": "配料1", // 配料名称
		"dosage": 1, // 用量
		"unit": "g", // 单位
		"station": "鲜料工位" // 工位
	}],
	"recipe": [{
		"id": "1",
		"bottomPotId": "1", // 锅底Id
		"name": "锅底配方1", // 配方名称
		"potTypeId": "1", // 锅型Id
		"tasteId": [
			"1",
			"2"
			], // 口味Id
		"materialsId": [
			"1",
			"2"
			], // 配料Id
		"waterLine": 1, // 加汤水位标准
		"bottomProperties": "1" // 锅底属性
	}],
}
*/

// 向MQTT服务器发送消息
func Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if _client == nil {
		return errors.New("_client is error")
	}
	t := _client.Publish(topic, qos, retained, payload)
	if t.Error() != nil {
		fmt.Println(t.Error())
	}
	return t.Error()
}

// 向海底捞发送mqtt消息
func SendToHDL(payload []byte, token string) (err error) {
	// 读取配置文件
	hdlConfig, err := ParseYaml()
	if err != nil {
		return err
	}
	// 发送消息
	err = Publish("device/attributes"+"/"+token, byte(hdlConfig.Qos), false, string(payload))
	if err != nil {
		logs.Error("publish error: %v", err)
		return err
	}
	return nil
}

// 向海底捞发送mqtt消息
func SendTimeToHDL(payload []byte, token string) (err error) {
	// 读取配置文件
	hdlConfig, err := ParseYaml()
	if err != nil {
		return err
	}
	// 发送消息
	err = Publish("device/command/"+token, byte(hdlConfig.Qos), false, string(payload))
	if err != nil {
		logs.Error("publish error: %v", err)
		return err
	}
	return nil
}

type HDLConfig struct {
	TopicToPublish string `yaml:"topicToPublish"`
	Qos            int    `yaml:"qos"`
}

// 通过gopkg.in/yaml.v3解析HDLConfig.yaml文件
func ParseYaml() (HDLConfig, error) {
	var hdlConfig HDLConfig
	// 读取配置文件
	file, err := os.Open("modules/dataService/HDLConfig.yaml")
	if err != nil {
		logs.Error("open file error: %v", err)
		return hdlConfig, err
	}
	defer file.Close()

	// 解析配置文件
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&hdlConfig)
	if err != nil {
		logs.Error("decode file error: %v", err)
		return hdlConfig, err
	}
	return hdlConfig, nil
}
