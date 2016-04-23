# processflac.sh
A Bash script to process FLAC files such as removing metadata, compressing, adding ReplayGain and much more

## Processing
This script can do the following processings:

- Replace metadata based on a whitelist
- Remove ID3 tags from FLAC files (ID3 found in a FLAC file are considered as a corrupted FLAC file)
- Remove cover from FLAC files
- Import existing cover to FLAC files (the cover file must be in the album's folder)
- Compress cover file using jpegoptim
- Export lyrics contained in FLAC to .srt files
- Import lyrics to FLAC files
- Delete lyrics from hard drive
- Add comment to FLAC files
- Compress FLAC files
- Add ReplayGain to FLAC files
- Move files to a destination folder
- Delete unwanted files from folder (such as .cue, .nfo, and so)

## Configuration
Just open the script file and change the constants to process the files as you want

Note: Your music collection must be structured as follow:

    artists/
    	albums/
    		FLAC files
    		cover file

## More informations
- The compress process is only made once. When a successfully compression is done, the script saves the current encoder used to process the file directly in the FLAC file. So when you relaunch the script, it will only try to compress files that were not already compressed (useful when having a huge collection)

## Dependencies
This script depends on the following commands:

- bash
- flac
- metaflac
- jpegoptim
- id3v2
- sed
- grep

## Todo / known possible improvements
- Export lyrics only if present in FLAC file
- Add the possibility to have a whitelist of files to keep instead of modifing the `find` command
- Change the structure for searching musics
- Make a better documentation (supported metadata, and much more...)

# Sources (unformated for the moment)
The following sources were used to make this script. Thanks to them for the help !

- http://www.hydrogenaud.io/forums/index.php?showtopic=84634
- http://stackoverflow.com/questions/3810709/how-to-evaluate-a-boolean-variable-in-an-if-block-in-bash
- http://stackoverflow.com/questions/9755068/bash-list-and-sort-files-and-their-sizes-and-by-name-and-size
- http://stackoverflow.com/questions/9011233/for-files-in-directory-only-echo-filename-no-path
- http://stackoverflow.com/questions/16365938/create-list-of-all-files-in-every-subdirectories-in-bash
- http://www.xiph.org/vorbis/doc/v-comment.html
- http://isrc.ifpi.org/en/
- https://wiki.xiph.org/VorbisComment
- http://age.hobba.nl/audio/mirroredpages/ogg-tagging.html
- https://wiki.xiph.org/Field_names
- http://www.legroom.net/2009/05/09/ogg-vorbis-and-flac-comment-field-recommendations
- http://stackoverflow.com/questions/20243467/write-bash-array-to-file-with-newlines
- http://www.unix.com/shell-programming-and-scripting/191217-using-whitelist-file-remove-entries.html
- https://hydrogenaud.io/index.php/topic,61371.0.html
- https://sourceforge.net/p/flac/feature-requests/98/
- https://code.google.com/archive/p/libkate/wikis/CreatingKateStreams.wiki
- http://unix.stackexchange.com/questions/153862/remove-all-files-directories-except-for-one-file
- http://www.unix.com/shell-programming-and-scripting/191217-using-whitelist-file-remove-entries.html
- http://stackoverflow.com/questions/7359527/removing-trailing-starting-newlines-with-sed-awk-tr-and-friends
- https://hydrogenaud.io/index.php/topic,36283.0.html
- http://stackoverflow.com/questions/2292847/how-to-silence-output-in-a-bash-script
- http://www.davidpashley.com/articles/writing-robust-shell-scripts/
- http://askubuntu.com/questions/385528/how-to-increment-a-variable-in-bash
- http://stackoverflow.com/questions/17830326/ignoring-specific-errors-on-shell-script
- http://superuser.com/questions/352289/bash-scripting-test-for-empty-directory
- http://unix.stackexchange.com/questions/156534/bash-script-error-with-strings-with-paths-that-have-spaces-and-wildcards
- http://www.bobulous.org.uk/misc/Replay-Gain-in-Linux.html
- https://bash.cyberciti.biz/guide/Create_an_integer_variable
- http://stackoverflow.com/questions/965053/extract-filename-and-extension-in-bash