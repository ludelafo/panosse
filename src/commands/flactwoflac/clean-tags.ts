import { workerData, parentPort } from 'worker_threads';
import { addTag, getTag, removeAllTags } from '../../utils/metaflac';
import { CommonOptions } from '../../utils/types';

export type TagsOptions = CommonOptions & {
	tagsToKeep: string[];
	fillMissingTags: boolean;
};

const worker = async (
	threadNumber: number,
	flacFiles: string[],
	options: TagsOptions,
) => {
	const {
		dryRun,
		tagsToKeep,
		fillMissingTags
	} = options;

	for await (const flacFile of flacFiles) {
		parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}'`);

		const tagsKept: Map<string, string> = new Map<string, string>();

		for await (const tag of tagsToKeep) {
			const tagContent = await getTag(flacFile, tag);

			if (!tagContent && fillMissingTags) {
				tagsKept.set(tag, `[No ${tag.toLowerCase()}]`);
			} else {
				tagsKept.set(tag, tagContent);
			}
		}

		if (tagsKept.get('TRACKNUMBER')) {
			const trackNumber = Number(tagsKept.get('TRACKNUMBER'));

			tagsKept.set('TRACKNUMBER', trackNumber.toString());
		}

		if (tagsKept.get('TOTALTRACKS')) {
			const trackNumber = Number(tagsKept.get('TOTALTRACKS'));

			tagsKept.set('TOTALTRACKS', trackNumber.toString());
		}

		if (tagsKept.get('DISCNUMBER')) {
			const trackNumber = Number(tagsKept.get('DISCNUMBER'));

			tagsKept.set('DISCNUMBER', trackNumber.toString());
		}

		if (tagsKept.get('TOTALDISKS')) {
			const trackNumber = Number(tagsKept.get('TOTALDISCKS'));

			tagsKept.set('TOTALDISCKS', trackNumber.toString());
		}

		parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}' - Removing all tags...`);

		if (!dryRun) {
			await removeAllTags(flacFile);
		}

		for await (const [tagName, tagContent] of tagsKept) {
			parentPort?.postMessage(`Thread ${threadNumber} - '${flacFile}' - Adding tag...`);

			if (!dryRun) {
				await addTag(flacFile, tagName, tagContent);
			}
		}
	}
};

worker(workerData.threadNumber, workerData.tasks, workerData.options);
