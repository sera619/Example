package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sera619/GOAssistant/calculator"
	"github.com/sera619/GOAssistant/httpserver"
)

// init http Client for random facts
var client *http.Client
var uselessAPI = "https://uselessfacts.jsph.pl/random.json?language=en"

//Resource contructor and generation
type randomFact struct {
	Text string `json:"text"`
}
type Resource interface {
	Name() string
	Content() []byte
}
type StaticResource struct {
	StaticName    string
	StaticContent []byte
}

func (r *StaticResource) Name() string {
	return r.StaticName
}

func (r *StaticResource) Content() []byte {
	return r.StaticContent
}

// Mainloop
func main() {
	TheApp()
}

//Api Request for a useless fact
func getRandomFact() (randomFact, error) {
	var fact randomFact
	resp, err := client.Get(uselessAPI)
	if err != nil {
		return randomFact{}, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&fact)
	if err != nil {
		return randomFact{}, err
	}
	return fact, nil
}

// init a static Resource for images etc
func NewStaticResource(name string, content []byte) *StaticResource {
	return &StaticResource{
		StaticName:    name,
		StaticContent: content,
	}
}

// exmaple for maths
func DummDummRechner() {
	fmt.Println("\n\nDumm Dumm Rechner in GO\n______________________________")
	fmt.Println("\nAddition:\n23 + 43 =", calculator.Add(23, 43))
	fmt.Println("\nMultiply:\n44 * 32 =", calculator.Multiply(44, 32))
	fmt.Println("\nSubstract:\n30 - 12 =", calculator.Substract(30, 12))
	fmt.Println("\nDevide:\n123242 / 123 =", calculator.Devide(123242, 123))
	fmt.Println("______________________________")
}

// Graphical Userinterface implementation
func TheApp() {
	client = &http.Client{Timeout: 10 * time.Second}
	r, _ := LoadResourceFromPath("./public/img/icon.png")
	myApp := app.New()
	// cyan font color for maintext
	mainTextColor := color.RGBA{
		R: 19,
		G: 183,
		B: 182,
		A: 255,
	}
	myWindow := myApp.NewWindow("GO Assistant © S3R43o3")
	//Sizing Window and adding icons
	//myWindow.Resize(fyne.NewSize(600, 300))
	myWindow.SetIcon(r)

	myWindow.CenterOnScreen()

	// create widgets for Text and Button
	httpServerLabel := canvas.NewText("Click the Start-Button to start the Webserver.", mainTextColor)
	httpServerLabel.Alignment = fyne.TextAlignCenter
	httpServerLabel.TextSize = 10
	httpServerLabel2 := canvas.NewText("Click the GetFact-Button to get a random useless fact.", mainTextColor)
	httpServerLabel2.TextSize = 10
	httpServerLabel2.Alignment = fyne.TextAlignCenter
	httpServerLabel3 := canvas.NewText("Click the 'Dumm Rechner' to get an example caluclation.", mainTextColor)

	httpServerLabel3.Alignment = fyne.TextAlignCenter
	httpServerLabel3.TextSize = 10

	myButton2 := widget.NewButton("Exit", func() {
		myApp.Quit()
	})

	factLabel := widget.NewLabel("")
	factLabel.Wrapping = fyne.TextWrapWord
	factButton := widget.NewButton("Get Fact", func() {
		fact, err := getRandomFact()
		if err != nil {
			dialog.ShowError(err, myWindow)
		} else {
			factLabel.SetText(fact.Text)

		}
	})
	calcButton := widget.NewButton("Dumm Rechner", func() {
		DummDummRechner()
	})
	// generate a Canvas and add text to it
	title := canvas.NewText("GO Assistant", mainTextColor)
	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 24

	startHTTPbutton := widget.NewButton("Start", func() {
		factLabel.SetText("HTTP Server is Running on\nhttp://localhost:8080/home\nClose the Window to exit.")
		defer httpserver.RunServer()
	})
	brandLabel := canvas.NewText("Development by S3R43o3", color.White)
	brandLabel.TextSize = 8
	// add buttons to Hbox Container
	buttonHBox := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), startHTTPbutton, factButton, calcButton, myButton2, layout.NewSpacer())
	// add widgets to V-BoxLayout Container
	vBox := container.New(layout.NewVBoxLayout(), title, layout.NewSpacer(), factLabel, httpServerLabel, httpServerLabel2, httpServerLabel3, layout.NewSpacer(), buttonHBox, layout.NewSpacer(), brandLabel, layout.NewSpacer())

	myWindow.SetContent(vBox)
	myWindow.ShowAndRun()

}

// helper funktion load resource from path and init as Staticresource
func LoadResourceFromPath(path string) (Resource, error) {
	bytes, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	name := filepath.Base(path)
	return NewStaticResource(name, bytes), nil
}
