package types

import (
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// String implements fmt.Stringer interface
func (cpp CurrencyPairProviders) String() string {
	out, _ := yaml.Marshal(cpp)
	return string(out)
}

func (cpp CurrencyPairProviders) Equal(cpp2 *CurrencyPairProviders) bool {
	if !strings.EqualFold(cpp.BaseDenom, cpp2.BaseDenom) || !strings.EqualFold(cpp.QuoteDenom, cpp2.QuoteDenom){
		return false
	}

	return reflect.DeepEqual(cpp.Providers, cpp2.Providers)
}

// CurrencyPairProvidersList is array of CurrencyPairProviders
type CurrencyPairProvidersList []CurrencyPairProviders

func (cppl CurrencyPairProvidersList) String() (out string) {
	for _, v := range cppl {
		out += v.String() + "\n"
	}

	return strings.TrimSpace(out)
}
