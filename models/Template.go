package models

import (
	"time"
	"github.com/lifei6671/mindoc/conf"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"errors"
)

type Template struct {
	TemplateId	int 		`orm:"column(template_id);pk;auto;unique;" json:"template_id"`
	TemplateName string 	`orm:"column(template_name);size(500);" json:"template_name"`
	MemberId int 			`orm:"column(member_id);index" json:"member_id"`
	BookId int				`orm:"column(book_id);index" json:"book_id"`
	//是否是全局模板：0 否/1 是; 全局模板在所有项目中都可以使用；否则只能在创建模板的项目中使用
	IsGlobal int			`orm:"column(is_global);default(0)" json:"is_global"`
	TemplateContent string 	`orm:"column(template_content);type(text);null" json:"template_content"`
	CreateTime time.Time     `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	ModifyTime time.Time     `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt   int           `orm:"column(modify_at);type(int)" json:"-"`
	Version    int64         `orm:"type(bigint);column(version)" json:"version"`
}

// TableName 获取对应数据库表名.
func (m *Template) TableName() string {
	return "templates"
}

// TableEngine 获取数据使用的引擎.
func (m *Template) TableEngine() string {
	return "INNODB"
}

func (m *Template) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewTemplate() *Template  {
	return &Template{}
}

//查询指定ID的模板
func (t *Template) Find(templateId int) (*Template,error) {
	if templateId <= 0 {
		return t, ErrInvalidParameter
	}

	o := orm.NewOrm()

	err := o.QueryTable(t.TableNameWithPrefix()).Filter("template_id",templateId).One(t)

	if err != nil {
		logs.Error("查询模板时失败 ->%s",err)
	}
	return t,err
}

//查询属于指定项目的模板.
func (t *Template) FindByBookId(bookId int) ([]*Template,error) {
	if bookId <= 0 {
		return nil,ErrInvalidParameter
	}
	o := orm.NewOrm()

	var templateList []*Template

	_,err := o.QueryTable(t.TableNameWithPrefix()).Filter("book_id",bookId).OrderBy("-template_id").All(&templateList)

	if err != nil {
		beego.Error("查询模板列表失败 ->",err)
	}
	return templateList,err
}

//查询指定项目所有可用模板列表.
func (t *Template) FindAllByBookId(bookId int) ([]*Template,error) {
	if bookId <= 0 {
		return nil,ErrInvalidParameter
	}
	o := orm.NewOrm()

	cond := orm.NewCondition()

	cond1 := cond.And("book_id",bookId).Or("is_global",1)

	qs := o.QueryTable(t.TableNameWithPrefix())

	var templateList []*Template

	_,err := qs.SetCond(cond1).OrderBy("-template_id").All(&templateList)

	if err != nil {
		beego.Error("查询模板列表失败 ->",err)
	}
	return templateList,err
}
//删除一个模板
func (t *Template) Delete(templateId int,memberId int) error {
	if templateId <= 0 {
		return ErrInvalidParameter
	}

	o := orm.NewOrm()

	qs := o.QueryTable(t.TableNameWithPrefix()).Filter("template_id",templateId)

	if memberId > 0 {
		qs = qs.Filter("member_id",memberId)
	}
	_,err := qs.Delete()

	if err != nil {
		beego.Error("删除模板失败 ->",err)
	}
	return err
}

//添加或更新模板
func (t *Template) Save(cols ...string) (err error) {

	if t.BookId <= 0 {
		return ErrInvalidParameter
	}
	o := orm.NewOrm()

	if !o.QueryTable(NewBook()).Filter("book_id",t.BookId).Exist() {
		return errors.New("项目不存在")
	}
	if !o.QueryTable(NewMember()).Filter("member_id",t.MemberId).Filter("status",0).Exist() {
		return errors.New("用户已被禁用")
	}
	if t.TemplateId > 0 {
		t.Version = time.Now().Unix()
		t.ModifyTime = time.Now()
		_,err = o.Update(t,cols...)
	}else{
		t.CreateTime = time.Now()
		_,err = o.Insert(t)
	}

	return
}








