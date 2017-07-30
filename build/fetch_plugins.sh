#!/bin/bash
SOURCE_FILE="sdb.go"
echo "Reading plugins.txt file..."

# Implement plugin stuff
FILE=""
if [ -e "$PWD/build_$filename" ] ; then
    rm $PWD/build_$SOURCE_FILE
fi
while IFS= read -r line; do
    if [ "$line" = ")" ]
    then
        if [ -e "$PWD/build/plugins.txt" ]
        then
            while IFS= read -r line2; do 
                FILE="$FILE    _ \"$line2\"\n"
                $GOPATH/bin/glide get $line2 || true
            done < $PWD/build/plugins.txt
        fi
    fi        
    FILE="$FILE$line\n"
done < $PWD/$SOURCE_FILE
printf "$FILE" > $PWD/build_$SOURCE_FILE 

# We are done
echo "Plugins successfully applied. Please run go build on 'build_$SOURCE_FILE' to complete the build process."
