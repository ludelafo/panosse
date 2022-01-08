import { workerData, parentPort } from 'worker_threads';
import { getFiles } from '../../utils/fs';
import {
	addReplayGain,
	addTag,
	getTag,
	removeTag,
} from '../../utils/metaflac';
import { CommonOptions } from '../../utils/types';

export type ReplayGainOptions = CommonOptions & {
	addReplayGainIfTagsAreNotPresent: boolean;
	addReplayGainIfTagMismatches: boolean;
	replayGainTags: string[];
	saveReplayGainSettingsInTag: boolean;
	replayGainTag: string;
	replayGainSettings: string;
};

const worker = async (
	threadNumber: number,
	flacDirectories: string[],
	options: ReplayGainOptions,
) => {
	const {
		dryRun,
		force,
		addReplayGainIfTagsAreNotPresent,
		addReplayGainIfTagMismatches,
		replayGainTags,
		saveReplayGainSettingsInTag,
		replayGainTag,
		replayGainSettings,
	} = options;

	for await (const flacDirectory of flacDirectories) {
		parentPort?.postMessage(`Thread ${threadNumber} - '${flacDirectory}'`);

		let needToAddReplayGain = false;
		let needToSaveReplayGainSettingsInTag = false;
		let allTagsPresent = true;

		const flacFiles = await getFiles(flacDirectory, {
			fileExtension: 'flac',
		});

		for await (const flacFile of flacFiles) {
			if (force) {
				parentPort?.postMessage(
					`Thread ${threadNumber} - '${flacDirectory}' - Force ReplayGain`,
				);

				needToAddReplayGain = true;
				needToSaveReplayGainSettingsInTag = true;
			} else {
				const replayGainTagContent = await getTag(flacFile, replayGainTag);

				for await (const tag of replayGainTags) {
					const tagContent = await getTag(flacFile, tag);

					if (!tagContent) {
						allTagsPresent = false;
						break;
					}
				}

				if (addReplayGainIfTagsAreNotPresent && !allTagsPresent) {
					parentPort?.postMessage(
						`Thread ${threadNumber} - '${flacFile}' - ReplayGain tags missing, need to add them`,
					);

					needToAddReplayGain = true;
					needToSaveReplayGainSettingsInTag = true;
				}

				if (addReplayGainIfTagMismatches && replayGainTagContent != replayGainSettings) {
					parentPort?.postMessage(
						`Thread ${threadNumber} - '${flacFile}' - ReplayGain tags mismatch, need to add them`,
					);

					needToAddReplayGain = true;
					needToSaveReplayGainSettingsInTag = true;
				}
			}

			if (saveReplayGainSettingsInTag && needToSaveReplayGainSettingsInTag) {
				parentPort?.postMessage(
					`Thread ${threadNumber} - '${flacFile}' - Saving ReplayGain settings in ReplayGain tag...`,
				);

				if (!dryRun) {
					await removeTag(flacFile, replayGainTag);
					await addTag(flacFile, replayGainTag, replayGainSettings);
				}
			}
		}

		if (needToAddReplayGain) {
			parentPort?.postMessage(
				`Thread ${threadNumber} - '${flacDirectory}' - Adding ReplayGain...`,
			);

			if (!dryRun) {
				await addReplayGain(flacFiles);
			}
		}
	}
};

worker(workerData.threadNumber, workerData.tasks, workerData.options);
