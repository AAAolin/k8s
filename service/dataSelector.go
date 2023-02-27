package service

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	nwv1 "k8s.io/api/networking/v1"
	"sort"
	"strings"
	"time"
)

type DataSelector struct {
	GenericDataList []DataCell
	DataSelectQuery *DataSelect
}

type DataCell interface {
	GetCreateTime() time.Time
	GetName() string
}

type DataSelect struct {
	FilterQuery   *Filter
	PaginateQuery *Paginate
}

type Filter struct {
	Name string
}

type Paginate struct {
	Limit int
	Page  int
}

// 实现自定义排序，需要重写Len、Swap、Less方法

// Len 获取数组长度
func (d *DataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap 在Less方法比较结果后，定义排序规则
func (d *DataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 定义数组中元素大小的比较方式
func (d *DataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreateTime()
	b := d.GenericDataList[j].GetCreateTime()
	// 如果a的时间在b之前，返回true
	return a.After(b)
}

// Sort 排序
func (d *DataSelector) Sort() *DataSelector {
	sort.Sort(d)
	return d
}

// Filter 过滤
func (d *DataSelector) Filter() *DataSelector {
	if d.DataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	FilterDataList := make([]DataCell, 0)
	for _, value := range d.GenericDataList {
		if strings.Contains(value.GetName(), d.DataSelectQuery.FilterQuery.Name) {
			FilterDataList = append(FilterDataList, value)
			continue
		} else {
			continue
		}
	}
	d.GenericDataList = FilterDataList
	return d
}

// Paginate 分页
func (d *DataSelector) Paginate() *DataSelector {
	limit := d.DataSelectQuery.PaginateQuery.Limit
	page := d.DataSelectQuery.PaginateQuery.Page
	if page <= 0 || limit <= 0 {
		return d
	}
	//total=5，limit=10，如果page=2，那么endIndex会报超过切片长度的错误，所以要判断page值
	if len(d.GenericDataList) > 0 {
		if len(d.GenericDataList)%limit > 0 && page > len(d.GenericDataList)/limit+1 {
			page = len(d.GenericDataList)/limit + 1
		} else if len(d.GenericDataList)%limit == 0 && page > len(d.GenericDataList)/limit {
			page = len(d.GenericDataList) / limit
		} else {
			page = page
		}
	}
	// 假如 共25条数据，limit=10，那么page应该是3
	// 第一页 数据的index是0-10
	// 第二页 数据的index是10-20
	// 第三页 数据的index是20-25
	startIndex := limit * (page - 1)
	endIndex := limit * page
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
		fmt.Println(startIndex, endIndex)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// 实现 DataMethods 接口
type podCell corev1.Pod

func (d podCell) GetCreateTime() time.Time {
	return d.CreationTimestamp.Time
}

func (d podCell) GetName() string {
	return d.Name
}

type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreateTime() time.Time {
	return d.CreationTimestamp.Time
}

func (d deploymentCell) GetName() string {
	return d.Name
}

type daemonSetCell appsv1.DaemonSet

func (d daemonSetCell) GetCreateTime() time.Time {
	return d.CreationTimestamp.Time
}

func (d daemonSetCell) GetName() string {
	return d.Name
}

type statefulSetCell appsv1.StatefulSet

func (d statefulSetCell) GetCreateTime() time.Time {
	return d.CreationTimestamp.Time
}

func (d statefulSetCell) GetName() string {
	return d.Name
}

type nodeCell corev1.Node

func (n nodeCell) GetCreateTime() time.Time {
	return n.CreationTimestamp.Time
}

func (n nodeCell) GetName() string {
	return n.Name
}

type namespaceCell corev1.Namespace

func (n namespaceCell) GetCreateTime() time.Time {
	return n.CreationTimestamp.Time
}

func (n namespaceCell) GetName() string {
	return n.Name
}

type serviceCell corev1.Service

func (s serviceCell) GetCreateTime() time.Time {
	return s.CreationTimestamp.Time
}

func (s serviceCell) GetName() string {
	return s.Name
}

type ingressCell nwv1.Ingress

func (s ingressCell) GetCreateTime() time.Time {
	return s.CreationTimestamp.Time
}

func (s ingressCell) GetName() string {
	return s.Name
}
