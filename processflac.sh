#!/usr/bin/env bash
######################################################
# Name:				processflac.sh
# Author:			Ludovic Delafontaine
# Date:				29.06.2015 - 22.04.2016
# Description:		Process FLACs files to add, edit or remove tags and/or compress the file and much more
# Documentation:	https://github.com/ludelafo/processflac.sh
######################################################

######################################################
# In case of problems, stop the script
set -e
set -u

######################################################
# Consts - Change here if needed
######################################################
REPLACE_METADATA=true
REMOVE_COVER=true
IMPORT_COVER=false
COMPRESS_COVER_FILE=true
DELETE_COVER_FILE=true
EXPORT_LYRICS=true
IMPORT_LYRICS=true
DELETE_LYRICS_FILE=true
ADD_COMMENT=false
COMPRESS_FLAC_FILE=true
ADD_REPLAYGAIN=true
MOVE_FILES=true
DELETE_UNWANTED_FILES=true # Only keeps FLACs, lyrics and cover.jpg by default

COMMENT="Share with love <3"
COVER_FILE="cover.jpg"
ENCODER="$(flac --version) --compression-level-8"

SOURCE_FOLDER="/home/ludelafo/Music/Temp"
DESTINATION_FOLDER="/home/ludelafo/Music"

TEMPLATE_METADATA="/tmp/template.txt"
TEMP_METADATA="/tmp/metadata.txt"
CLEAN_METADATA="/tmp/clean-metadata.txt"
LOG_FILE="/tmp/workflac.logs"

METADATA_TO_KEEP=(
	ARTIST
	ALBUM
	MEDIA
	DATE
	TITLE
	LANGUAGE
	GENRE
	TRACKNUMBER
	TRACKTOTAL
	DISCNUMBER
	DISCTOTAL
	ENCODER
)

######################################################
# Functions
######################################################
verifyFlac() {
	music="$1"
	flac --warnings-as-errors --test --silent "$music"
}

removeCover() {
	music="$1"
	metaflac --remove --block-type=PICTURE "$music"
}

compressCover() {
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
	id3v2 --delete-all "$music" &> /dev/null
}

exportLyrics() {
	music="$1"
	lyricsFile="$2"
	metaflac --show-tag LYRICS "$music" > "$lyricsFile"
}

cleanLyrics() {
	lyricsFile="$1"
	sed -i -e 's/LYRICS=//g' "$lyricsFile"
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

logs() {
	message="$1"
	echo "$message" >> "$LOG_FILE"
}

######################################################
# Script
######################################################
# Enable "Internal Field Separateur" to allow filenames with spaces
IFS=$'\n'

rm --force "$LOG_FILE"

for artist in $(find $SOURCE_FOLDER -mindepth 1 -maxdepth 1 -type d | sort); do

	echo "${artist##*/}"

	for album in $(find $artist -mindepth 1 -type d | sort); do

		echo "  ${album##*/}"

		if $REPLACE_METADATA; then
			printf "%s\n" "${METADATA_TO_KEEP[@]}" > "$TEMPLATE_METADATA"
		fi

		if $COMPRESS_COVER_FILE && [[ -f "$album/$COVER_FILE" ]]; then
			compressCover "$album/$COVER_FILE"
		fi

		declare -i errors=0

		for music in $(find $album -type f -name "*.flac" | sort); do
			
			echo "    ${music##*/}"
				
			if [[ $(verifyFlac "$music" 2>&1) == "" ]]; then

				if $REMOVE_COVER; then
					removeCover "$music"
				fi

				if $IMPORT_COVER && [[ -f "$album/$COVER_FILE" ]]; then
					importCover "$music" "$album/$COVER_FILE"
				fi

				if $EXPORT_LYRICS; then
					exportLyrics "$music" "${music%.*}.srt"
					cleanLyrics "${music%.*}.srt"
				fi

				if $REPLACE_METADATA; then

					exportMetadata "$music" "$TEMP_METADATA"
					cleanMetadata "$TEMP_METADATA" "$TEMPLATE_METADATA" "$CLEAN_METADATA"
					deleteID3Comments "$music"
					importMetadata "$CLEAN_METADATA" "$music"
					
					if $IMPORT_LYRICS && [[ -f "${music%.*}.srt" ]]; then
						importLyrics "${music%.*}.srt" "$music"
					fi

				fi

				if $DELETE_LYRICS_FILE; then
					deleteLyricsFile "${music%.*}.srt"
				fi

				if $ADD_COMMENT; then
					addComment "$COMMENT" "$music"
				fi

				if $COMPRESS_FLAC_FILE && [[ "$(metaflac --show-tag ENCODER $music)" != "ENCODER=$ENCODER" ]]; then
					removeEncoder "$music"

					if [[ $(compressFlac "$music" 2>&1) == "" ]]; then
						addEncoder "$ENCODER" "$music"
					fi

				fi
				
			else
				echo "      The file has a problem. This will be logged..."
				logs "The file '$music' was not succesfully verified."
				errors=$((errors + 1))
			fi

		done
		
		if [[ $errors == 0 ]]; then

			if $DELETE_COVER_FILE; then
				deleteCoverFile "$album/$COVER_FILE"
			fi

			if $ADD_REPLAYGAIN; then
				addReplayGain "$album"
			fi
			
			if $DELETE_UNWANTED_FILES; then
				deleteUnwantedFiles "$album"
			fi

			if $MOVE_FILES; then
				moveFiles "$album" "$DESTINATION_FOLDER/${artist##*/}/"
			fi

			echo "  The album has been correctly processed."

		fi

	done

	if [[ ! "$(ls --almost-all $artist)" ]]; then
		rm -r "$artist"
	fi

done

sync

if [[ -f "$LOG_FILE" ]] && [[ -s $LOG_FILE ]]; then
	echo "Some errors were found. Please check $LOG_FILE for more informations."
fi

echo "Everthing is done. Exiting program."
