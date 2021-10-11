import { workerData, parentPort } from 'worker_threads';
import { compressFile, getEncoderVersion } from '../utils/flac';
import { getFiles } from '../utils/fs';
import {
	addReplayGain,
	addTag,
	getTag,
	getVendorTag,
	removeTag,
} from '../utils/metaflac';

export type ReplayGainOptions = {
	tags: string[];
	force: boolean;
};

const worker = async (
	threadNumber: number,
	flacDirectories: string[],
	options: ReplayGainOptions,
) => {
	const { force, tags } = options;

	for await (const flacDirectory of flacDirectories) {
		parentPort?.postMessage(`Thread ${threadNumber} - '${flacDirectory}'`);

		let needToAddReplayGain = false;
		let allTagsPresent = true;

		const flacFiles = await getFiles(flacDirectory, {
			fileExtension: 'flac',
		});

		if (force) {
			parentPort?.postMessage(
				`Thread ${threadNumber} - '${flacDirectory}' - Force ReplayGain`,
			);

			needToAddReplayGain = true;
		} else {
			for await (const flacFile of flacFiles) {
				allTagsPresent = true;

				for await (const tag of tags) {
					const tagContent = await getTag(flacFile, tag);

					if (!tagContent) {
						parentPort?.postMessage(
							`Thread ${threadNumber} - '${flacFile}' - ReplayGain tag missing, need to add them`,
						);

						allTagsPresent = false;
						break;
					}
				}

				if (!allTagsPresent) {
					needToAddReplayGain = true;
					break;
				}
			}
		}

		if (needToAddReplayGain) {
			parentPort?.postMessage(
				`Thread ${threadNumber} - '${flacDirectory}' - Adding ReplayGain...`,
			);

			await addReplayGain(flacFiles);
		}
	}
};

worker(workerData.threadNumber, workerData.files, workerData.options);
