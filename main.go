package main

import (
	"Personal-web/connection"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type Projects struct {
	Id          int
	Title       string
	Sdate       time.Time
	Edate       time.Time
	Duration    string
	Descript    string
	Technologys []string
	Tech1       bool
	Tech2       bool
	Tech3       bool
	Tech4       bool
	Image       string
}

var dataProject = []Projects{
	{
		Title: "Mobile App 2019",
		// Sdate:    "26 January 2019",
		// Edate:    "05 March 2019",
		Duration: "3 Month",
		Descript: "App that used for dumbways student, it was deployed and can downloaded on playstore.<br />Happy download",
		Tech1:    true,
		Tech2:    false,
		Tech3:    true,
		Tech4:    true,
	},
	{
		Title: "Web App 2020",
		// Sdate:    "26 August 2020",
		// Edate:    "05 December 2020",
		Duration: "2 Month",
		Descript: "App that used for dumbways student, it was deployed and can downloaded on playstore.<br />Happy download",
		Tech1:    false,
		Tech2:    true,
		Tech3:    true,
		Tech4:    true,
	},
	{
		Title: "Web App 2023",
		// Sdate:    "26 March 2023",
		// Edate:    "05 September 2023",
		Duration: "3 Month",
		Descript: "App that used for dumbways student, it was deployed and can downloaded on playstore.<br />Happy download",
		Tech1:    true,
		Tech2:    true,
		Tech3:    false,
		Tech4:    true,
	},
	{
		Title: "Mobile App 2024",
		// Sdate:    "26 Jully 2024",
		// Edate:    "05 November 2024",
		Duration: "5 Month",
		Descript: "App that used for dumbways student, it was deployed and can downloaded on playstore.<br />Happy download",
		Tech1:    true,
		Tech2:    true,
		Tech3:    true,
		Tech4:    true,
	},
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	connection.ConnectDB()
	e := echo.New()

	// root statis untuk mengakses folder public
	e.Static("/public", "public") //public

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	// renderer
	e.Renderer = t

	// routing
	e.GET("/", home)
	e.GET("/contact", contactMe)
	e.GET("/project", myProject)
	e.GET("/project-detail/:id", projectDetail)
	e.POST("/add-project", addProject)
	e.GET("/delete/:id", deleteProject)
	e.GET("/edit-project/:id", editProject)
	e.POST("/edit/:id", edit)

	fmt.Println("localhost: 5000 sucssesfully")
	e.Logger.Fatal(e.Start("localhost: 5000"))
}

func home(c echo.Context) error {

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, technologys, description, image FROM public.tb_project;")

	var result []Projects
	for data.Next() {
		var each = Projects{}

		err := data.Scan(&each.Id, &each.Title, &each.Sdate, &each.Edate, &each.Technologys, &each.Descript, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}
		// Duration
		// formatDate := "2006/01/02"

		durasi := each.Edate.Sub(each.Sdate)
		var Durations string

		if durasi.Hours()/24 < 7 {
			Durations = strconv.FormatFloat(durasi.Hours()/24, 'f', 0, 64) + " Days"
		} else if durasi.Hours()/24/7 < 4 {
			Durations = strconv.FormatFloat(durasi.Hours()/24/7, 'f', 0, 64) + " Weeks"
		} else if durasi.Hours()/24/30 < 12 {
			Durations = strconv.FormatFloat(durasi.Hours()/24/30, 'f', 0, 64) + " Months"
		} else {
			Durations = strconv.FormatFloat(durasi.Hours()/24/30/12, 'f', 0, 64) + " Years"
		}

		each.Duration = Durations

		result = append(result, each)
	}
	project := map[string]interface{}{
		"Projects": result,
	}
	return c.Render(http.StatusOK, "index.html", project)
}

func contactMe(c echo.Context) error {
	return c.Render(http.StatusOK, "contact-me.html", nil)
}

func myProject(c echo.Context) error {
	return c.Render(http.StatusOK, "myProject.html", nil)
}

func addProject(c echo.Context) error {
	title := c.FormValue("project-name")
	sDate := c.FormValue("start-date")
	eDate := c.FormValue("end-date")
	desc := c.FormValue("description")
	tech1 := false
	tech2 := false
	tech3 := false
	tech4 := false

	formatDate := "2006-01-02"
	SD, _ := time.Parse(formatDate, sDate)
	// sDateFormat := SD.Format("02 January 2006")
	ED, _ := time.Parse(formatDate, eDate)
	// eDateFormat := ED.Format("02 January 2006")

	durasi := ED.Sub(SD)

	var Durations string
	if durasi.Hours()/24 < 7 {
		Durations = strconv.FormatFloat(durasi.Hours()/24, 'f', 0, 64) + " Days"
	} else if durasi.Hours()/24/7 < 4 {
		Durations = strconv.FormatFloat(durasi.Hours()/24/7, 'f', 0, 64) + " Weeks"
	} else if durasi.Hours()/24/30 < 12 {
		Durations = strconv.FormatFloat(durasi.Hours()/24/30, 'f', 0, 64) + " Months"
	} else {
		Durations = strconv.FormatFloat(durasi.Hours()/24/30/12, 'f', 0, 64) + " Years"
	}

	// if checked
	if c.FormValue("Python") != "" {
		tech1 = true
	}
	if c.FormValue("reactJS") != "" {
		tech2 = true
	}
	if c.FormValue("Javascript") != "" {
		tech3 = true
	}
	if c.FormValue("nodeJS") != "" {
		tech4 = true
	}

	var newAdd = Projects{
		Title: title,
		// Sdate:    sDateFormat,
		// Edate:    eDateFormat,
		Duration: Durations,
		Descript: desc,
		Tech1:    tech1,
		Tech2:    tech2,
		Tech3:    tech3,
		Tech4:    tech4,
	}

	dataProject = append(dataProject, newAdd)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Projects{}

	for i, data := range dataProject {
		if id == i {
			ProjectDetail = Projects{
				Title:    data.Title,
				Sdate:    data.Sdate,
				Edate:    data.Edate,
				Duration: data.Duration,
				Descript: data.Descript,
				Tech1:    data.Tech1,
				Tech2:    data.Tech2,
				Tech3:    data.Tech3,
				Tech4:    data.Tech4,
			}
		}
	}

	detailProject := map[string]interface{}{
		"Projects": ProjectDetail,
	}
	return c.Render(http.StatusOK, "projectDetail.html", detailProject)
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	dataProject = append(dataProject[:id], dataProject[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	edit := Projects{}
	for i, data := range dataProject {
		// convertSD, _ := time.Parse("02 January 2006", data.Sdate)
		// convertED, _ := time.Parse("02 January 2006", data.Edate)

		if id == i {
			edit = Projects{
				Title: data.Title,
				// Sdate:    convertSD.Format("2006-01-02"),
				// Edate:    convertED.Format("2006-01-02"),
				Duration: data.Duration,
				Descript: data.Descript,
				Tech1:    data.Tech1,
				Tech2:    data.Tech2,
				Tech3:    data.Tech3,
				Tech4:    data.Tech4,
			}
		}
	}

	editResult := map[string]interface{}{
		"Edit": edit,
		"Id":   id,
	}

	return c.Render(http.StatusOK, "updateProject.html", editResult)

}

func edit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	title := c.FormValue("project-name")
	sDate := c.FormValue("start-date")
	eDate := c.FormValue("end-date")
	desc := c.FormValue("description")
	tech1 := false
	tech2 := false
	tech3 := false
	tech4 := false

	formatDate := "2006-01-02"
	SD, _ := time.Parse(formatDate, sDate)
	// sDateFormat := SD.Format("02 January 2006")
	ED, _ := time.Parse(formatDate, eDate)
	// eDateFormat := ED.Format("02 January 2006")

	durasi := ED.Sub(SD)

	var Durations string
	if durasi.Hours()/24 < 7 {
		Durations = strconv.FormatFloat(durasi.Hours()/24, 'f', 0, 64) + " Days"
	} else if durasi.Hours()/24/7 < 4 {
		Durations = strconv.FormatFloat(durasi.Hours()/24/7, 'f', 0, 64) + " Weeks"
	} else if durasi.Hours()/24/30 < 12 {
		Durations = strconv.FormatFloat(durasi.Hours()/24/30, 'f', 0, 64) + " Months"
	} else {
		Durations = strconv.FormatFloat(durasi.Hours()/24/30/12, 'f', 0, 64) + " Years"
	}

	// if checked
	if c.FormValue("Python") != "" {
		tech1 = true
	}
	if c.FormValue("reactJS") != "" {
		tech2 = true
	}
	if c.FormValue("Javascript") != "" {
		tech3 = true
	}
	if c.FormValue("nodeJS") != "" {
		tech4 = true
	}

	for i := range dataProject {
		edit := &dataProject[id]
		if id == i {
			(*edit).Title = title
			// (*edit).Sdate = sDateFormat
			// (*edit).Edate = eDateFormat
			(*edit).Duration = Durations
			(*edit).Descript = desc
			(*edit).Tech1 = tech1
			(*edit).Tech2 = tech2
			(*edit).Tech3 = tech3
			(*edit).Tech4 = tech4
		}
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}
