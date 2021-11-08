#!/usr/bin/env ts-node
import 'reflect-metadata';
import os from 'os';
import path from 'path';
import { App, flag, cli, Command } from '@deepkit/app';
import { t } from '@deepkit/type';
import { CompressOptions } from './compress';
import { getDirectories, getFiles } from '../utils/fs';
import { runThreads, splitFiles } from '../utils/exec';
import { ReplayGainOptions } from './add-replay-gain';
import { EmbeddedBlocksOptions } from './remove-embedded-blocks';
import { TagsOptions } from './remove-tags';

@cli.controller('compress', {
	description: 'Compress',
})
export class Compression implements Command {
	// TODO: Replace this with DI when Logger class will work properly
	protected readonly logger = console;

	async execute(
		// Folders
		@flag.description('input folder') input: string,
		// Logging
		@(flag.description('logfile').default('flactwoflac.log')) logfile: string,
		// Processing
		@(flag
			.description('number of threads to use for processing')
			.default(os.cpus().length))
		@t.positive(true)
		threadCount: number,
		// Compression
		@(flag.description('compress if encoder tag is not present').default(true))
		ifEncoderTagIsNotPresent: boolean,
		// TODO: Check the following conditions
		// TODO: Skip if encoder tag is present
		// TODO: Skip if encoder tag matches
		// TODO: Skip if encoder version matches
		@(flag.description('compress if encoder version mismatches').default(true))
		ifEncoderVersionMismatches: boolean,
		@(flag.description('force compression').default(false))
		force: boolean,
		@(flag.description('add encoder tag with encoder settings').default(true))
		addEncoderTag: boolean,
		@(flag.description('encoder tag').default('COMMENT'))
		encoderTag: string,
		@(flag
			.description('encoder settings')
			.default(
				'flac --compression-level-8 --delete-input-file --no-padding --force --verify --warnings-as-errors --silent',
			))
		encoderSettings: string,
	) {
		const inputPath = path.resolve(input);

		const files = await getFiles(inputPath, {
			fileExtension: 'flac',
			recursive: true,
		});

		const filesPerThread = splitFiles(files, threadCount);

		const options: CompressOptions = {
			ifEncoderTagIsNotPresent,
			ifEncoderVersionMismatches,
			force,
			addEncoderTag,
			encoderTag,
			encoderSettings,
		};

		await runThreads(
			filesPerThread,
			options,
			path.join(__dirname, './compress.ts'),
			this.logger,
		);

		this.logger.log('Done!');
	}
}

@cli.controller('remove-embedded-blocks', {
	description: 'Remove embedded blocks',
})
export class Embedded implements Command {
	// TODO: Replace this with DI when Logger class will work properly
	protected readonly logger = console;

	async execute(
		// Folders
		@flag.description('input folder') input: string,
		// Logging
		@(flag.description('logfile').default('flactwoflac.log')) logfile: string,
		// Processing
		@(flag
			.description('number of threads to use for processing')
			.default(os.cpus().length))
		@t.positive(true)
		threadCount: number,
		// Embedded blocks
		@(flag.description('remove application blocks').default(true))
		removeApplicationBlocks: boolean,
		@(flag.description('remove cuesheet blocks').default(true))
		removeCuesheetBlocks: boolean,
		@(flag.description('remove padding blocks').default(true))
		removePaddingBlocks: boolean,
		@(flag.description('remove picture blocks').default(true))
		removePictureBlocks: boolean,
		@(flag.description('remove seek table blocks').default(true))
		removeSeektableBlocks: boolean,
	) {
		const inputPath = path.resolve(input);

		const files = await getFiles(inputPath, {
			fileExtension: 'flac',
			recursive: true,
		});

		const filesPerThread = splitFiles(files, threadCount);

		const options: EmbeddedBlocksOptions = {
			removeApplicationBlocks,
			removeCuesheetBlocks,
			removePaddingBlocks,
			removePictureBlocks,
			removeSeektableBlocks,
		};

		await runThreads(
			filesPerThread,
			options,
			path.join(__dirname, './remove-embedded-blocks.ts'),
			this.logger,
		);

		this.logger.log('Done!');
	}
}

@cli.controller('remove-tags', {
	description: 'Remove tags',
})
export class Tags implements Command {
	// TODO: Replace this with DI when Logger class will work properly
	protected readonly logger = console;

	async execute(
		// Folders
		@flag.description('input folder') input: string,
		// Logging
		@(flag.description('logfile').default('flactwoflac.log')) logfile: string,
		// Processing
		@(flag
			.description('number of threads to use for processing')
			.default(os.cpus().length))
		@t.positive(true)
		threadCount: number,
		// FLAC tags
		@(flag
			.description('Fill missings tags with a "[No <tag>]" in file')
			.default(true))
		fillMissingTags: boolean,
		@(flag
			.description('Vorbis tags to keep in files')
			.default([
				'ALBUM',
				'ALBUMARTIST',
				'ARTIST',
				'DISCNUMBER',
				'REPLAYGAIN_REFERENCE_LOUDNESS',
				'REPLAYGAIN_ALBUM_GAIN',
				'REPLAYGAIN_ALBUM_PEAK',
				'REPLAYGAIN_TRACK_GAIN',
				'REPLAYGAIN_TRACK_PEAK',
				'TITLE',
				'TRACKNUMBER',
				'TOTALDISCS',
				'TOTALTRACKS',
				'YEAR',
			]))
		@t.array(String)
		tagsToKeep: string[],
	) {
		const inputPath = path.resolve(input);

		const files = await getFiles(inputPath, {
			fileExtension: 'flac',
			recursive: true,
		});

		const filesPerThread = splitFiles(files, threadCount);

		const options: TagsOptions = {
			tagsToKeep,
			fillMissingTags,
		};

		await runThreads(
			filesPerThread,
			options,
			path.join(__dirname, './remove-tags.ts'),
			this.logger,
		);

		this.logger.log('Done!');
	}
}

@cli.controller('add-replay-gain', {
	description: 'Add ReplayGain',
})
export class ReplayGain implements Command {
	// TODO: Replace this with DI when Logger class will work properly
	protected readonly logger = console;

	async execute(
		// Folders
		@flag.description('input folder') input: string,
		// Logging
		@(flag.description('logfile').default('flactwoflac.log')) logfile: string,
		// Processing
		@(flag
			.description('number of threads to use for processing')
			.default(os.cpus().length))
		@t.positive(true)
		threadCount: number,
		// ReplayGain
		@(flag.description('force ReplayGain tags').default(false))
		force: boolean,
		@(flag
			.description('ReplayGain tags to look in files')
			.default([
				'REPLAYGAIN_REFERENCE_LOUDNESS',
				'REPLAYGAIN_TRACK_GAIN',
				'REPLAYGAIN_TRACK_PEAK',
				'REPLAYGAIN_ALBUM_GAIN',
				'REPLAYGAIN_ALBUM_PEAK',
			]))
		@t.array(String)
		tags: string[],
	) {
		const inputPath = path.resolve(input);

		const directories = await getDirectories(inputPath, {
			fileExtension: 'flac',
			recursive: true,
		});

		const filesPerThread = splitFiles(directories, threadCount);

		const options: ReplayGainOptions = {
			force,
			tags,
		};

		await runThreads(
			filesPerThread,
			options,
			path.join(__dirname, './add-replay-gain.ts'),
			this.logger,
		);

		this.logger.log('Done!');
	}
}

new App({
	controllers: [Compression, Embedded, Tags, ReplayGain],
}).run();
