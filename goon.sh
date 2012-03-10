#!/bin/sh

function gogo {
	echo
    echo __________________________________________________________________ 
    go test github.com/afajl/ctrl...
    #go vet github.com/afajl/ctrl...
}
    

function vimcheck { 
	touch /tmp/vim_saved_file
    while true; do
		if [ -e /tmp/vim_saved_file ]; then
			rm -f /tmp/vim_saved_file
            gogo
        fi
        sleep 0.5
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
    vimcheck
fi
