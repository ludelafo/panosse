import { workerData, parentPort } from 'worker_threads';
import { compressFile, getEncoderVersion } from '../utils/flac';
import { addTag, getTag, getVendorTag, removeTag } from '../utils/metaflac';

export type CompressOptions = {
	ifEncoderTagIsNotPresent: boolean;
	ifEncoderVersionMismatches: boolean;
	force: boolean;
	addEncoderTag: boolean;
	encoderTag: string;
	encoderSettings: string;
};

const worker = async (
	threadNumber: number,
	flacFiles: string[],
	options: CompressOptions,
) => {
	const {
		ifEncoderTagIsNotPresent,
		ifEncoderVersionMismatches,
		force,
		addEncoderTag,
		encoderTag,
		encoderSettings,
	} = options;

	for await (const flacFile of flacFiles) {
		parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}'`);

		let needToCompress = false;
		let needToAddEncoderTag = false;

		const encoderTagContent = await getTag(flacFile, encoderTag);
		const fileEncoderVersion = await getVendorTag(flacFile);
		const hostEncoderVersion = await getEncoderVersion();

		if (force) {
			parentPort?.postMessage(
				`Thread ${threadNumber} - '${flacFile}' - Force compression`,
			);

			needToCompress = true;
			needToAddEncoderTag = true;
		} else {
			if (ifEncoderTagIsNotPresent && !encoderTagContent) {
				parentPort?.postMessage(
					`Thread ${threadNumber} - '${flacFile}' - Compress because encoder tag is not present`,
				);

				needToCompress = true;
				needToAddEncoderTag = true;
			}

			if (
				ifEncoderVersionMismatches &&
				fileEncoderVersion != hostEncoderVersion
			) {
				parentPort?.postMessage(
					`Thread ${threadNumber} - '${flacFile}' - Compress because encoder version mismatches`,
				);

				needToCompress = true;
				needToAddEncoderTag = true;
			}
		}

		if (needToCompress) {
			parentPort?.postMessage(
				`Thread ${threadNumber} - '${flacFile}' - Compressing...`,
			);

			await compressFile(flacFile, encoderSettings);
		}

		if (addEncoderTag && needToAddEncoderTag) {
			await removeTag(flacFile, encoderTag);
			await addTag(flacFile, encoderTag, encoderSettings);
		}
	}
};

worker(workerData.threadNumber, workerData.files, workerData.options);
