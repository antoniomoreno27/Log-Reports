package main

import (
	"time"

	"github.com/antoniomoreno27/logs-LS/internal/controllers"
	"github.com/antoniomoreno27/logs-LS/internal/domain"
	"github.com/antoniomoreno27/logs-LS/internal/services/logger"
	"github.com/antoniomoreno27/logs-LS/internal/services/report"
)

func main() {
	const (
		cookie string = `_ga_28MSZRPJS9=GS1.1.1680804503.41.0.1680804503.0.0.0; security_storage=Fe26.2**c1c44afb2478714e188d490cbf9ad3b88571534b0afb5fe7860e3459ccc2c72c*CY_ayTRSgwziZQ70BaZu8w*C3ckN2Yt3yO4F5-IrkwEfz53dJEx_hJ72QbXqbUIiDNSa4azTnMVrdsYgYe7xZA-B590bhZMW7Q2WzDMKQpZmr0h_Zir9juJtJ7Jpy6hepMfcyFWhhgCJF5M02D3vIcOtXbG4FNj3cgynk1jdzqh7SrsOMm8SqiPsMwaDPHSWxir7bT8aNNGGzNQPmLaG5MMq3SfYj0GVOpXGpHa166h2CTreBnzgb11n2LsbrZZrRL45-X6wYeU9ZWbnzRLMUpz**0ad0dd21dbaf519e39296b71c8632042e155792ea84dfdee5bc776252072fd39*YGpKv6b02g1OnPQuCf9crCDbcPuM1jyUuxSiMHaT2EI; security_authentication=Fe26.2**f95a18516973eceb899d855a88cc3635e0c4a5256546a1a2aeefe37d4176c3bd*5RsCT7WXUfqfNBqcF6p6lQ*v8n5bLfc7GDBHK3f79zeiDaQUAiJ3A3ngkzbYvjiuGx_AQmmOqCe81JzSXheMsxulb4SgyGtF3PuJdR1TX0OEzV4CcMxr45DRi_zO820gyCIOFh6wRjjS3qEldJqQM420aszK3tRDR56RKxSg99FspRH_qy_vnM5mkl2tfp6JV6B1gMoawLwiCPxcbQ-pUGq6E3Fs6oDtm0YWCo9auxg3D_lceLx7ezyb8DgAnJYhrityhJsC48pfqaaFmNYAyakw4Vs-9J03sOvpR0esIj-GG2Kj9j95vamBwGoTxh4EoBjVExhKgs27vvnXM--Gh8WCe-a-IEQhndsLyBahDwL3hScyjG-fDCFJBtYt_k2nuvkYOoPNpsHf3LqdeIImrKKsQvVMV4Ft-84Z8oEmegL7Q6UJVQalcmM_Ep_MovUAyoxy0MTVbIMjM9titKH1_c_y-JtJHSYzExTDF71dQuGumH6lcjkQKEwNSH5-0uo0eGFqRHiztzo28tesEmmWZDXoo3Iky-028eSylNrz1KwyplDqpq7BWWN_ynyqCIs0DnKcJeHfO2IcHlFXGpd9U1r8tWciJysynPzzHuG5NdOZC8zk4JoDgAz_h83xNpbbm0**6051e931773fcce91439a203b522ff00581f93c031cb2e6150f6cee0851c965e*Rv1oZYAdIEiHg6vkFst9ByLCXcddVsXa96FDSh_BVEU`
	)
	logger.Warnf("Configurating Report")
	timeWindowSample, err := time.ParseDuration("2h")
	if err != nil {
		logger.Panic("bat time window sample definition")
	}
	rs, err := report.New(
		cookie,
		domain.Report{
			Name:             "parse_error.csv",
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
	logger.Warnf("Creating report")
	reportController.Create()
}
