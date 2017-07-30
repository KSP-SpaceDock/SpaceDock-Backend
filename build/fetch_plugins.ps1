$script:THIS_PATH = $myinvocation.mycommand.path
$script:BASE_DIR = split-path (resolve-path "$THIS_PATH/..") -Parent
$script:DIR_NAME = split-path $BASE_DIR -Leaf
$script:PLUGIN_FILE = $BASE_DIR + "/build/plugins.txt"
$script:SOURCE_FILE = "sdb.go"

Write-Host "Reading plugins.txt file..."

# Implement plugin stuff
Get-Content ($BASE_DIR + "/" + $SOURCE_FILE) | Foreach-Object {
    if ($_ -eq ")")  {
    # Add Lines before the selected pattern 
        if (Test-Path $PLUGIN_FILE) {
            Get-Content $PLUGIN_FILE | Foreach-Object {
                Write-Host "Fetching $_"
                & "$env:GOPATH/bin/glide.exe" get $_
                '    _ "' + $_ + '"'
            }
        }
    }
    $_ # send the current line to output
} | Set-Content ($BASE_DIR + "/build_" + $SOURCE_FILE)

# We are done
Write-Host "Plugins successfully applied. Please run go build on 'build_$SOURCE_FILE' to complete the build process."
