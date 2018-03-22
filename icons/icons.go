package icons

import "github.com/mattn/go-gtk/gtk"
import "github.com/mattn/go-gtk/gdkpixbuf"

type Icon = *gdkpixbuf.Pixbuf

type Icons struct {
	Dir  Icon
	File Icon
}

func New () *Icons {
	var icons Icons
	mk := func (stock string) Icon { return gtk.NewImage().RenderIcon(stock, gtk.ICON_SIZE_SMALL_TOOLBAR, "") }
	icons.Dir  = mk (gtk.STOCK_DIRECTORY)
	icons.File = mk (gtk.STOCK_FILE     )
	return &icons
}

