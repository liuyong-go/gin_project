package ycommon

import "math"

func Paginator(page, pagesize int, nums int64) map[string]interface{} {

	var prepage int  //前一页地址
	var nextpage int //后一页地址
	//根据nums总数，和prepage每页数量 生成分页总数
	totalpages := int(math.Ceil(float64(nums) / float64(pagesize))) //page总数
	if page > totalpages {
		page = totalpages
	}
	if page <= 0 {
		page = 1
	}
	var pages []int
	switch {
	case page >= totalpages-5 && totalpages > 5: //最后5页
		start := totalpages - 5 + 1
		prepage = page - 1
		nextpage = int(math.Min(float64(totalpages), float64(page+1)))
		pages = make([]int, 5)
		for i, _ := range pages {
			pages[i] = start + i
		}
	case page >= 3 && totalpages > 5:
		start := page - 3 + 1
		pages = make([]int, 5)
		prepage = page - 3
		for i, _ := range pages {
			pages[i] = start + i
		}
		prepage = page - 1
		nextpage = page + 1
	default:
		pages = make([]int, int(math.Min(5, float64(totalpages))))
		for i, _ := range pages {
			pages[i] = i + 1
		}
		prepage = int(math.Max(float64(1), float64(page-1)))
		nextpage = page + 1
		//fmt.Println(pages)
	}
	paginatorMap := make(map[string]interface{})
	paginatorMap["pages"] = pages
	paginatorMap["totalpages"] = totalpages
	paginatorMap["prepage"] = prepage
	paginatorMap["nextpage"] = nextpage
	paginatorMap["page"] = page
	return paginatorMap
}
