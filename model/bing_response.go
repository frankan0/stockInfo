package model

type BingResponse struct {
	Images []BingImage
}

type BingImage struct {

	Startdate string
	Enddate string
	Url string
	Urlbase string
	Copyright string
	Title string

}
