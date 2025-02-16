package models

import (
	"github.com/iancoleman/strcase"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"
)

var g_ *gen.Generator

func g() *gen.Generator {
	if g_ == nil {
		cfg := gen.Config{
			OutPath:       "services/repo",
			ModelPkgPath:  "./models",
			FieldNullable: true,
		}
		g_ = gen.NewGenerator(cfg)
	}
	return g_
}

// This method is called by the dbgen tool and it will generate the models files
func Generate() {
	g().Execute()
}

// call this method to add a modelSpec to generation
func generateModel(modelSpec helper.Object) {
	g := g()
	// g.ApplyBasic(g.GenerateModelFrom(modelSpec))
	g.ApplyInterface(func(Querier) {}, g.GenerateModelFrom(modelSpec))
}

type Querier interface {
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id int) (*gen.T, error) // returns data as pointer and error
}

var idFieldSpec *FieldSpec = &FieldSpec{
	name:    "ID",
	typ:     "int64",
	gormTag: "primarykey",
	comment: "",
}

var createdAtFieldSpec *FieldSpec = &FieldSpec{
	name:    "CreatedAt",
	typ:     "time.Time",
	gormTag: `autoCreateTime`,
	comment: "",
}

var deletedAtFieldSpec *FieldSpec = &FieldSpec{
	name:    "DeletedAt",
	typ:     "gorm.DeletedAt",
	gormTag: "index",
	comment: "",
}
var updatedAtFieldSpec *FieldSpec = &FieldSpec{
	name:    "UpdatedAt",
	typ:     "time.Time",
	gormTag: "autoUpdateTime",
	comment: "",
}

// ModelSpec defines a model
type ModelSpec struct {
	structName string
	tableName  string
	fields     []helper.Field
}

// TableName return table name
func (d *ModelSpec) TableName() string { return d.tableName }

// StructName return struct name
func (d *ModelSpec) StructName() string { return d.structName }

// FileName return file name
func (d *ModelSpec) FileName() string { return d.tableName }

// ImportPkgPaths return import package paths
func (d *ModelSpec) ImportPkgPaths() []string { return nil }

// Fields return fields
func (d *ModelSpec) Fields() []helper.Field { return d.fields }

// FieldSpec demo field
type FieldSpec struct {
	name    string
	typ     string
	gormTag string
	comment string
}

// Name return name
func (f *FieldSpec) Name() string { return f.name }

// Type return field type
func (f *FieldSpec) Type() string { return f.typ }

// ColumnName return column name
func (f *FieldSpec) ColumnName() string { return strcase.ToSnake(f.name) }

// GORMTag return gorm tag
func (f *FieldSpec) GORMTag() string { return f.gormTag }

// JSONTag return json tag
func (f *FieldSpec) JSONTag() string { return "-" }

// Tag return new tag
func (f *FieldSpec) Tag() field.Tag { return field.Tag{} }

// Comment return comment
func (f *FieldSpec) Comment() string { return f.comment }
