#!/bin/sh

function gogo {
    echo _____________________________________________________ 
    go test github.com/afajl/ctrl...
}
    

function gitdiff { 
    prev=
    while true; do
        diff="$(git diff --no-color | sum)"
        if [ "$diff" != "$prev" ]; then
            gogo
        fi
        prev="$diff"
        sleep 1
    done
}

function inotify {
    while true; do
        inotifywait -r . -q -e modify -e create -e delete --exclude '^.*\.sw?'
        gogo
    done
}

if type -t inotifywait; then
    inotify
else
    gitdiff
fi
