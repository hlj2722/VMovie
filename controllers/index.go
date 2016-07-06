package controllers

import (
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/zituocn/VMovie/models"
	"strconv"
	"strings"
	//"time"
)

///前台页面handel
type IndexHandel struct {
	baseController
}

func (this *IndexHandel) Start() {
	this.TplName = "_start.html"
}

///今日更新
func (this *IndexHandel) Today() {
	var (
		info models.MovieInfo
		list []*models.MovieInfo
	)
	this.Ctx.Output.Header("Cache-Control", "public")

	list = info.GetWeekList(100)

	///内页热门列表
	hostlist := info.GetHotList(0, 10)
	this.Data["hostlist"] = hostlist

	//随机列表
	randlist := info.GetRandList(10)
	this.Data["randlist"] = randlist

	this.Data["list"] = list

	this.Data["week"] = this.GetWeekString()
	this.TplName = "_today.html"
}

///专题页面详情
func (this *IndexHandel) Page() {
	var (
		minfo models.MovieInfo
		info  = new(models.PageInfo)
		ename string
		err   error
	)
	//页面cache控制
	this.Ctx.Output.Header("Cache-Control", "public")
	ename = strings.TrimSpace(this.Ctx.Input.Param(":ename"))
	if len(ename) == 0 {
		this.Abort("404")
	}

	//直接orm查询
	o := orm.NewOrm()
	o.Using("default")
	err = o.QueryTable("page_info").Filter("ename", ename).One(info)

	if err != nil {
		this.Abort("404")
	}
	///内页热门列表
	hostlist := minfo.GetHotList(0, 10)
	this.Data["hostlist"] = hostlist

	//随机列表
	randlist := minfo.GetRandList(10)
	this.Data["randlist"] = randlist

	this.Data["info"] = info
	this.TplName = "_pageinfo.html"
}

//search页面
func (this *IndexHandel) Search() {
	var (
		keyword  string
		page     int64
		pagesize int64 = 40
		offset   int64
		info     models.MovieInfo
		list     []*models.MovieInfo
		pager    string
		nodate   int64
	)
	this.Ctx.Output.Header("Cache-Control", "public")
	keyword = strings.TrimSpace(this.Ctx.Input.Param(":key"))
	pagestr := this.Ctx.Input.Param(":page")
	page, _ = strconv.ParseInt(pagestr, 10, 64)
	if page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize

	query := info.Query()
	query = query.Filter("status__gte", 0)
	query = query.Filter("title__icontains", keyword).OrderBy("-id")

	count, _ := query.Count()
	if count > 0 {
		query.Limit(pagesize, offset).All(&list, "id", "name", "title", "ename", "photo", "Hasepisode", "Episode")
	} else {
		nodate = 1
	}

	pager = this.PageList(pagesize, page, count, false, "/search/"+keyword)
	this.Data["pager"] = pager
	this.Data["list"] = list

	///内页热门列表
	hostlist := info.GetHotList(0, 10)
	this.Data["hostlist"] = hostlist

	//随机列表
	randlist := info.GetRandList(10)
	this.Data["randlist"] = randlist

	this.Data["nodate"] = nodate

	this.TplName = "_search.html"
}

//Json输出页面
func (this *IndexHandel) Json() {
	this.Ctx.Output.Header("Cache-Control", "public")

	var (
		info models.MovieInfo
		list []*models.MovieInfo
	)
	info.Query().OrderBy("-id").Limit(30, 0).All(&list)
	this.Data["json"] = &list
	this.ServeJSON()
}

///前台首页
func (this *IndexHandel) Index() {
	var (
		info   models.MovieInfo
		mphoto []*models.MovieInfo
	)
	this.Ctx.Output.Header("Cache-Control", "public")
	list100 := info.GetList(100, 14)
	list200 := info.GetList(200, 14)
	list1 := info.GetList(1, 7)
	list2 := info.GetList(2, 7)
	list3 := info.GetList(3, 7)
	list4 := info.GetList(4, 7)
	list5 := info.GetList(5, 7)
	list6 := info.GetList(6, 7)
	wlist := info.GetWeekList(6)

	info.Query().Filter("status", 1).Limit(6, 0).OrderBy("-Updatetime").All(&mphoto, "id", "name", "ename", "iphoto")

	this.Data["list100"] = list100
	this.Data["list200"] = list200
	this.Data["list1"] = list1
	this.Data["list2"] = list2
	this.Data["list3"] = list3
	this.Data["list4"] = list4
	this.Data["list5"] = list5
	this.Data["list6"] = list6
	this.Data["wlist"] = wlist
	this.Data["mphoto"] = mphoto
	this.Data["week"] = this.GetWeekString()
	this.TplName = "_index.html"
}

///前台分类列表页
func (this *IndexHandel) List() {
	var (
		cid      int64
		page     int64
		pagesize int64 = 40
		offset   int64
		info     models.MovieInfo
		list     []*models.MovieInfo
		cinfo    *models.MovieClassInfo = new(models.MovieClassInfo)
		pager    string
		err      error
	)
	this.Ctx.Output.Header("Cache-Control", "public")

	cidstr := this.Ctx.Input.Param(":cid")
	cid, err = strconv.ParseInt(cidstr, 10, 64)
	if err != nil || cid <= 0 {
		this.Abort("404")
	}

	//查询分类信息
	cinfo.Id = cid
	err = cinfo.Read()
	if err != nil {
		this.Abort("404")
	}

	pagestr := this.Ctx.Input.Param(":page")
	page, _ = strconv.ParseInt(pagestr, 10, 64)
	if page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize

	query := info.Query()
	query = query.Filter("status__gte", 0)
	if cid > 0 && cid != 100 && cid != 200 {
		query = query.Filter("cid", cid).OrderBy("-id")
	} else if cid == 100 {
		query = query.OrderBy("-id")
	} else if cid == 200 {
		query = query.OrderBy("-views")
	}

	count, _ := query.Count()
	if count > 0 {
		query.Limit(pagesize, offset).All(&list, "id", "name", "title", "ename", "photo", "Hasepisode", "Episode")
	}

	pager = this.PageList(pagesize, page, count, false, "/m/"+cidstr)
	this.Data["pager"] = pager
	this.Data["list"] = list
	this.Data["cid"] = cid
	this.Data["cinfo"] = cinfo

	///内页热门列表
	hostlist := info.GetHotList(cinfo.Id, 10)
	this.Data["hostlist"] = hostlist

	//随机列表
	randlist := info.GetRandList(10)
	this.Data["randlist"] = randlist

	this.TplName = "_list.html"
}

//前台详细页
func (this *IndexHandel) Detail() {
	var (
		id     int64
		info   *models.MovieInfo   = new(models.MovieInfo)
		rmlist []*models.MovieInfo //相关影片数据
		rinfo  models.RelationInfo //影片关系
		down   models.DownAddrInfo
		//downlist string
		downitem string
		err      error
		isend    string
	)
	//页面cache控制
	this.Ctx.Output.Header("Cache-Control", "public")

	idstr := this.Ctx.Input.Param(":id")
	id, err = strconv.ParseInt(idstr, 10, 64)

	if err != nil || id <= 0 {
		this.Abort("404")
		return
	}

	//读取数据
	info.Id = id
	err = info.Read()
	if err != nil || info.Status < 0 {
		this.Abort("404")
		return
	}

	//相关影片
	query := rinfo.Query().Filter("mids__icontains", ","+idstr+",")
	query.OrderBy("-Id").One(&rinfo)

	ids := make([]int64, 0)
	midstr := strings.Split(rinfo.Mids, ",")
	for _, s := range midstr {
		i, _ := strconv.ParseInt(s, 10, 64)
		if i > 0 && i != id {
			ids = append(ids, i)
		}
	}
	if len(ids) > 0 {
		q := info.Query().Filter("id__in", ids)
		count, _ := q.Count()
		if count > 0 {
			q.OrderBy("-Id").Limit(10, 0).All(&rmlist, "Id", "Name", "Ename")
		}
	}

	///相关影片
	var item string
	liststring := []string{}
	if len(rmlist) > 0 {
		item = "<p>相关影片：<br />"
		for _, i := range rmlist {
			item = item + fmt.Sprintf("<a href=\"/v/%d/\" target=\"_blank\">%s(%s)</a><br />", i.Id, i.Name, i.Ename)
		}
		item = item + "</p>"
		liststring = append(liststring, item)
	}

	//下载地址json数据
	list := make([]*models.DownAddrInfo, 0)
	down.Query().Filter("mid", id).OrderBy("ep").All(&list)
	count := int64(len(list))
	for i := 1; int64(i) < (count + 1); i++ {
		hdurl := list[i-1].Hdtvurl
		if strings.Contains(hdurl, "mkv") {
			if i < 10 {
				downitem = downitem + fmt.Sprintf("<li id=\"hdtv%d\"><a href=\"%s\" target=\"_blank\">第0%d集.HDTV.1024.中文字幕.mkv</a></li>", i, hdurl, i)
			} else {
				downitem = downitem + fmt.Sprintf("<li id=\"hdtv%d\"><a href=\"%s\" target=\"_blank\">第%d集.HDTV.1024.中文字幕.mkv</a></li>", i, hdurl, i)
			}
		} else {
			if i < 10 {
				downitem = downitem + fmt.Sprintf("<li id=\"hdtv%d\"><a href=\"%s\" target=\"_blank\">第0%d集.HDTV.1024.中文字幕.mp4</a></li>", i, hdurl, i)
			} else {
				downitem = downitem + fmt.Sprintf("<li id=\"hdtv%d\"><a href=\"%s\" target=\"_blank\">第%d集.HDTV.1024.中文字幕.mp4</a></li>", i, hdurl, i)
			}
		}
	}
	if count < (info.Episode + 1) {
		for i := (count + 1); int64(i) < (info.Episode + 1); i++ {
			if i < 10 {
				downitem = downitem + fmt.Sprintf("<li id=\"hdtv%d\">第0%d集.HDTV.1024.中文字幕.mp4</li>", i, i)
			} else {
				downitem = downitem + fmt.Sprintf("<li id=\"hdtv%d\">第%d集.HDTV.1024.中文字幕.mp4</li>", i, i)
			}
		}
	}

	//更新点击
	info.Views++
	info.Update("Views")

	///内页热门列表
	hostlist := info.GetHotList(info.Cid, 16)
	this.Data["hostlist"] = hostlist

	//随机列表
	randlist := info.GetRandList(10)
	this.Data["randlist"] = randlist

	info.Content = strings.Replace(info.Content, "\r\n\r\n", "\r\n", -1)
	info.Content = strings.Replace(info.Content, "\r\n", "<br />", -1)
	//this.Data["downlist"] = downlist
	this.Data["downitem"] = downitem
	this.Data["info"] = info
	if info.Isend == 1 {
		isend = "已完结"
	} else {
		isend = fmt.Sprintf("每周%d更新并播出", info.Updateweek)
	}
	this.Data["isend"] = isend
	this.Data["rmlist"] = strings.Join(liststring, "\n") //相关影片输出
	this.TplName = "_detail.html"
}
