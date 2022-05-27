package typescript

import (
	"github.com/maasasia/donggu/dictionary"
)

type ArgumentFormatter interface {
	Format(key string, format dictionary.TemplateKeyFormat) string
}
