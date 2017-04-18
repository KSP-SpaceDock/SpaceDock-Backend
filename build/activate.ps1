# This file must be dot sourced from PoSh; you cannot run it
# directly. Do this: . ./activate.ps1

# FIXME: clean up unused vars.
$script:THIS_PATH = $myinvocation.mycommand.path
$script:BASE_DIR = split-path (resolve-path "$THIS_PATH/..") -Parent
$script:DIR_NAME = split-path $BASE_DIR -Leaf
$script:PLUGIN_FILE = $BASE_DIR + "/build/plugins.txt"

function global:deactivate ( [switch] $NonDestructive ){

    if ( test-path variable:_OLD_VIRTUAL_PATH ) {
        $env:PATH = $variable:_OLD_VIRTUAL_PATH
        remove-variable "_OLD_VIRTUAL_PATH" -scope global
    }

    if ( test-path variable:_OLD_GOPATH ) {
        $env:GOPATH = $variable:_OLD_GOPATH
        remove-variable "_OLD_GOPATH" -scope global
    }

    if ( test-path function:_old_virtual_prompt ) {
        $function:prompt = $function:_old_virtual_prompt
        remove-item function:\_old_virtual_prompt
    }

    if ($env:VIRTUAL_ENV) {
        $old_env = split-path $env:VIRTUAL_ENV -leaf
        remove-item env:VIRTUAL_ENV -erroraction silentlycontinue
    }

    if ( !$NonDestructive ) {
        # Self destruct!
        remove-item function:deactivate
        remove-item function:build
    }
}

function global:build($projectname) {
    # Vars
    $filename = $projectname + ".go"
    $binname = $projectname + ".exe"
    
    # Update deps
    go get -v -u github.com/KSP-SpaceDock/SpaceDock-Backend/install

    # Implement plugin stuff
    Get-Content ($BASE_DIR + "/" + $filename) | Foreach-Object {
        if ($_ -eq ")")  {
            # Add Lines before the selected pattern 
            if (Test-Path $PLUGIN_FILE) {
                Get-Content $PLUGIN_FILE | Foreach-Object {
                    '    _ "' + $_ + '"'
                    go get -v -u $_
                }
            }
        }        
        $_ # send the current line to output
    } | Set-Content ($BASE_DIR + "/build_" + $filename)
    
    # Build the binary
    go build -v -o $BASE_DIR/$binname $BASE_DIR/build_$filename
}

# unset irrelevant variables
deactivate -nondestructive

$VIRTUAL_ENV = $BASE_DIR
$env:VIRTUAL_ENV = $VIRTUAL_ENV

$global:_OLD_VIRTUAL_PATH = $env:PATH
$global:_OLD_GOPATH = $env:GOPATH
$env:PATH = "$env:VIRTUAL_ENV/build;" + $env:PATH
$env:GOPATH = "$env:VIRTUAL_ENV/vendor;" + $env:VIRTUAL_ENV
if (! $env:VIRTUAL_ENV_DISABLE_PROMPT) {
    function global:_old_virtual_prompt { "" }
    $function:_old_virtual_prompt = $function:prompt
    function global:prompt {
        # Add a prefix to the current prompt, but don't discard it.
        write-host "($(split-path $env:VIRTUAL_ENV -leaf)) " -nonewline
        & $function:_old_virtual_prompt
    }
}