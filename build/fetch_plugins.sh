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
            for line2 in `cat $PWD/build/plugins.txt`; do
                echo $line2
                FILE="$FILE    _ \"$line2\"\n"
                glide get $line2 || true
            done
        fi
    fi        
    FILE="$FILE$line\n"
done < $PWD/$SOURCE_FILE
printf "$FILE" > $PWD/build_$SOURCE_FILE 

# We are done
echo "Plugins successfully applied. Please run go build on 'build_$SOURCE_FILE' to complete the build process."
