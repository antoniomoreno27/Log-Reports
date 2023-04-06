package main

import (
	"time"

	"github.com/antoniomoreno27/logs-LS/internal/controllers"
	"github.com/antoniomoreno27/logs-LS/internal/domain"
	"github.com/antoniomoreno27/logs-LS/internal/services/report"
)

func main() {
	const (
		cookie string = `_ga_28MSZRPJS9=GS1.1.1680705390.40.0.1680705390.0.0.0; security_storage=Fe26.2**f2b2b2ab4757d4c8885262906dfb23214e0a0f9a45fba672a0e0aabe6dcc04f5*6k5eC0wqyK45iQzYshhc0Q*5TkVQljGvNsS9Blv04vCd-Gy57lcxdGxdXd_dHBLaJeIiHU11yhvn3OK3bbCQz-4e3ebcEdkXPp2tRh183_bkuJTdHAKHaizQfqxz3m-wyl8GDcuwMJ0jUh8tJJFkcJm5mm1KL5VLhksVmyy2tJ1Kxxu0h-ppowyJs7QwiGFjy6VGsQpo7fI1BebRqDk3lgpZvrKj-hh_gvfKn1021vaUetCQ8DZ_koitpxFEMaeLOmCAVxanrFX593bDxu0ti7E**3afb5a553ca536fcb6894caa0ae9fd4c354019f220d3734c0d046b5a3f6e27d1*Z4p05jT-38ZStJztfn-M268-mmsro8Le9l6fef6rGOE; security_authentication=Fe26.2**6cb797e228660501c48526c2c57eecca444e7f23c120d56fe5f1215e5c716c40*D6OI4erx3Y0abwXeXfpqvg*im4mjeFEZFXK8mb7AkWGDsPFYBF6_zeXkYFBNNPoouyVP5pw6O5A7cnSnoy2vZQOnTacXK9QBf5AfxtuB5XkG__EZJ78sSlOQ8QE0Bh7gjnHwNOsw0GR60kS7rlJY_ETxI7QSPSySEeTkhRNZ04mPP_A3PLOUwLoTxWCJoODWnRqQvmXPB14lm7zI0mAhXz5_yg_1rSOSVIyFLIrQuSwi9pcDkTRa5i4rpUfWtS6Tdc9pxM-k0pyoAEuFf2MzBZPNo8ojd7p_TfSZYNZDtPcFYHWBmGi4QwHpWWKNL1D0dY1vgQ7F_FCcXmvfzEFbSD6ej6fkqgf_CiOOlrhRkcxey34taapdGyZEuZUVYcECGC5GiTuLGmCP24Zmp4wh32nmGWy8dgZrhGanbShUB1VnnebR2JZLFeI3_rwToQESQYf9rmKz4SMtEXX1S-_ksH8Hq1r4K5gd0r5G6WhzF6xKoWZm2C8d7nF3_C5PznDcebcrNcAT9kgT1NApUsJOi6VlsArz9ZZ4R9qaWD4QmbfeK27IyMxrXomgefenEvrd1H1gq4t9SpPRk8j4uztEV4ReISDB9dmjqkHG3bUiABFj9uP1Y2f0DH3UBfCN6xbqco**46d5bc0e2bbdb2eade4fdf9ac9103146236a239064d8613a4748947142d50d61*MvdzWNaHqAaFHBJNDJrINzcL5cQ5HQSMoKNljC371us`
	)
	timeWindowSample, err := time.ParseDuration("2h")
	if err != nil {
		panic("bat time window sample definition")
	}
	rs, err := report.New(
		cookie,
		domain.Report{
			Name:             "ParseError",
			Path:             "./reports",
			Query:            "Error retrieving iterator next data",
			TimeWindowSample: timeWindowSample,
			StartDate:        time.Now().AddDate(0, 0, -9),
			FinishDate:       time.Now(),
		},
	)
	if err != nil {
		panic(err)
	}
	reportController := controllers.Report{
		ReportService: rs,
	}
	reportController.Create()
}
