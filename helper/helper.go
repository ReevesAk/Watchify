package helper

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	bolt "go.etcd.io/bbolt"
)

type dbParams struct {
	MonthOfWork     string
	Timer           string
	Title           string
	WorkDescription string
}

type Objects struct {
	button1    *widget.Label
	label      *widget.Label
	but2, but3 fyne.CanvasObject
}

func openDatabase() *bolt.DB {

	database, err := bolt.Open("fyne.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return database
}

func launchApp() fyne.App {
	myApp := app.New()
	return myApp
}

var (
	counter        = 0
	incrementTimer *time.Ticker
)

func convertTicker(timeInSeconds int) (str string) {
	minutes := timeInSeconds / 60
	seconds := timeInSeconds % 60
	hours := timeInSeconds / 60 / 60
	str = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	return
}

func saveLog(Bucket, key string, values []byte) {
	fyneDb := openDatabase()

	if err := fyneDb.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(Bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		if err := bucket.Put([]byte(key), []byte(values)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func retrieveLog(Bucket, key string) (str []byte) {
	fyneDb := openDatabase()

	if err := fyneDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))
		str = bucket.Get([]byte(key))
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return str
}

func startTimer(text *widget.Label) {
	interval := float64(1000)
	incrementTimer = time.NewTicker(time.Duration(interval) * time.Millisecond)

	for _ = range incrementTimer.C {
		counter++
		text.SetText(fmt.Sprint(convertTicker(counter)))
	}
}

func Run() {
	timerApp := launchApp()

	obj := Objects{}

	theme := theme.LightTheme()
	theme.PrimaryColor()
	timerApp.Settings().SetTheme(theme)
	window := timerApp.NewWindow("Watchify")

	window.Resize(fyne.NewSize(600, 200))

	obj.button1 = widget.NewLabelWithStyle("00:00:00",
		fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})

	month := widget.NewEntry()
	month.Text = ""

	projectTitle := widget.NewEntry()
	projectTitle.Text = ""

	workDescription := widget.NewMultiLineEntry()
	workDescription.Text = ""

	window.SetContent(widget.NewVBox(
		obj.button1,
		widget.NewButton("Start/Resume", func() {
			go startTimer(obj.button1)
			obj.button1.SetText(fmt.Sprint(counter))
		}),

		widget.NewButton("Pause/Stop", func() {
			incrementTimer.Stop()
		}),

		fyne.NewContainerWithLayout(layout.NewFormLayout(),
			widget.NewLabel("Enter Month"),
			month),

		fyne.NewContainerWithLayout(layout.NewFormLayout(),
			widget.NewLabel("Project Title:"),
			projectTitle),

		fyne.NewContainerWithLayout(layout.NewFormLayout(),
			widget.NewLabel("Description:"),
			widget.NewHBox(
				workDescription,
			),
		),

		widget.NewButton("...", func() {
			widget.ShowPopUp(
				widget.NewVBox(
					widget.NewButton("Save", func() {
						save := dbParams{
							MonthOfWork:     month.Text,
							Title:           projectTitle.Text,
							Timer:           obj.button1.Text,
							WorkDescription: workDescription.Text,
						}

						b := fmt.Sprintf("Month: %s\nTitle: %s\nTimer: %s\n"+
							"Description: %s\n",
							save.MonthOfWork, save.Title,
							save.Timer, save.WorkDescription)

						saveLog(month.Text, month.Text, []byte(b))
						fmt.Println(b)
					}),
					widget.NewButton("WorkLog", func() {
						window2 := timerApp.NewWindow("Log")
						window2.Show()
						window2.Resize(fyne.NewSize(400, 300))

						c := retrieveLog(month.Text, month.Text)

						fmt.Println(c)
						obj.label = widget.NewLabel(string(c))

						window2.SetContent(
							obj.label,
						)
					}),
				),
				window.Canvas(),
			)
		}),

		widget.NewButton("Quit App", func() {
			timerApp.Quit()
		}),
	))

	window.ShowAndRun()
}
