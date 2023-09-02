package repository

type ITemplateRepository interface{}

type templateRepository struct {}

func NewTemplate(
) ITemplateRepository {
	return templateRepository{}
}
