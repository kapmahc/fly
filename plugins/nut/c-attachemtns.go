package nut

import (
	"encoding/base64"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// IndexAttachments list
// @router /attachments [get]
func (p *Plugin) IndexAttachments() {
	p.LayoutDashboard()
	var items []Attachment
	if _, err := orm.NewOrm().QueryTable(new(Attachment)).
		OrderBy("-updated_at").
		Filter("user_id", p.CurrentUser().ID).
		All(&items); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[TITLE] = Tr(p.Locale(), "nut.attachments.index.title")

	p.TplName = "nut/attachments/index.html"
}

// DestroyAttachment remove
// @router /attachments/:id [delete]
func (p *Plugin) DestroyAttachment() {
	p.MustSignIn()
	var item Attachment
	o := orm.NewOrm()
	err := o.QueryTable(&item).
		Filter("id", p.Ctx.Input.Param(":id")).
		One(&item)
	if err == nil {
		if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
			err = Te(p.Locale(), "errors.not-allow")
		}
	}
	if err == nil {
		_, err = o.Delete(&item)
	}
	if err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	p.Data["json"] = H{"ok": true}
	p.ServeJSON()
}

// PostAttachmentsUeditor ueditor
// @router /attachments/ueditor [get,post]
func (p *Plugin) PostAttachmentsUeditor() {
	p.MustSignIn()
	switch p.GetString("action") {
	case "config":
		p.ueditorConfig()
	case "uploadimage":
		p.ueditorUploadFile()
	case "uploadscrawl":
		p.ueditorUploadScrawl()
	case "uploadvideo":
		p.ueditorUploadFile()
	case "uploadfile":
		p.ueditorUploadFile()
	case "listimage":
		p.ueditorImagesManager()
	case "listfile":
		p.ueditorFilesManager()
	case "catchimage":
		p.ueditorCatchImage()
	default:
		beego.Warn("TODO: ueditor action")
	}
}

const (
	stateSuccess = "SUCCESS"
	stateFailed  = "FAILED"
)

func (p *Plugin) ueditorUploadScrawl() {
	buf, err := base64.StdEncoding.DecodeString(p.GetString("upfile"))
	var att *Attachment
	if err == nil {
		att, err = p.writeToS3("scrawl.png", buf, int64(len(buf)))
	}
	p.ueditorWrite(att, err)
}

func (p *Plugin) ueditorWrite(att *Attachment, err error) {
	if err == nil {
		p.Data["json"] = H{
			"state":    stateSuccess,
			"url":      att.URL,
			"title":    "",
			"original": att.Title,
		}
	} else {
		beego.Error(err)
		p.Data["json"] = H{
			"state": stateFailed,
		}
	}
	p.ServeJSON()
}

func (p *Plugin) ueditorUploadFile() {
	att, err := p.UploadFile("upfile")
	p.ueditorWrite(att, err)
}

func (p *Plugin) ueditorManager(f func(*Attachment) bool) {
	o := orm.NewOrm()
	var items []Attachment
	_, err := o.QueryTable(new(Attachment)).
		Filter("user_id", p.CurrentUser().ID).
		All(&items, "media_type", "url")
	var list []H
	if err == nil {
		for _, it := range items {
			if f(&it) {
				list = append(list, H{"url": it.URL})
			}
		}
	}
	if err == nil {
		p.Data["json"] = H{
			"state": stateSuccess,
			"list":  list,
			"start": 0,
			"total": len(list),
		}
	} else {
		p.Data["json"] = H{"state": stateFailed}
	}
	p.ServeJSON()
}
func (p *Plugin) ueditorImagesManager() {
	p.ueditorManager(func(a *Attachment) bool {
		return a.IsPicture()
	})
}
func (p *Plugin) ueditorFilesManager() {
	p.ueditorManager(func(a *Attachment) bool {
		return !a.IsPicture()
	})
}
func (p *Plugin) ueditorCatchImage() {
	// TODO
}

func (p *Plugin) ueditorConfig() {
	/* 前后端通信相关的配置,注释只允许使用多行方式 */
	p.Data["json"] = H{
		/* 上传图片配置项 */
		"imageActionName":     "uploadimage",                                     /* 执行上传图片的action名称 */
		"imageFieldName":      "upfile",                                          /* 提交的图片表单名称 */
		"imageMaxSize":        2048000,                                           /* 上传大小限制，单位B */
		"imageAllowFiles":     []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 上传图片格式显示 */
		"imageCompressEnable": true,                                              /* 是否压缩图片,默认是true */
		"imageCompressBorder": 1600,                                              /* 图片压缩最长边限制 */
		"imageInsertAlign":    "none",                                            /* 插入的图片浮动方式 */
		"imageUrlPrefix":      "",                                                /* 图片访问路径前缀 */
		"imagePathFormat":     "/images/{yyyy}{mm}{dd}/{time}{rand:6}",           /* 上传保存路径,可以自定义保存路径和文件名格式 */
		/* {filename} 会替换成原文件名,配置这项需要注意中文乱码问题 */
		/* {rand:6} 会替换成随机数,后面的数字是随机数的位数 */
		/* {time} 会替换成时间戳 */
		/* {yyyy} 会替换成四位年份 */
		/* {yy} 会替换成两位年份 */
		/* {mm} 会替换成两位月份 */
		/* {dd} 会替换成两位日期 */
		/* {hh} 会替换成两位小时 */
		/* {ii} 会替换成两位分钟 */
		/* {ss} 会替换成两位秒 */
		/* 非法字符 \ : * ? " < > | */
		/* 具请体看线上文档: fex.baidu.com/ueditor/#use-format_upload_filename */
		/* 涂鸦图片上传配置项 */
		"scrawlActionName":  "uploadscrawl",                          /* 执行上传涂鸦的action名称 */
		"scrawlFieldName":   "upfile",                                /* 提交的图片表单名称 */
		"scrawlPathFormat":  "/images/{yyyy}{mm}{dd}/{time}{rand:6}", /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"scrawlMaxSize":     2048000,                                 /* 上传大小限制，单位B */
		"scrawlUrlPrefix":   "",                                      /* 图片访问路径前缀 */
		"scrawlInsertAlign": "none",
		/* 截图工具上传 */
		"snapscreenActionName":  "uploadimage",                           /* 执行上传截图的action名称 */
		"snapscreenPathFormat":  "/images/{yyyy}{mm}{dd}/{time}{rand:6}", /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"snapscreenUrlPrefix":   "",                                      /* 图片访问路径前缀 */
		"snapscreenInsertAlign": "none",                                  /* 插入的图片浮动方式 */
		/* 抓取远程图片配置 */
		"catcherLocalDomain": []string{"127.0.0.1", "localhost", "image.baidu.com"},
		"catcherActionName":  "catchimage",                                      /* 执行抓取远程图片的action名称 */
		"catcherFieldName":   "source",                                          /* 提交的图片列表表单名称 */
		"catcherPathFormat":  "/images/{yyyy}{mm}{dd}/{time}{rand:6}",           /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"catcherUrlPrefix":   "",                                                /* 图片访问路径前缀 */
		"catcherMaxSize":     2048000,                                           /* 上传大小限制，单位B */
		"catcherAllowFiles":  []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 抓取图片格式显示 */
		/* 上传视频配置 */
		"videoActionName": "uploadvideo",                           /* 执行上传视频的action名称 */
		"videoFieldName":  "upfile",                                /* 提交的视频表单名称 */
		"videoPathFormat": "/videos/{yyyy}{mm}{dd}/{time}{rand:6}", /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"videoUrlPrefix":  "",                                      /* 视频访问路径前缀 */
		"videoMaxSize":    102400000,                               /* 上传大小限制，单位B，默认100MB */
		"videoAllowFiles": []string{
			".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
			".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
		}, /* 上传视频格式显示 */
		/* 上传文件配置 */
		"fileActionName": "uploadfile",                           /* controller里,执行上传视频的action名称 */
		"fileFieldName":  "upfile",                               /* 提交的文件表单名称 */
		"filePathFormat": "/files/{yyyy}{mm}{dd}/{time}{rand:6}", /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"fileUrlPrefix":  "",                                     /* 文件访问路径前缀 */
		"fileMaxSize":    51200000,                               /* 上传大小限制，单位B，默认50MB */
		"fileAllowFiles": []string{
			".png", ".jpg", ".jpeg", ".gif", ".bmp",
			".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
			".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
			".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
			".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml",
		}, /* 上传文件格式显示 */
		/* 列出指定目录下的图片 */
		"imageManagerActionName":  "listimage",                                       /* 执行图片管理的action名称 */
		"imageManagerListPath":    "/images/",                                        /* 指定要列出图片的目录 */
		"imageManagerListSize":    20,                                                /* 每次列出文件数量 */
		"imageManagerUrlPrefix":   "",                                                /* 图片访问路径前缀 */
		"imageManagerInsertAlign": "none",                                            /* 插入的图片浮动方式 */
		"imageManagerAllowFiles":  []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 列出的文件类型 */
		/* 列出指定目录下的文件 */
		"fileManagerActionName": "listfile", /* 执行文件管理的action名称 */
		"fileManagerListPath":   "/files/",  /* 指定要列出文件的目录 */
		"fileManagerUrlPrefix":  "",         /* 文件访问路径前缀 */
		"fileManagerListSize":   20,         /* 每次列出文件数量 */
		"fileManagerAllowFiles": []string{
			".png", ".jpg", ".jpeg", ".gif", ".bmp",
			".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
			".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
			".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
			".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml",
		}, /* 列出的文件类型 */
	}

	p.ServeJSON()
}
