package icons

import "github.com/mattn/go-gtk/gtk"
import "github.com/mattn/go-gtk/gdkpixbuf"

type Icons struct {
	Dir  *gdkpixbuf.Pixbuf
	File *gdkpixbuf.Pixbuf
}

func New () *Icons {
	var icons Icons
	mk := func (stock string) *gdkpixbuf.Pixbuf { return gtk.NewImage().RenderIcon(stock, gtk.ICON_SIZE_SMALL_TOOLBAR, "") }
	icons.Dir  = mk (gtk.STOCK_DIRECTORY)
	icons.File = mk (gtk.STOCK_FILE     )
	return &icons
}

