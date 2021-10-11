import { workerData, parentPort } from 'worker_threads';
import { addTag, getTag, removeAll } from '../utils/metaflac';

export type TagsOptions = {
	tagsToKeep: string[];
	fillMissingTags: boolean;
};

const worker = async (
	threadNumber: number,
	flacFiles: string[],
	options: TagsOptions,
) => {
	const { tagsToKeep, fillMissingTags } = options;

	for await (const flacFile of flacFiles) {
		parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}'`);

		const tagsKept: Map<string, string> = new Map<string, string>();

		for await (const tag of tagsToKeep) {
			const tagContent = await getTag(flacFile, tag);

			if (!tagContent && fillMissingTags) {
				tagsKept.set(tag, `No ${tag.toLowerCase()}`);
			} else {
				tagsKept.set(tag, tagContent);
			}
		}

		await removeAll(flacFile);

		for await (const [tagName, tagContent] of tagsKept) {
			await addTag(flacFile, tagName, tagContent);
		}
	}
};

worker(workerData.threadNumber, workerData.files, workerData.options);
