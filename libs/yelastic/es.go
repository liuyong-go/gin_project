package yelastic

//文档 https://www.elastic.co/guide/cn/elasticsearch/guide/current/geohash-mapping.html
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/liuyong-go/gin_project/config"
	"github.com/tidwall/gjson"
)

//"github.com/olivere/elastic/v7"
var (
	elas *ES
	once sync.Once
	err  error
)
var booltrue bool = true

type ES struct {
	Client *es7.Client
}

//单词相关搜索，按照相关分排序
type wordMatch struct {
	Query struct {
		Match map[string]interface{} `json:"match"`
	} `json:"query"`
}
type multiFiledMatch struct {
	Query struct {
		MultiMatch struct {
			Query  string   `json:"query"`
			Type   string   `json:"type"`
			Fields []string `json:"fields"`
		} `json:"multi_match"`
	} `json:"query"`
}

//多个单词同时包含且以短语形式紧挨着
type phraseMatch struct {
	Query struct {
		MatchPhrase map[string]interface{} `json:"match_phrase"`
	} `json:"query"`
}

func NewES() *ES {
	once.Do(func() {
		esConfig := config.Config.Es
		cfg := es7.Config{
			Addresses: esConfig.Address,
		}
		//var err error
		var els = &ES{}
		els.Client, err = es7.NewClient(cfg)
		if err != nil {
			fmt.Println("elasticsearch err:", err)
		}
		elas = els
	})
	return elas
}
func (elas *ES) Delete(ctx context.Context, index string, documentType string, documentID string) ([]byte, error) {
	req := esapi.DeleteRequest{
		Index:        index,
		DocumentID:   documentID,
		DocumentType: documentType,
		Refresh:      "true",
	}
	res, err := req.Do(ctx, elas.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//CreateDocument 创建文档 index 索引，DocumentType 类型，documentID 文章id,text 内容，结构体或map
func (elas *ES) UpdateDocument(ctx context.Context, index string, documentType string, documentID string, text interface{}) ([]byte, error) {
	textDoc := map[string]interface{}{
		"doc": text,
	}
	mjson, _ := json.Marshal(textDoc)
	var content = string(mjson)
	req := esapi.UpdateRequest{
		Index:        index,
		DocumentID:   documentID,
		DocumentType: documentType,
		Body:         strings.NewReader(content),
		Refresh:      "true",
	}
	fmt.Println("req", content)
	res, err := req.Do(ctx, elas.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//CreateDocument 创建文档/或修改 index 索引，DocumentType 类型，documentID 文章id,text 内容，结构体或map

func (elas *ES) PutDocument(ctx context.Context, index string, documentType string, documentID string, text interface{}) ([]byte, error) {
	mjson, _ := json.Marshal(text)
	var content = string(mjson)
	req := esapi.IndexRequest{
		Index:        index,
		DocumentID:   documentID,
		DocumentType: documentType,
		Body:         strings.NewReader(content),
		Refresh:      "true",
	}
	res, err := req.Do(ctx, elas.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//CreateDocument 创建文档 index 索引，DocumentType 类型，documentID 文章id,text 内容，结构体或map
func (elas *ES) CreateDocument(ctx context.Context, index string, documentType string, documentID string, text interface{}) ([]byte, error) {
	mjson, _ := json.Marshal(text)
	var content = string(mjson)
	req := esapi.CreateRequest{
		Index:        index,
		DocumentID:   documentID,
		DocumentType: documentType,
		Body:         strings.NewReader(content),
		Refresh:      "true",
	}
	res, err := req.Do(ctx, elas.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//根据id拿信息
func (elas *ES) GetByID(ctx context.Context, index string, documentType string, documentID string) (result string, err error) {
	res, err := elas.Client.Get(index, documentID, elas.Client.Get.WithDocumentType(documentType), elas.Client.Get.WithContext(ctx))
	if err != nil {
		return
	}
	resultb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	result = gjson.GetBytes(resultb, "_source").String()
	return
}

type MgetData struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
	ID    string `json:"_id"`
}

//同index，同type下根据多个id拿信息，不同分区，需要拼凑docs = []MgetData
func (elas *ES) MgetByIds(ctx context.Context, index string, documentType string, documentIDs []string) ([]byte, error) {

	textDoc := map[string]interface{}{
		"ids": documentIDs,
	}
	mjson, _ := json.Marshal(textDoc)
	var content = string(mjson)
	req := esapi.MgetRequest{
		Index:        index,
		DocumentType: documentType,
		Body:         strings.NewReader(content),
	}
	res, err := req.Do(ctx, elas.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//Search 复杂检索，自定义query
func (elas *ES) Search(ctx context.Context, query interface{}, index string, documentType string, from int, size int) ([]byte, error) {
	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	requestArr := []func(*esapi.SearchRequest){
		elas.Client.Search.WithContext(ctx),
		elas.Client.Search.WithIndex(index),
		elas.Client.Search.WithBody(&buf),
		elas.Client.Search.WithTrackTotalHits(true),
		elas.Client.Search.WithPretty(),
		elas.Client.Search.WithFrom(from),
		elas.Client.Search.WithSize(size),
	}
	if documentType != "" {
		requestArr = append(requestArr, elas.Client.Search.WithDocumentType(documentType))
	}
	res, err := elas.Client.Search(requestArr...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//创建mapping
func (elas *ES) GetMapping(index string) ([]byte, error) {
	req := esapi.IndicesGetMappingRequest{
		Index:           []string{index},
		IncludeTypeName: &booltrue,
	}
	res, err := req.Do(context.Background(), elas.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//创建索引
//yelastic.NewES().CreateIndex("localcities")
func (elas *ES) CreateIndex(index string) ([]byte, error) {
	req := esapi.IndicesCreateRequest{
		Index:           index,
		IncludeTypeName: &booltrue,
	}
	res, err := req.Do(context.Background(), elas.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//生成mapping
//	mapData := map[string]interface{}{
//		"properties": map[string]interface{}{
//			"location": map[string]interface{}{
//				"type":  "geo_point",
//				"index": "true",
//			},
//			"name": map[string]interface{}{
//				"type": "keyword",
//			},
//		},
//	}
//	yelastic.NewES().PutMapping(mapData, []string{"cities/doc"})
func (elas *ES) PutMapping(content string, indexs []string) ([]byte, error) {
	// mjson, _ := json.Marshal(text)
	// var content = string(mjson)
	//直接给json
	req := esapi.IndicesPutMappingRequest{
		Index:           indexs,
		Body:            strings.NewReader(content),
		IncludeTypeName: &booltrue,
	}
	res, err := req.Do(context.Background(), elas.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//WordSearch 关键词检索 index 索引名字，key 检索字段，value 检索字段值
func (elas *ES) WordMatch(ctx context.Context, key string, value string, index string, documentType string) ([]byte, error) {
	var query wordMatch
	query.Query.Match = map[string]interface{}{key: value}
	return elas.Search(ctx, query, index, documentType, 0, 10)
}

//短语搜索
func (elas *ES) PhrasseMatch(ctx context.Context, key string, value string, index string, documentType string) ([]byte, error) {
	var query phraseMatch
	query.Query.MatchPhrase = map[string]interface{}{key: value}
	return elas.Search(ctx, query, index, documentType, 0, 10)
}

//WordMultiSearch 多字段关键词检索 index 索引名字，key 检索字段，value 检索字段值
func (elas *ES) WordMultiSearch(ctx context.Context, keyword string, fields []string, index string, documentType string) ([]byte, error) {
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"multi_match": map[string]interface{}{
	// 			"query":  keyword,
	// 			"type":   "most_fields",
	// 			"fields": fields,
	// 		},
	// 	},
	// }
	var query multiFiledMatch
	query.Query.MultiMatch.Type = "most_fields"
	query.Query.MultiMatch.Fields = fields
	query.Query.MultiMatch.Query = keyword
	return elas.Search(ctx, query, index, documentType, 0, 10)
}
func (elas *ES) SearchByLocation(ctx context.Context, distance string, lat string, lon string, index string, documentType string, page int, pageSize int) ([]byte, error) {
	offset := (page - 1) * pageSize
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": map[string]interface{}{
					"geo_distance": map[string]interface{}{
						"distance":      distance + "km",
						"distance_type": "plane",
						"location": map[string]interface{}{
							"lat": lat,
							"lon": lon,
						},
					},
				},
			},
		},
		"size": pageSize,
		"from": offset,
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": lat,
						"lon": lon,
					},
					"order": "asc",
					"unit":  "km",
				},
			},
		},
	}
	return elas.Search(ctx, query, index, documentType, 0, 10)
}
