package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"k8s.io/api/admission/v1beta1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func handler(c echo.Context) error {
	fmt.Println("got request...")

	adrev := v1beta1.AdmissionReview{}

	if err := c.Bind(&adrev); err != nil {
		return responseBadReview(c, adrev, err.Error())
	}

	adrev.Response = &v1beta1.AdmissionResponse{
		UID:     adrev.Request.UID,
		Allowed: true,
	}
	if adrev.Request.Kind.Kind != "Service" {
		return c.JSON(http.StatusOK, adrev)
	}

	service := core.Service{}
	err := json.Unmarshal(adrev.Request.Object.Raw, &service)
	if err != nil {
		return responseBadReview(c, adrev, err.Error())
	}

	for _, port := range service.Spec.Ports {
		name := strings.ToLower(port.Name)
		if port.Port == 443 && !strings.HasPrefix(name, "https") {
			adrev.Response = &v1beta1.AdmissionResponse{
				UID:     adrev.Request.UID,
				Allowed: false,
				Result: &meta.Status{
					Message: fmt.Sprintf(
						`invalid service %s: port=443 but name="%s" (not prefix with https)`,
						service.Name, name),
				},
			}

			break
		}
	}

	return c.JSON(http.StatusOK, adrev)
}

func responseBadReview(c echo.Context, adrev v1beta1.AdmissionReview, msg string) error {
	adrev.Response = &v1beta1.AdmissionResponse{
		Allowed: false,
		Result: &meta.Status{
			Message: msg,
			Code:    http.StatusBadRequest,
		},
	}

	return c.JSON(http.StatusBadRequest, adrev)
}
