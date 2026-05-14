package clipboard

import atotto "github.com/atotto/clipboard"

// WriteText 写入系统剪贴板。
func WriteText(text string) error {
	return atotto.WriteAll(text)
}
