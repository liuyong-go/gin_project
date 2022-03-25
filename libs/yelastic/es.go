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
)

var (
	elas *ES
	once sync.Once
	err  error
)
var booltrue bool = true

type ES struct {
	Client *es7.Client
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

//CreateDocument 创建文档
func (elas *ES) CreateDocument(index string, DocumentType string, documentID string, text map[string]interface{}) ([]byte, error) {
	mjson, _ := json.Marshal(text)
	var content = string(mjson)
	req := esapi.IndexRequest{
		Index:        index,
		DocumentID:   documentID,
		DocumentType: DocumentType,
		Body:         strings.NewReader(content),
		Refresh:      "true",
	}
	res, err := req.Do(context.Background(), elas.Client)
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
func (elas *ES) PutMapping(text map[string]interface{}, indexs []string) ([]byte, error) {
	mjson, _ := json.Marshal(text)
	var content = string(mjson)
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

//Search 复杂检索，自定义query
func (elas *ES) Search(query interface{}, index string, DocumentType string) ([]byte, error) {
	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := elas.Client.Search(
		elas.Client.Search.WithContext(context.Background()),
		elas.Client.Search.WithIndex(index),
		elas.Client.Search.WithDocumentType(DocumentType),
		elas.Client.Search.WithBody(&buf),
		elas.Client.Search.WithTrackTotalHits(true),
		elas.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//WordSearch 关键词检索 index 索引名字，key 检索字段，value 检索字段值
func (elas *ES) WordSearch(key string, value string, index string, documentType string) ([]byte, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				key: value,
			},
		},
	}
	return elas.Search(query, index, documentType)
}

//WordMultiSearch 多字段关键词检索 index 索引名字，key 检索字段，value 检索字段值
func (elas *ES) WordMultiSearch(keyword string, fields []string, index string, documentType string) ([]byte, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"type":   "most_fields",
				"fields": fields,
			},
		},
	}
	return elas.Search(query, index, documentType)
}
func (elas *ES) SearchByLocation(distance string, lat string, lon string, index string, documentType string, page int, pageSize int) ([]byte, error) {
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
	return elas.Search(query, index, documentType)
}
