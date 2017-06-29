#!/usr/bin/env bash
######################################################
# Name:             processflac.sh
# Author:           Ludovic Delafontaine
# Creation:         29.06.2015
# Description:      Process FLACs files to add, edit or remove tags and/or compress the file and much more
# Documentation:    https://github.com/ludelafo/processflac.sh
######################################################

######################################################
# In case of problems, stops the script
set -e
set -u

######################################################
# Consts - Change here if needed
######################################################
REMOVE_UNWANTED_METADATA=true
DELETE_COVER_FROM_FILE=true
WRITE_COVER_TO_FILE=false
OPTIMIZE_COVER_FILE=true
CHECK_FOR_COVER=true
DELETE_COVER_FROM_FILESYSTEM=false
EXPORT_LYRICS=true
IMPORT_LYRICS=true
DELETE_LYRICS_FILE=true
ADD_COMMENT=false
COMPRESS_FLAC_FILE=true
FORCE_COMPRESS=false
ADD_REPLAYGAIN=true
MOVE_FILES=false
DELETE_UNWANTED_FILES=true # Only keeps FLACs, lyrics and cover.jpg by default

COMMENT="Share with love <3"
COVER_FILE="cover.jpg"
ENCODER="$(flac --version) --compression-level-8"

SOURCE_FOLDER="/media/0fbf8e50-3c72-4b0e-9810-b05c9fb9b7e9/Media/ToSort"
DESTINATION_FOLDER="/media/0fbf8e50-3c72-4b0e-9810-b05c9fb9b7e9/Media/ToSort"

TEMP_PATH="/dev/shm"
TEMPLATE_METADATA="$TEMP_PATH/template.txt"
TEMP_METADATA="$TEMP_PATH/metadata.txt"
CLEAN_METADATA="$TEMP_PATH/clean-metadata.txt"
LOG_FILE="/tmp/processflac.log"

METADATA_TO_KEEP=(
    ARTIST
    ALBUM
    DATE
    TITLE
    LANGUAGE
    GENRE
    TRACKNUMBER
    TRACKTOTAL
    DISCNUMBER
    DISCTOTAL
    ENCODER
    REPLAYGAIN_ALBUM_GAIN
    REPLAYGAIN_ALBUM_PEAK
    REPLAYGAIN_TRACK_GAIN
    REPLAYGAIN_TRACK_PEAK
)

######################################################
# Functions
######################################################
verifyFlac() {
    music="$1"
    flac --silent --test "$music"
}

removeCover() {
    music="$1"
    metaflac --remove --block-type=PICTURE "$music"
}

optimizeCover() {
    coverFile="$1"
    jpegoptim --quiet "$coverFile"
}

importCover() {
    music="$1"
    cover="$2"
    metaflac --import-picture-from="$cover" "$music"
}

deleteCoverFile() {
    coverFile="$1"
    rm $coverFile
}

exportMetadata() {
    music="$1"
    exportTo="$2"
    metaflac --export-tags-to="$exportTo" "$music"
}

cleanMetadata() {
    metadata="$1"
    templateMetadata="$2"
    output="$3"
    grep --word-regexp --file "$templateMetadata" "$metadata"  > "$output"
}

importMetadata() {
    metadata="$1"
    music="$2"
    metaflac --remove-all-tags --import-tags-from="$metadata" "$music"
}

deleteID3Comments() {
    music="$1"
    id3v2 --delete-all "$music" > /dev/null
}

exportLyrics() {
    music="$1"
    lyricsFile="$2"
    metaflac --show-tag LYRICS "$music" > "$lyricsFile"
}

cleanLyrics() {
    lyricsFile="$1"
    sed -i 's/LYRICS=//g' "$lyricsFile"
}

importLyrics() {
    lyricsFile="$1"
    music="$2"
    metaflac --set-tag-from-file="LYRICS=$lyricsFile" "$music"
}

deleteLyricsFile() {
    lyricsFile="$1"
    rm --force "$lyricsFile"
}

addComment() {
    comment="$1"
    music="$2"
    metaflac --set-tag="COMMENT=$comment" "$music"
}

removeEncoder() {
    music="$1"
    metaflac --remove-tag ENCODER "$music"
}

compressFlac() {
    music="$1"
    flac --silent --delete-input-file --compression-level-8 --force --verify "$music"
}

addEncoder() {
    encoder="$1"
    music="$2"
    metaflac --set-tag="ENCODER=$encoder" "$music"
}

removePadding() {
    music="$1"
    metaflac --dont-use-padding --remove --block-type=PADDING "$music"
}

addReplayGain() {
    album="$1"
    metaflac --add-replay-gain "$album"/*.flac
}

deleteUnwantedFiles() {
    directory="$1"
    find $directory ! -name "*.flac" ! -name "$COVER_FILE" ! -name "*.srt" -type f -exec rm {} +
}

moveFiles() {
    sourceDirectory="$1"
    destinationDirectory="$2"
    mkdir -p "$destinationDirectory"
    mv "$sourceDirectory" "$destinationDirectory"
}

log() {
    message="$1"
    echo "$message"
    echo "$message" >> "$LOG_FILE"
}

######################################################
# Script
######################################################
# Enable "Internal Field Separateur" to allow filenames with spaces
IFS=$'\n'

rm --force "$LOG_FILE"

printf "%s\n" "${METADATA_TO_KEEP[@]}" > "$TEMPLATE_METADATA"

for artist in $(find $SOURCE_FOLDER -mindepth 1 -maxdepth 1 -type d | sort); do

    echo "${artist##*/}"

    for album in $(find $artist -mindepth 1 -type d | sort); do

        echo "  ${album##*/}"

        error=false

        if $CHECK_FOR_COVER && [[ ! -f "$album/$COVER_FILE" ]]; then
            log "    The file '$album/$COVER_FILE' was not found !"
        else
            cover=true
        fi

        if $OPTIMIZE_COVER_FILE && [[ -f "$album/$COVER_FILE" ]]; then
            optimizeCover "$album/$COVER_FILE"
        fi

        if $DELETE_UNWANTED_FILES; then
            deleteUnwantedFiles "$album"
        fi

        for music in $(find $album -type f -name "*.flac" | sort); do

            echo "    ${music##*/}"

            if [[ $(verifyFlac "$music") != "" ]]; then
                log "The file '$music' is corrupted."
                error=true
                continue
            fi

            if $DELETE_COVER_FROM_FILE; then
                removeCover "$music"
            fi

            if $WRITE_COVER_TO_FILE && [[ -f "$album/$COVER_FILE" ]]; then
                importCover "$music" "$album/$COVER_FILE"
            fi

            if $EXPORT_LYRICS; then
                exportLyrics "$music" "${music%.*}.srt"
            fi

            if [[ -s "${music%.*}.srt" ]]; then
                deleteLyricsFile "${music%.*}.srt"
            else
                cleanLyrics "${music%.*}.srt"
            fi

            if $REMOVE_UNWANTED_METADATA; then

                exportMetadata "$music" "$TEMP_METADATA"
                cleanMetadata "$TEMP_METADATA" "$TEMPLATE_METADATA" "$CLEAN_METADATA"
                deleteID3Comments "$music"
                importMetadata "$CLEAN_METADATA" "$music"

            fi

            if $IMPORT_LYRICS && [[ -f "${music%.*}.srt" ]]; then
                importLyrics "${music%.*}.srt" "$music"
            fi

            if $DELETE_LYRICS_FILE && [[ -f "${music%.*}.srt" ]]; then
                deleteLyricsFile "${music%.*}.srt"
            fi

            if $ADD_COMMENT; then
                addComment "$COMMENT" "$music"
            fi

            if $FORCE_COMPRESS || ($COMPRESS_FLAC_FILE && [[ "$(metaflac --show-tag ENCODER $music)" != *"--compression-level-8" ]]); then

                removeEncoder "$music"

                if [[ $(compressFlac "$music" 2>&1) == "" ]]; then
                    addEncoder "$ENCODER" "$music"
                fi

            fi

            removePadding "$music"

        done

        if $DELETE_COVER_FROM_FILESYSTEM && [[ -f "$album/$COVER_FILE" ]]; then
            deleteCoverFile "$album/$COVER_FILE"
        fi

        if ! "$error"; then

            if $ADD_REPLAYGAIN; then
                addReplayGain "$album"
            fi

            if $MOVE_FILES; then
                moveFiles "$album" "$DESTINATION_FOLDER/${artist##*/}/"
            fi

            echo "  The album has been correctly processed."

        fi

    done

    sync
    wait

    if [[ ! "$(ls --almost-all $artist)" ]]; then
        rm -r "$artist"
    fi

done

if [[ -s $LOG_FILE ]]; then
    echo "Log not empty. Please check $LOG_FILE for more informations."
fi

echo "Everthing is done. Exiting program."
