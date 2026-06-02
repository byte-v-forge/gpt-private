package plugins

import (
	"github.com/byte-v-forge/gpt-private/plugins/gopay"
	"github.com/byte-v-forge/gpt-private/plugins/privateflows"
	"github.com/byte-v-forge/gpt/pkg/gptplugin"
)

func Plugins() []gptplugin.Plugin {
	return []gptplugin.Plugin{
		gopay.Plugin(),
		privateflows.Plugin(),
	}
}
