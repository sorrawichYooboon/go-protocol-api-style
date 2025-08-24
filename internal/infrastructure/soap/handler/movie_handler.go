package soaphandler

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/usecase"
)

const soapEnv = "http://schemas.xmlsoap.org/soap/envelope/"

type SOAPEnvelope struct {
	XMLName xml.Name    `xml:"Envelope"`
	Header  *SOAPHeader `xml:"Header,omitempty"`
	Body    SOAPBody    `xml:"Body"`
}
type SOAPHeader struct {
	Raw string `xml:",innerxml"`
}
type SOAPBody struct {
	Raw string `xml:",innerxml"`
}

type SOAPEnvelopeResponse struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Xmlns   string   `xml:"xmlns:soapenv,attr"`
	Body    SOAPBody `xml:"soapenv:Body"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"soapenv:Fault"`
	Code    string   `xml:"faultcode"`
	String  string   `xml:"faultstring"`
}

type GetMovieRequest struct {
	XMLName xml.Name `xml:"GetMovieRequest"`
	ID      int64    `xml:"id"`
}

type GetMovieResponse struct {
	XMLName     xml.Name `xml:"GetMovieResponse"`
	ID          int64    `xml:"id"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	ReleaseDate string   `xml:"releaseDate"`
}

type MovieSOAPHandler struct {
	movieUsecase usecase.MovieUsecase
}

func NewMovieSOAPHandler(movieUsecase usecase.MovieUsecase) *MovieSOAPHandler {
	return &MovieSOAPHandler{movieUsecase: movieUsecase}
}

func (h *MovieSOAPHandler) Handle(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		h.writeSOAPFault(c, "Client", "Only POST is allowed")
		return
	}
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		h.writeSOAPFault(c, "Client", "Could not read request body")
		return
	}

	var envelope SOAPEnvelope
	if err := xml.Unmarshal(body, &envelope); err != nil {
		h.writeSOAPFault(c, "Client", "Invalid SOAP envelope")
		return
	}

	action, actionBody, err := getSOAPAction(envelope.Body.Raw)
	if err != nil {
		h.writeSOAPFault(c, "Client", err.Error())
		return
	}

	switch action {
	case "GetMovieRequest":
		var req GetMovieRequest
		if err := xml.Unmarshal([]byte(actionBody), &req); err != nil {
			h.writeSOAPFault(c, "Client", "Invalid GetMovieRequest")
			return
		}
		if req.ID == 0 {
			h.writeSOAPFault(c, "Client", "id is required")
			return
		}
		h.processGetMovie(c, req)
	default:
		h.writeSOAPFault(c, "Client", "Unknown SOAP action: "+action)
	}
}

func (h *MovieSOAPHandler) processGetMovie(c *gin.Context, req GetMovieRequest) {
	movie, err := h.movieUsecase.GetMovieByID(req.ID)
	if err != nil || movie == nil {
		h.writeSOAPFault(c, "Server", "Movie not found")
		return
	}

	resp := GetMovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		ReleaseDate: movie.ReleaseDate,
	}
	out, _ := xml.Marshal(resp)
	envelope := SOAPEnvelopeResponse{
		Xmlns: soapEnv,
		Body:  SOAPBody{Raw: string(out)},
	}
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	xml.NewEncoder(&buf).Encode(envelope)
	c.Data(http.StatusOK, "text/xml; charset=utf-8", buf.Bytes())
}

func (h *MovieSOAPHandler) writeSOAPFault(c *gin.Context, code, msg string) {
	fault := SOAPFault{Code: code, String: msg}
	out, _ := xml.Marshal(fault)
	envelope := SOAPEnvelopeResponse{
		Xmlns: soapEnv,
		Body:  SOAPBody{Raw: string(out)},
	}
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	xml.NewEncoder(&buf).Encode(envelope)
	c.Data(http.StatusOK, "text/xml; charset=utf-8", buf.Bytes())
}

func getSOAPAction(body string) (string, string, error) {
	d := xml.NewDecoder(strings.NewReader(body))
	for {
		tok, err := d.Token()
		if err != nil {
			break
		}
		switch se := tok.(type) {
		case xml.StartElement:
			name := se.Name.Local
			var elemBuf bytes.Buffer
			e := xml.NewEncoder(&elemBuf)
			e.EncodeToken(se)
			depth := 1
			for depth > 0 {
				tok, err := d.Token()
				if err != nil {
					return "", "", err
				}
				switch tok.(type) {
				case xml.StartElement:
					depth++
				case xml.EndElement:
					depth--
				}
				e.EncodeToken(tok)
			}
			e.Flush()
			return name, elemBuf.String(), nil
		}
	}
	return "", "", errors.New("no SOAP action found")
}

func (h *MovieSOAPHandler) ServeWSDL(c *gin.Context) {
	wsdlBytes, err := os.ReadFile("internal/infrastructure/soap/wsdl/movie_service.wsdl")
	if err != nil {
		c.String(http.StatusInternalServerError, "WSDL not found")
		return
	}
	c.Data(http.StatusOK, "text/xml; charset=utf-8", wsdlBytes)
}
