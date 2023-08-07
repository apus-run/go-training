package model

type Model[E any] interface {
	TableName() string
	ToEntity() E
	FromEntity(E) any
}
