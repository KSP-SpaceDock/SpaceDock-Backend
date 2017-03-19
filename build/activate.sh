# This file must be used with "source bin/activate" *from bash*
# you cannot run it directly

deactivate () {

    # reset old environment variables
    # ! [ -z ${VAR+_} ] returns true if VAR is declared at all
    if ! [ -z "${_OLD_VIRTUAL_PATH+_}" ] ; then
        PATH="$_OLD_VIRTUAL_PATH"
        export PATH
        unset _OLD_VIRTUAL_PATH
    fi
    if ! [ -z "${_OLD_GOPATH+_}" ] ; then
        GOPATH="$_OLD_GOPATH"
        export GOPATH
        unset _OLD_GOPATH
    fi

    # This should detect bash and zsh, which have a hash command that must
    # be called to get it to forget past commands.  Without forgetting
    # past commands the $PATH changes we made may not be respected
    if [ -n "${BASH-}" ] || [ -n "${ZSH_VERSION-}" ] ; then
        hash -r 2>/dev/null
    fi

    if ! [ -z "${_OLD_VIRTUAL_PS1+_}" ] ; then
        PS1="$_OLD_VIRTUAL_PS1"
        export PS1
        unset _OLD_VIRTUAL_PS1
    fi

    unset VIRTUAL_ENV
    if [ ! "${1-}" = "nondestructive" ] ; then
    # Self destruct!
        unset -f deactivate
        unset -f build
    fi
}

build () {
    # Vars
    $filename = $1 + ".go"
    $binname = $1
    
    # Update deps
    go get -u ./...

    # Implement plugin stuff
    rm $VIRTUAL_ENV/build_$filename
    filelines=`cat $VIRTUAL_ENV/$filename`
    for line in $filelines ; do
        if ["$line" = ")"]
        then
            if [-e "$VIRTUAL_ENV/build/plugins.txt"]
            then
                filelines2=`cat $VIRTUAL_ENV/build/plugins.txt`
                for line2 in $filelines ; do        
                    echo $line >> $VIRTUAL_ENV/build_$filename                    
                done
            fi
        fi        
        echo $line >> $VIRTUAL_ENV/build_$filename
    done
    
    # Build the binary
    go build -v -o $VIRTUAL_ENV/build/$binname $VIRTUAL_ENV/build_$filename
}

# unset irrelevant variables
deactivate nondestructive

VIRTUAL_ENV="$(dirname "$PWD")"
export VIRTUAL_ENV

_OLD_VIRTUAL_PATH="$PATH"
PATH="$VIRTUAL_ENV/build:$PATH"
export PATH

_OLD_GOPATH="$GOPATH"
GOPATH="$VIRTUAL_ENV/vendor:$VIRTUAL_ENV"
export GOPATH

if [ -z "${VIRTUAL_ENV_DISABLE_PROMPT-}" ] ; then
    _OLD_VIRTUAL_PS1="$PS1"
    if [ "x" != x ] ; then
        PS1="$PS1"
    else
        PS1="(`basename \"$VIRTUAL_ENV\"`) $PS1"
    fi
    export PS1
fi

# This should detect bash and zsh, which have a hash command that must
# be called to get it to forget past commands.  Without forgetting
# past commands the $PATH changes we made may not be respected
if [ -n "${BASH-}" ] || [ -n "${ZSH_VERSION-}" ] ; then
    hash -r 2>/dev/null
fi
