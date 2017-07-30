$script:THIS_PATH = $myinvocation.mycommand.path
$script:BASE_DIR = split-path (resolve-path "$THIS_PATH/..") -Parent
$GLIDE_VERSION = "v0.12.3"
$GLIDE_FILE = "glide-$GLIDE_VERSION-windows-amd64.zip"

New-Item -Path "$BASE_DIR" -Name "tmp" -ItemType "directory"
Write-Host "Downloading the glide package manager..."
[System.Net.ServicePointManager]::ServerCertificateValidationCallback = { $true } # I hate cert stores
$client = New-Object System.Net.WebClient
$client.DownloadFile("https://github.com/Masterminds/glide/releases/download/$GLIDE_VERSION/$GLIDE_FILE", "$BASE_DIR/tmp/$GLIDE_FILE")
Write-Host "Extracting the glide executable..."
Expand-Archive "$BASE_DIR/tmp/$GLIDE_FILE" -DestinationPath "$BASE_DIR/tmp"    
New-Item -Path "$env:GOPATH" -Name "bin" -ItemType "directory"
Move-Item -path "$BASE_DIR/tmp/windows-amd64/glide.exe" -destination "$env:GOPATH/bin/glide.exe"
Remove-Item -Force -Recurse "$BASE_DIR/tmp"
Write-Host "Done."
