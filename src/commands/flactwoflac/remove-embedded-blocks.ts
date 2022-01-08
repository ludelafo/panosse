import { workerData, parentPort } from 'worker_threads';
import { removeBlock } from '../../utils/metaflac';
import { CommonOptions } from '../../utils/types';

export type EmbeddedBlocksOptions = CommonOptions & {
	removeApplicationBlocks: boolean;
	removeCuesheetBlocks: boolean;
	removePaddingBlocks: boolean;
	removePictureBlocks: boolean;
	removeSeektableBlocks: boolean;
};

const worker = async (
	threadNumber: number,
	flacFiles: string[],
	options: EmbeddedBlocksOptions,
) => {
	const {
		dryRun,
		force,
		removeApplicationBlocks,
		removeCuesheetBlocks,
		removePaddingBlocks,
		removePictureBlocks,
		removeSeektableBlocks,
	} = options;

	for await (const flacFile of flacFiles) {
		parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}'`);

		if (removeApplicationBlocks || force) {
			parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}' - Remove application blocks`);

			if (!dryRun) {
				await removeBlock(flacFile, 'APPLICATION');
			}
		}

		if (removeCuesheetBlocks || force) {
			parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}' - Remove cuesheet blocks`);

			if (!dryRun) {
				await removeBlock(flacFile, 'CUESHEET');
			}
		}

		if (removePictureBlocks || force) {
			parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}' - Remove picture blocks`);

			if (!dryRun) {
				await removeBlock(flacFile, 'PICTURE');
			}
		}

		if (removeSeektableBlocks || force) {
			parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}' - Remove seektable blocks`);

			if (!dryRun) {
				await removeBlock(flacFile, 'SEEKTABLE');
			}
		}

		if (removePaddingBlocks || force) {
			parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}' - Remove padding blocks`);

			if (!dryRun) {
				await removeBlock(flacFile, 'PADDING');
			}
		}
	}
};

worker(workerData.threadNumber, workerData.tasks, workerData.options);
