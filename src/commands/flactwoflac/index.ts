import {Command, flags} from '@oclif/command';
import { cpus as threads } from 'os';
import { resolve as getPath, join as joinPath } from 'path';
import { runThreads, splitTasksPerThread } from '../../utils/exec';
import { getDirectories, getFiles } from '../../utils/fs';
import { ReplayGainOptions } from './add-replay-gain';
import { TagsOptions } from './clean-tags';
import { CompressionOptions } from './compress';
import { EmbeddedBlocksOptions } from './remove-embedded-blocks';

export default class FlacTwoFlac extends Command {
	static description = 'Process FLAC files to FLAC files'

	static flags = {
		// Display the help for the command
		help: flags.help({char: 'h'}),
		// Processing
		threadCount: flags.integer({description: 'number of threads to use for processing', default: threads().length}),
		force: flags.boolean({description: 'force command'}),
		dryRun: flags.boolean({description: 'run without actually doing anything'}),
		// Compression
		compress: flags.boolean({description: 'compress FLAC files'}),
		compressIfTagIsNotPresent: flags.boolean({description: 'compress if tag is not present', default: true}),
		compressIfTagMismatches: flags.boolean({description: 'compress if tag mismatches', default: true}),
		compressIfVersionMismatches: flags.boolean({description: 'compress if version mismatches', default: true}),
		saveEncoderSettingsInTag: flags.boolean({description: 'add encoder tag with encoder settings', default: true}),
		encoderTag: flags.string({description: 'encoder tag name', default: 'ENCODER_SETTINGS'}),
		encoderSettings: flags.string({description: 'encoder settings', default: 'flac --compression-level-8 --delete-input-file --no-padding --force --verify --warnings-as-errors --silent'}),
		// Embedded blocks
		removeEmbeddedBlocks: flags.boolean({description: 'remove embedded blocks from FLAC files'}),
		removeApplicationBlocks: flags.boolean({description: 'remove application blocks', default: true}),
		removeCuesheetBlocks: flags.boolean({description: 'remove cuesheet blocks', default: true}),
		removePaddingBlocks: flags.boolean({description: 'remove padding blocks', default: true}),
		removePictureBlocks: flags.boolean({description: 'remove picture blocks', default: true}),
		removeSeektableBlocks: flags.boolean({description: 'remove seek table blocks', default: true}),
		// Tags
		cleanTags: flags.boolean({description: 'clean tags from FLAC files'}),
		fillMissingTags: flags.boolean({description: 'Fill missings tags with a "[No <tag>]" in file', default: true}),
		tagsToKeep: flags.string({ multiple: true, default: [
			'ALBUM',
			'ALBUMARTIST',
			'ARTIST',
			'COMMENT',
			'DISCNUMBER',
			'ENCODER_SETTINGS',
			'GENRE',
			'REPLAYGAIN_SETTINGS',
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
		]}),
		// ReplayGain
		addReplayGain: flags.boolean({description: 'add ReplayGain to FLAC files'}),
		addReplayGainIfTagsAreNotPresent: flags.boolean({description: 'add ReplayGain if tags are not present', default: true}),
		addReplayGainIfTagMismatches: flags.boolean({description: 'add ReplayGain if tag mismatches', default: true}),
		replayGainTags: flags.string({ multiple: true, default: [
			'REPLAYGAIN_REFERENCE_LOUDNESS',
			'REPLAYGAIN_TRACK_GAIN',
			'REPLAYGAIN_TRACK_PEAK',
			'REPLAYGAIN_ALBUM_GAIN',
			'REPLAYGAIN_ALBUM_PEAK',
		]}),
		saveReplayGainSettingsInTag: flags.boolean({description: 'add ReplayGain tag with ReplayGain settings', default: true}),
		replayGainTag: flags.string({description: 'ReplayGain tag name', default: 'REPLAYGAIN_SETTINGS'}),
		replayGainSettings: flags.string({description: 'ReplayGain settings', default: 'metaflac --add-replay-gain'}),
	}

	static args = [{name: 'input', description: 'input folder', required: true}]

	async run() {
		const {args, flags} = this.parse(FlacTwoFlac);

		const { dryRun, force, compress, removeEmbeddedBlocks, addReplayGain, cleanTags } = flags;

		const inputPath = getPath(args.input);

		const files = await getFiles(inputPath, {
			fileExtension: 'flac',
			recursive: true,
		});

		const directories = await getDirectories(inputPath, {
			fileExtension: 'flac',
			recursive: true,
		});

		const filesPerThread = splitTasksPerThread(files, flags.threadCount);

		const directoriesPerThread = splitTasksPerThread(directories, flags.threadCount);

		if (compress) {
			const {
				compressIfTagIsNotPresent,
				compressIfTagMismatches,
				compressIfVersionMismatches,
				saveEncoderSettingsInTag,
				encoderTag,
				encoderSettings,
			} = flags;

			const options: CompressionOptions = {
				dryRun,
				force,
				compressIfTagIsNotPresent,
				compressIfTagMismatches,
				compressIfVersionMismatches,
				saveEncoderSettingsInTag,
				encoderTag,
				encoderSettings,
			};

			await runThreads(
				filesPerThread,
				options,
				joinPath(__dirname, './compress.ts'),
				this.log,
			);
		}

		if (addReplayGain) {
			const {
				addReplayGainIfTagsAreNotPresent,
				addReplayGainIfTagMismatches,
				replayGainTags,
				saveReplayGainSettingsInTag,
				replayGainTag,
				replayGainSettings,
			} = flags;

			const options: ReplayGainOptions = {
				dryRun,
				force,
				addReplayGainIfTagsAreNotPresent,
				addReplayGainIfTagMismatches,
				replayGainTags,
				saveReplayGainSettingsInTag,
				replayGainTag,
				replayGainSettings,
			};

			await runThreads(
				directoriesPerThread,
				options,
				joinPath(__dirname, './add-replay-gain.ts'),
				this.log,
			);
		}

		if (cleanTags) {
			const {
				tagsToKeep,
				fillMissingTags,
			} = flags;

			const options: TagsOptions = {
				dryRun,
				force,
				tagsToKeep,
				fillMissingTags,
			};

			await runThreads(
				filesPerThread,
				options,
				joinPath(__dirname, './clean-tags.ts'),
				this.log,
			);
		}

		if (removeEmbeddedBlocks) {
			const {
				removeApplicationBlocks,
				removeCuesheetBlocks,
				removePaddingBlocks,
				removePictureBlocks,
				removeSeektableBlocks,
			} = flags;

			const options: EmbeddedBlocksOptions = {
				dryRun,
				force,
				removeApplicationBlocks,
				removeCuesheetBlocks,
				removePaddingBlocks,
				removePictureBlocks,
				removeSeektableBlocks,
			};

			await runThreads(
				filesPerThread,
				options,
				joinPath(__dirname, './remove-embedded-blocks.ts'),
				this.log,
			);
		}
	}
}
