package donggu

import "github.com/ghost/donggu/generated"

func NewDonggu(resolver generated.ResolverFunc) *generated.Donggu {
	return &generated.Donggu{Resolver: resolver}
}
