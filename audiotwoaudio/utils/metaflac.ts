import { executeCommand } from './exec';

export const getTag = async (flacFile: string, tagName: string) => {
	const stdout = await executeCommand(
		`metaflac --show-tag=${tagName} "${flacFile}"`,
	);

	return stdout.replace(`${tagName}=`, '').trim();
};

export const getVendorTag = async (flacFile: string) => {
	const stdout = await executeCommand(
		`metaflac --show-vendor-tag "${flacFile}"`,
	);

	const vendorTagRegex = new RegExp('(libFLAC) (d+[.]d+[.]d+)');

	return vendorTagRegex.exec(stdout)?.at(1);
};

export const addTag = async (
	flacFile: string,
	tagName: string,
	tagContent: string,
) =>
	await executeCommand(
		`metaflac "--set-tag=${tagName}=${tagContent}" "${flacFile}"`,
	);

export const removeTag = async (flacFile: string, tagName: string) =>
	await executeCommand(`metaflac "--remove-tag=${tagName}" "${flacFile}"`);

export const removeBlock = async (flacFile: string, blockType: string) =>
	await executeCommand(
		`metaflac --remove "--block-type=${blockType}" --dont-use-padding "${flacFile}"`,
	);

export const removeAll = async (flacFile: string) =>
	await executeCommand(
		`metaflac --remove-all --dont-use-padding "${flacFile}"`,
	);

export const addReplayGain = async (flacFiles: string[]) =>
	await executeCommand(
		`metaflac --add-replay-gain ${flacFiles
			.map((flacFile) => `"${flacFile}"`)
			.join(' ')}`,
	);
