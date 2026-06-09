package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"simplest_script/core/conf"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	deleteElast "github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

var elasticsearchClient *elasticsearch.TypedClient

func newElasticsearchClient() *elasticsearch.TypedClient {
	if elasticsearchClient != nil {
		return elasticsearchClient
	}

	// 创建官方类型化客户端
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: strings.Split(conf.Conf.Elastic.Addr, ";"),
		Username:  conf.Conf.Elastic.Username,
		Password:  conf.Conf.Elastic.Password,
	})
	if err != nil {
		hlog.Info("elasticsearch 创建客户端失败")
		return nil
	}

	elasticsearchClient = client

	return client
}

type Elasticsearch struct {
	Client *elasticsearch.TypedClient
}

func NewElasticsearch() *Elasticsearch {
	return &Elasticsearch{
		Client: newElasticsearchClient(),
	}
}

func (l *Elasticsearch) Create(ctx context.Context, index, id string, doc any) (*index.Response, error) {
	return l.Client.Index(index).Id(id).Document(doc).Do(ctx)
}

func (l *Elasticsearch) Delete(ctx context.Context, index, id string) (*deleteElast.Response, error) {
	return l.Client.Delete(index, id).Do(ctx)
}

func (l *Elasticsearch) Get(ctx context.Context, index, id string) (*get.Response, error) {
	return l.Client.Get(index, id).Do(ctx)
}

func (l *Elasticsearch) Update(ctx context.Context, index, id string, doc any) (*update.Response, error) {
	return l.Client.Update(index, id).Doc(doc).Do(ctx)
}

func (l *Elasticsearch) MultiCreate(ctx context.Context, index, id string, docs []any) (*bulk.Response, error) {
	multi := l.Client.Bulk().Index(index)

	for _, v := range docs {
		multi.CreateOp(types.CreateOperation{}, v)
	}

	return multi.Do(ctx)
}

// PageResult 分页查询返回结果
type PageResult struct {
	Total    int64             // 总命中数
	Data     []json.RawMessage // 当前页的文档源数据
	Page     int               // 当前页码
	PageSize int               // 每页大小
}

func (l *Elasticsearch) Search(ctx context.Context, index string, query *types.Query, page, pageSize int, sort ...types.SortCombinations) (*PageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	from := (page - 1) * pageSize

	searchReq := l.Client.Search().Index(index).Query(query).From(from).Size(pageSize)
	if len(sort) > 0 {
		searchReq = searchReq.Sort(sort...)
	}

	resp, err := searchReq.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("execute search failed: %w", err)
	}

	// 提取源数据
	hits := resp.Hits.Hits
	data := make([]json.RawMessage, 0, len(hits))
	for _, hit := range hits {
		data = append(data, hit.Source_)
	}

	return &PageResult{
		Total:    resp.Hits.Total.Value,
		Data:     data,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
