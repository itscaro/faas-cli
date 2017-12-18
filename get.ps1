$version = (Invoke-WebRequest "https://api.github.com/repos/openfaas/faas-cli/releases/latest" | ConvertFrom-Json)[0].tag_name

(New-Object System.Net.WebClient).DownloadFile("https://github.com/openfaas/faas-cli/releases/download/$version/faas-cli.exe", "faas-cli.exe")
