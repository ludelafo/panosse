import { workerData, parentPort } from 'worker_threads';
import { removeBlock } from '../utils/metaflac';

export type EmbeddedBlocksOptions = {
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
		removeApplicationBlocks,
		removeCuesheetBlocks,
		removePaddingBlocks,
		removePictureBlocks,
		removeSeektableBlocks,
	} = options;

	for await (const flacFile of flacFiles) {
		parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}'`);

		if (removeApplicationBlocks) {
			await removeBlock(flacFile, 'APPLICATION');
		}

		if (removeCuesheetBlocks) {
			await removeBlock(flacFile, 'CUESHEET');
		}

		if (removePaddingBlocks) {
			await removeBlock(flacFile, 'PICTURE');
		}

		if (removePictureBlocks) {
			await removeBlock(flacFile, 'PADDING');
		}

		if (removeSeektableBlocks) {
			await removeBlock(flacFile, 'SEEKTABLE');
		}
	}
};

worker(workerData.threadNumber, workerData.files, workerData.options);
