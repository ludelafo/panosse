import { workerData, parentPort } from 'worker_threads';
import { compressFile, getEncoderVersion } from '../../utils/flac';
import { addTag, getTag, getVendorTag, removeTag } from '../../utils/metaflac';
import { CommonOptions } from '../../utils/types';

export type CompressionOptions = CommonOptions & {
	compressIfTagIsNotPresent: boolean;
	compressIfTagMismatches: boolean;
	compressIfVersionMismatches: boolean;
	saveEncoderSettingsInTag: boolean;
	encoderTag: string;
	encoderSettings: string;
};

const worker = async (
	threadNumber: number,
	flacFiles: string[],
	options: CompressionOptions,
) => {
	const {
		dryRun,
		force,
		compressIfTagIsNotPresent,
		compressIfTagMismatches,
		compressIfVersionMismatches,
		saveEncoderSettingsInTag,
		encoderTag,
		encoderSettings,
	} = options;

	for await (const flacFile of flacFiles) {
		parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}'`);

		let needToCompress = false;
		let needToSaveEncoderSettingsInTag = false;

		const encoderTagContent = await getTag(flacFile, encoderTag);
		const fileEncoderVersion = await getVendorTag(flacFile);
		const hostEncoderVersion = await getEncoderVersion();

		if (force) {
			parentPort?.postMessage(
				`Thread ${threadNumber} - '${flacFile}' - Force compression`,
			);

			needToCompress = true;
			needToSaveEncoderSettingsInTag = true;
		} else {
			if (compressIfTagIsNotPresent && !encoderTagContent) {
				parentPort?.postMessage(
					`Thread ${threadNumber} - '${flacFile}' - Compress because encoder tag is not present`,
				);

				needToCompress = true;
				needToSaveEncoderSettingsInTag = true;
			}

			if (compressIfTagMismatches && encoderTagContent !== encoderSettings) {
				parentPort?.postMessage(
					`Thread ${threadNumber} - '${flacFile}' - Compress because encoder tag mistaches`,
				);

				needToCompress = true;
				needToSaveEncoderSettingsInTag = true;
			}

			if (
				compressIfVersionMismatches &&
				fileEncoderVersion != hostEncoderVersion
			) {
				parentPort?.postMessage(
					`Thread ${threadNumber} - '${flacFile}' - Compress because encoder version mismatches`,
				);

				needToCompress = true;
				needToSaveEncoderSettingsInTag = true;
			}
		}

		if (needToCompress) {
			parentPort?.postMessage(
				`Thread ${threadNumber} - '${flacFile}' - Compressing...`,
			);

			if (!dryRun) {
				await compressFile(flacFile, encoderSettings);
			}
		}

		if (saveEncoderSettingsInTag && needToSaveEncoderSettingsInTag || force) {
			parentPort?.postMessage(
				`Thread ${threadNumber} - '${flacFile}' - Saving encoder settings in encoder tag...`,
			);

			if (!dryRun) {
				await removeTag(flacFile, encoderTag);
				await addTag(flacFile, encoderTag, encoderSettings);
			}
		}
	}
};

worker(workerData.threadNumber, workerData.tasks, workerData.options);
