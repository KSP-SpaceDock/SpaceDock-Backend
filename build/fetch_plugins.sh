#!/bin/bash
SOURCE_FILE="sdb.go"
echo "Reading plugins.txt file..."

# Implement plugin stuff
if [ -e "$PWD/build_$filename" ] ; then
    rm $PWD/build_$SOURCE_FILE
fi
touch $PWD/build_$SOURCE_FILE
while IFS= read -r line; do
    if [ "$line" = ")" ]
    then
        if [ -e "$PWD/build/plugins.txt" ]
        then
            while IFS= read -r line2; do 
                printf "    _ \"$line2\"\n" >> $PWD/build_$SOURCE_FILE
                glide get $line2 || true
            done < $PWD/build/plugins.txt
        fi
    fi        
    printf "$line\n" >> $PWD/build_$SOURCE_FILE
done < $PWD/$SOURCE_FILE 

# We are done
echo "Plugins successfully applied. Please run go build on 'build_$SOURCE_FILE' to complete the build process."
