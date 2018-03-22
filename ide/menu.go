package ide

import "log"
import "github.com/mattn/go-gtk/gdk"
import "github.com/mattn/go-gtk/gtk"
import "fmt"

var _ = log.Printf

func (ide *IDE) MakeMenu () {

	ide.Menubar = gtk.NewMenuBar ()
	addCascade := func (label string, fill func (*gtk.Menu)) {
		cascademenu := gtk.NewMenuItemWithMnemonic(label)
		ide.Menubar.Append(cascademenu)
		submenu := gtk.NewMenu()
		cascademenu.SetSubmenu(submenu)
		fill (submenu)
	}
	makeItem := func (label string, action func ()) *gtk.MenuItem {
		menuitem := gtk.NewMenuItemWithMnemonic(label)
		menuitem.Connect("activate", action)
		return menuitem
	}
	
	addCascade ("_File", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("E_xit", gtk.MainQuit))
		//submenu.Append(makeItem ("_Open", ide.Editor.LoadFileFromDialog))
	})

	addCascade ("_Edit", func (submenu *gtk.Menu) {
		{
			item := makeItem ("_Complete", ide.Complete)
			item.AddAccelerator ("activate", ide.Accel, gdk.KEY_space, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
			submenu.Append(item)
		}
		submenu.Append(makeItem ("_Fonts", func () {
			fsd := gtk.NewFontSelectionDialog("Font")
			fsd.Response(func() {
				fmt.Println(fsd.GetFontName())
				fsd.Destroy()
			})
			fsd.SetTransientFor(ide.Window)
			fsd.Run()
		}))
		//submenu.Append(makeItem ("_Open", ide.Editor.LoadFileFromDialog))
	})

	addCascade ("_View", func (submenu *gtk.Menu) {
		{
			item := makeItem ("_Files", ide.SelectFiles)
			item.AddAccelerator ("activate", ide.Accel, gdk.KEY_F, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
			submenu.Append(item)
		}
		{
			item := makeItem ("_Names", ide.SelectNames)
			item.AddAccelerator ("activate", ide.Accel, gdk.KEY_N, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
			submenu.Append(item)
		}
	})

	addCascade ("_Build", func (submenu *gtk.Menu) {
		item := makeItem ("_Make", ide.Make)
		item.AddAccelerator ("activate", ide.Accel, gdk.KEY_F9, 0, gtk.ACCEL_VISIBLE)
		submenu.Append(item)
		//submenu.Append(makeItem ("_Open", ide.Editor.LoadFileFromDialog))
	})

	addCascade ("_Tech", func (submenu *gtk.Menu) {
		item := makeItem ("_Show Selection", ide.TechShowSelection)
		submenu.Append(item)
	})

	addCascade ("_Help", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("_About", func () {
			dialog := gtk.NewAboutDialog()
			dialog.SetName("Go-Gtk Demo!")
			dialog.SetProgramName("demo")
			dialog.SetAuthors([]string{"Kouzdra"})
			//dir, _ := filepath.Split(os.Args[0])
			//imagefile := filepath.Join(dir, "../../data/mattn-logo.png")
			//pixbuf, _ := gdkpixbuf.NewPixbufFromFile(imagefile)
			//dialog.SetLogo(pixbuf)
			dialog.SetLicense("The library is available under the same terms and conditions as the Go, the BSD style license, and the LGPL (Lesser GNU Public License). The idea is that if you can use Go (and Gtk) in a project, you should also be able to use go-gtk.")
			dialog.SetWrapLicense(true)
			dialog.Run()
			dialog.Destroy()
		}))
	})
}

