package mybeego

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
)

type SqlExcutor struct {
	table  string
	wherestr  []string
	orderstr []string
	islimit bool
	offset int
	pagenum int
}

func NewSqlExutor()*SqlExcutor{
	return &SqlExcutor{}
}

func(s *SqlExcutor)Table(str string)*SqlExcutor{
	s.table = str
	return s
}

func(s *SqlExcutor)Where(condition ...string)*SqlExcutor{
	if len(condition)<=0{
		return s
	}
	s.wherestr = append(s.wherestr,condition...)
	return s
}

func(s *SqlExcutor)Order(condition ...string)*SqlExcutor{
	if len(condition)<=0{
		return s
	}
	s.orderstr = append(s.orderstr,condition...)
	return s
}

func(s *SqlExcutor)Limit(page,pagenum int)*SqlExcutor{
	offset :=0
	if page>0{
		offset = (page-1)*pagenum
	}
	s.islimit =true
	s.offset = offset
	s.pagenum = pagenum
	return s
}

func(s *SqlExcutor)Strings()( string, string, error){
	sqlRow :=bytes.NewBufferString(" select * ")
	sqlCount :=bytes.NewBufferString("select count(0) ")
	sql :=bytes.NewBufferString("")
	if len(s.table)<0{
		err := fmt.Errorf("table name is empty")
		return "","",err
	}
	sql.WriteString(fmt.Sprintf(" from %v",s.table))
	if  len(s.wherestr)>0{
		sql.WriteString(" where ")
		for i:=range s.wherestr{
			if i!=0{
				sql.WriteString( " AND ")
			}
			sql.WriteString(s.wherestr[i])
		}
	}
	if len(s.orderstr)>0{
		sql.WriteString("\n order by ")
		sql.WriteString(strings.Join(s.orderstr,","))
	}
	sqlCount.WriteString(sql.String())
	if s.islimit{
		sql.WriteString(fmt.Sprintf("\n limit %v,%v",s.offset,s.pagenum))
	}
	sqlRow.WriteString(sql.String())
	return sqlRow.String(),sqlCount.String(),nil
}

func(s *SqlExcutor)Querys(v interface{})(total int,err error){
	o :=orm.NewOrm()
	var  sqlRow,sqlCount string
	if sr,sc,e :=s.Strings();e!=nil{
		err =e
		return
	}else{
		sqlRow = sr
		sqlCount = sc
	}
	if err=o.Raw(sqlCount).QueryRow(&total);err!=nil{
		return
	}
	if _,err=o.Raw(sqlRow).QueryRows(v);err!=nil{
		return
	}
	return
}
